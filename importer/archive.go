package importer

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/h2non/filetype"
	"github.com/pkg/errors"
	"io"
	"mime"
	"path/filepath"
	"strings"
)

var (
	fileMatchersToIgnore = []matcher{
		//hidden files
		func(dir, name string) bool { return name[0] == '.' },
		//MACOSX created when right-click, compress: https://superuser.com/questions/104500/what-is-macosx-folder
		func(dir, name string) bool { return name == "__MACOSX" || strings.Contains(dir, "__MACOSX") },
		//https://en.wikipedia.org/wiki/Windows_thumbnail_cache
		func(dir, name string) bool { return name == "Thumbs.db" },
	}
)

type matcher func(string, string) bool

type Archive struct {
	Context    context.Context
	ReadCloser io.ReadCloser
	Files      []*File
}

type File struct {
	Context     context.Context
	ReadCloser  io.ReadCloser
	Name        string
	MimeType    string
	SizeInBytes int64
	Closed      bool
}

func (a *Archive) OpenAndValidate() error {
	//need to read it all for archives
	raw, err := io.ReadAll(a.ReadCloser)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		if IsRegular(f) {
			mimetype, err := MimeType(f)
			if err != nil {
				return fmt.Errorf("cannot determine mime type: %s %w", f.Name, err)
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			size := int64(f.UncompressedSize64)
			a.Files = append(a.Files, &File{
				Context:     a.Context,
				Name:        f.Name,
				ReadCloser:  rc,
				SizeInBytes: size,
				MimeType:    mimetype,
			})
		}
	}

	return nil
}

func (a *Archive) Close() {
	if e := a.ReadCloser.Close(); e != nil {
		log.Warn(a.Context, "cannot close archive", log.Data{"error": e.Error()})
	}
}

func IsRegular(f *zip.File) bool {
	ignore := !f.Mode().IsRegular()
	for _, m := range fileMatchersToIgnore {
		if ignore {
			break
		}
		ignore = m(filepath.Dir(f.Name), filepath.Base(f.Name))
	}
	return !ignore
}

func MimeType(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(f.Name)
	if extension == ".geojson" {
		return "application/geo+json", rc.Close()
	}

	mimetype := mime.TypeByExtension(extension)
	if mimetype == "" {
		kind, _ := filetype.MatchReader(rc)
		if kind == filetype.Unknown {
			return "", errors.New("type unknown")
		}
		mimetype = kind.MIME.Value
	}

	return mimetype, rc.Close()
}
