package importer_test

import (
	"github.com/ONSdigital/dp-interactives-importer/importer"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var (
	doNothing = func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error { return nil }
)

func getConcatFunc(content *[]byte) func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error {
	return func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error {
		tmpFileContent, e := ioutil.ReadFile(tmpFile.Name())
		if e != nil {
			return e
		}
		*content = append(*content, tmpFileContent...)
		return nil
	}
}

func TestFile(t *testing.T) {

	Convey("Given an empty file with an valid extension", t, func() {
		var zero int64
		f := &importer.File{
			Name:        "testing.css",
			ReadCloser:  io.NopCloser(strings.NewReader("")),
			SizeInBytes: &zero,
		}

		Convey("When split and concatenate response into memory", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(zero+1024, getConcatFunc(&content))

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the content should be empty with no chunks processed", func() {
				So(string(content), ShouldBeEmpty)
				So(totalChunks, ShouldBeZeroValue)
			})

			Convey("And the file is closed", func() {
				So(f.Closed, ShouldBeTrue)
			})
		})
	})

	testContent := "this is some dummy content"
	size := int64(len(testContent))

	Convey("Given a file with an invalid extension", t, func() {
		f := &importer.File{
			Name:        "testing",
			ReadCloser:  io.NopCloser(strings.NewReader(testContent)),
			SizeInBytes: &size,
		}

		Convey("When attempt to split and do nothing with each chunk", func() {
			totalChunks, err := f.SplitAndClose(size, doNothing)

			Convey("Then there should an error returned with the filename", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, f.Name)
				So(totalChunks, ShouldBeZeroValue)
			})

			Convey("And the file is closed", func() {
				So(f.Closed, ShouldBeTrue)
			})
		})
	})

	Convey("Given a file with valid content and extension", t, func() {
		f := &importer.File{
			Name:        "testing.css",
			ReadCloser:  io.NopCloser(strings.NewReader(testContent)),
			SizeInBytes: &size,
		}

		Convey("When split into 4 chunks", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(size/4, getConcatFunc(&content))

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the content should match with a total of 5 chunks processed", func() {
				So(string(content), ShouldEqual, testContent)
				So(totalChunks, ShouldEqual, 5)
			})

			Convey("And the file is closed", func() {
				So(f.Closed, ShouldBeTrue)
			})
		})

		Convey("When processed in full with a bigger chunk size", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(size+1024, getConcatFunc(&content))

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the content should match with 1 chunk processed", func() {
				So(string(content), ShouldEqual, testContent)
				So(totalChunks, ShouldEqual, 1)
			})

			Convey("And the file is closed", func() {
				So(f.Closed, ShouldBeTrue)
			})
		})
	})
}
