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

type FileProcessor func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error

type File struct {
	Name        string
	ReadCloser  io.ReadCloser
	SizeInBytes *int64
}

func (f *File) SplitAndClose(chunkSize int64, doFunc FileProcessor) (totalChunks, totalSize int64, err error) {
	mimetype := mime.TypeByExtension(filepath.Ext(f.Name))
	if mimetype == "" {
		err = fmt.Errorf("invalid file extension cannot determine mime type: %s", f.Name)
		return
	}

	currentChunk, totalChunks, totalSize := 1, *f.SizeInBytes/chunkSize, *f.SizeInBytes
	if totalChunks == 0 {
		totalChunks = 1
	}

	r := bufio.NewReader(f.ReadCloser)
	for {
		var n int
		buf := make([]byte, chunkSize)
		if n, err = r.Read(buf); n == 0 || err != nil {
			if err != nil && err != io.EOF {
				return
			}
			break //all done
		}

		//todo can we just send the []bytes instead?
		var tmp *os.File
		if tmp, err = ioutil.TempFile("", "chunk_*"); err != nil {
			return
		}

		if _, err = tmp.Write(buf[:n]); err != nil {
			return
		}

		if err = doFunc(int32(currentChunk), int32(totalChunks), int32(totalSize), mimetype, tmp); err != nil {
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
