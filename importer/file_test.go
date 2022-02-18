package importer_test

import (
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	testContent := "this is some dummy content"
	size := int64(len(testContent))

	f := &importer.File{
		Name:        "testing.css",
		ReadCloser:  io.NopCloser(strings.NewReader(testContent)),
		SizeInBytes: &size,
	}

	var content []byte
	concat := func(currentChunk, totalChunks, totalSize int, mimetype string, tmpFile *os.File) error {
		tmpFileContent, err := ioutil.ReadFile(tmpFile.Name())
		if err != nil {
			return err
		}
		content = append(content, tmpFileContent...)
		return nil
	}

	totalChunks, totalSize, err := f.SplitAndClose(size/4, concat)

	assert.Nil(t, err)
	assert.Equal(t, size, totalSize)
	assert.EqualValues(t, 4, totalChunks)
	assert.Equal(t, testContent, string(content))
}
