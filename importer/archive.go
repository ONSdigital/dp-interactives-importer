package importer

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/pkg/errors"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	EmptyProcessor       = func(uint64, string, *zip.File) error { return nil }
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

type File struct {
	Context     context.Context
	ReadCloser  io.ReadCloser
	Name        string
	MimeType    string
	SizeInBytes int64
	Closed      bool
}

func Process(batchSize int, z string, processor func(count uint64, mimetype string, zip *zip.File) error) error {
	zipReader, err := zip.OpenReader(z)
	if err != nil {
		return err
	}

	var count uint64
	var validationErrs []error
	var wg sync.WaitGroup

	ch := make(chan struct{}, batchSize)

	for _, f := range zipReader.File {
		file := f
		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer wg.Done()
			skip, mimetype, err := ValidateZipFile(file)
			if err != nil {
				validationErrs = append(validationErrs, fmt.Errorf("cannot open zip file: %s %w", file.Name, err))
			}

			if !skip {
				currentCount := atomic.AddUint64(&count, 1)
				err = processor(currentCount, mimetype, file)
				if err != nil {
					//should we hit the kill switch here...
					validationErrs = append(validationErrs, fmt.Errorf("cannot process zip file: %s %w", file.Name, err))
				}
			}

			<-ch
		}()
	}
	wg.Wait()

	if len(validationErrs) > 0 {
		return fmt.Errorf("found %d validation errors: %v", len(validationErrs), validationErrs)
	}

	return nil
}

func ValidateZipFile(file *zip.File) (skip bool, mimetype string, err error) {
	if IsRegular(file) {
		mimetype, err = MimeType(file)
		if err != nil {
			err = fmt.Errorf("cannot determine mime type: %s %w", file.Name, err)
			return
		}
	} else {
		return true, "", nil
	}
	return
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
