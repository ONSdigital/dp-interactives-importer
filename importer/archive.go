package importer

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/h2non/filetype"
	"github.com/pkg/errors"
	"io"
	"mime"
	"os"
	"path/filepath"
)

type Archive struct {
	Context    context.Context
	ReadCloser io.ReadCloser
	Files      []*File
	TmpRemoved bool
}

func (a *Archive) OpenAndValidate() error {
	tmp, err := os.CreateTemp("", "zip_*")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if e := os.Remove(f.Name()); e != nil {
			logData := log.Data{"error": e.Error(), "name": f.Name()}
			log.Warn(a.Context, "cannot remove tmp file", logData)
		}
	}(tmp)

	if _, err = io.Copy(tmp, a.ReadCloser); err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return err
	}

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
		log.Warn(a.Context, "cannot close zip file", log.Data{"error": e.Error()})
	}
}

func IsRegular(f *zip.File) bool {
	//MACOSX created when right-click, compress: https://superuser.com/questions/104500/what-is-macosx-folder
	b := filepath.Base(f.Name)
	return f.Mode().IsRegular() && b[0] != '.' && b != "__MACOSX"
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
