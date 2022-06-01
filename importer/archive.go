package importer

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/h2non/filetype"
	"github.com/pkg/errors"
)

var (
	ErrNoIndexHtml          = errors.New("interactive must contain 1 index.html (in root folder)")
	ErrMoreThanOneIndexHtml = errors.New("there can only be 1 index.html in an interactive")
	fileMatchersToIgnore    = []matcher{
		//hidden files
		func(f string) bool { return f[0] == '.' },
		//MACOSX created when right-click, compress: https://superuser.com/questions/104500/what-is-macosx-folder
		func(f string) bool { return f == "__MACOSX" },
		//https://en.wikipedia.org/wiki/Windows_thumbnail_cache
		func(f string) bool { return f == "Thumbs.db" },
	}
)

type matcher func(string) bool

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

	hasIndexHtml := false
	for _, f := range zipReader.File {
		if IsRegular(f) {
			mimetype, err := MimeType(f)
			if err != nil {
				return fmt.Errorf("cannot determine mime type: %s", f.Name)
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			if strings.EqualFold(filepath.Base(f.Name), "index.html") {
				if hasIndexHtml {
					return ErrMoreThanOneIndexHtml
				}
				// Check that the above index file is in root
				if strings.EqualFold(filepath.Clean(f.Name), "index.html") {
					hasIndexHtml = true
				}
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

	if !hasIndexHtml {
		return ErrNoIndexHtml
	}

	return nil
}

func (a *Archive) Close() {
	if e := a.ReadCloser.Close(); e != nil {
		log.Warn(a.Context, "cannot close archive", log.Data{"error": e.Error()})
	}
}

func IsRegular(f *zip.File) bool {
	b := filepath.Base(f.Name)
	ignore := !f.Mode().IsRegular()
	for _, m := range fileMatchersToIgnore {
		if ignore {
			break
		}
		ignore = m(b)
	}
	return !ignore
}

func MimeType(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(f.Name)
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
