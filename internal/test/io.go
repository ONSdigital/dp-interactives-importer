package test

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

func CreateTestZip(filenames ...string) (string, error) {
	archive, err := os.CreateTemp("", "test-zip_*.zip")
	if err != nil {
		return "", err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	for _, f := range filenames {
		w, err := zipWriter.Create(f)
		if err != nil {
			return "", err
		}
		if _, err = io.Copy(w, strings.NewReader(f)); err != nil {
			return "", err
		}
	}
	if err := zipWriter.Flush(); err != nil {
		return "", err
	}

	return archive.Name(), zipWriter.Close()
}
