package importer

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/ONSdigital/log.go/v2/log"
	"io"
	"mime"
	"os"
	"path/filepath"
)

type Archive struct {
	Context    context.Context
	ReadCloser io.ReadCloser
	Files      []*File
}

func (a *Archive) OpenAndValidate() error {
	tmp, err := os.CreateTemp("", "zip_*")
	if err != nil {
		return err
	}
	if _, err = io.Copy(tmp, a.ReadCloser); err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		if f.Mode().IsRegular() {
			extension := filepath.Ext(f.Name)
			mimetype := mime.TypeByExtension(extension)
			if mimetype == "" {
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
