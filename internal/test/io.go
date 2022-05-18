package test

import (
	"archive/zip"
	"bufio"
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

func CreateTempFile(content string) (string, int64, error) {
	f, err := os.CreateTemp("", "test-tmp_*.txt")
	if err != nil {
		return "", 0, err
	}

	w := bufio.NewWriter(f)
	if _, err = io.Copy(w, strings.NewReader(content)); err != nil {
		return "", 0, err
	}
	if err := w.Flush(); err != nil {
		return "", 0, err
	}

	stat, err := f.Stat()
	if err != nil {
		return "", 0, err
	}

	return f.Name(), stat.Size(), f.Close()
}
