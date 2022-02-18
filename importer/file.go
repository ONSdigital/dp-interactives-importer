package importer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
)

type FileProcessor func(currentChunk, totalChunks, totalSize int, mimetype string, tmpFile *os.File) error

type File struct {
	Name        string
	ReadCloser  io.ReadCloser
	SizeInBytes *int64
}

func (f *File) SplitAndClose(chunkSize int64, doFunc FileProcessor) (totalChunks, totalSize int64, err error) {
	mimetype := mime.TypeByExtension(filepath.Ext(f.Name))
	if mimetype == "" {
		return 0, 0, fmt.Errorf("invalid file extension cannot determine mime type: %s", f.Name)
	}

	currentChunk, totalChunks, totalSize := 1, *f.SizeInBytes/chunkSize, *f.SizeInBytes

	r := bufio.NewReader(f.ReadCloser)
	for {
		var n int
		buf := make([]byte, chunkSize)
		if n, err = r.Read(buf); n == 0 {
			if err != io.EOF && err != nil {
				return
			}
			break //all done
		}

		var tmp *os.File
		if tmp, err = ioutil.TempFile("", "chunk_*"); err != nil {
			return
		}

		if _, err = tmp.Write(buf[:n]); err != nil {
			return
		}

		if err = doFunc(currentChunk, int(totalChunks), int(totalSize), mimetype, tmp); err != nil {
			return
		}

		if err = os.Remove(tmp.Name()); err != nil {
			return
		}

		currentChunk++
	}

	if err = f.ReadCloser.Close(); err != nil {
		return
	}

	return
}
