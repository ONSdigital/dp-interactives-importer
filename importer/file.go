package importer

import (
	"bufio"
	"context"
	"github.com/ONSdigital/log.go/v2/log"
	"io"
	"math"
	"os"
)

const (
	DefaultChunkSize = 5 << (10 * 2) //https://github.com/ONSdigital/dp-s3#chunk-size
)

type FileProcessor func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error

type File struct {
	Context     context.Context
	ReadCloser  io.ReadCloser
	Name        string
	MimeType    string
	SizeInBytes int64
	Closed      bool
}

func (f *File) SplitAndClose(chunkSize int64, doFunc FileProcessor) (totalChunks int64, err error) {
	defer func(r io.ReadCloser) {
		if e := r.Close(); e != nil {
			logData := log.Data{"error": e.Error(), "name": f.Name}
			log.Warn(f.Context, "cannot close file", logData)
		}
		f.Closed = true
	}(f.ReadCloser)

	var totalExpectedChunks, totalSize int32
	if f.SizeInBytes > 0 {
		totalExpectedChunks, totalSize = int32(math.Ceil(float64(f.SizeInBytes/chunkSize))), int32(f.SizeInBytes)
	}

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
			if n == 0 {
				//Note: zip files behave slightly differently to normal readers: https://github.com/golang/go/issues/32858
				//      you get n>0 and the EOF error not 0 and EOF like others
				err = nil //dont return io.EOF
				break     //all done
			}
		}

		//todo can we just send the []bytes instead?
		var tmp *os.File
		if tmp, err = os.CreateTemp("", "chunk_*"); err != nil {
			return
		}

		if _, err = tmp.Write(buf[:n]); err != nil {
			return
		}

		if err = doFunc(currentChunk, totalExpectedChunks, totalSize, f.MimeType, tmp); err != nil {
			return
		}

		if err = os.Remove(tmp.Name()); err != nil {
			return
		}
	}

	totalChunks = int64(currentChunk - 1)

	return
}
