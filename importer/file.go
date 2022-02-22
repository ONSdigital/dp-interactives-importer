package importer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"os"
	"path/filepath"
)

type FileProcessor func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error

type File struct {
	Name        string
	ReadCloser  io.ReadCloser
	SizeInBytes *int64
	Closed      bool
}

func (f *File) Stat(chunkSize int64) (totalExpectedChunks, totalSize int32) {
	if *f.SizeInBytes == 0 {
		return 0, 0
	}
	return int32(math.Ceil(float64(*f.SizeInBytes / chunkSize))), int32(*f.SizeInBytes)
}

func (f *File) SplitAndClose(chunkSize int64, doFunc FileProcessor) (totalChunks int64, err error) {
	defer func(ReadCloser io.ReadCloser) {
		e := ReadCloser.Close()
		if e != nil {
			err = e
			return
		}
		f.Closed = true
	}(f.ReadCloser)

	mimetype := mime.TypeByExtension(filepath.Ext(f.Name))
	if mimetype == "" {
		err = fmt.Errorf("invalid file extension cannot determine mime type: %s", f.Name)
		return
	}

	totalSize, totalExpectedChunks := f.Stat(chunkSize)

	r := bufio.NewReader(f.ReadCloser)
	var currentChunk int32
	for {
		currentChunk++

		var n int
		buf := make([]byte, chunkSize)
		if n, err = r.Read(buf); n == 0 || err != nil {
			if err != nil && err != io.EOF {
				return
			}
			err = nil //dont return io.EOF
			break     //all done
		}

		//todo can we just send the []bytes instead?
		var tmp *os.File
		if tmp, err = ioutil.TempFile("", "chunk_*"); err != nil {
			return
		}

		if _, err = tmp.Write(buf[:n]); err != nil {
			return
		}

		if err = doFunc(currentChunk, totalExpectedChunks, totalSize, mimetype, tmp); err != nil {
			return
		}

		if err = os.Remove(tmp.Name()); err != nil {
			return
		}
	}

	totalChunks = int64(currentChunk - 1)

	return
}
