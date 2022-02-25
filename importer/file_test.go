package importer_test

import (
	"archive/zip"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/internal/test"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
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

	Convey("Given an empty file", t, func() {
		f := &importer.File{
			Name:        "empty",
			ReadCloser:  io.NopCloser(strings.NewReader("")),
			SizeInBytes: 0,
		}

		Convey("When split and concatenate response into memory", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(1024, getConcatFunc(&content))

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

	Convey("Given a valid file", t, func() {
		testContent := "this is some dummy content"
		size := int64(len(testContent))

		f := &importer.File{
			Name:        "valid",
			ReadCloser:  io.NopCloser(strings.NewReader(testContent)),
			SizeInBytes: size,
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

		Convey("When processed in full with default chunk size", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(importer.DefaultChunkSize, getConcatFunc(&content))

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

	Convey("Given an actual file", t, func() {
		testContent := "this is some dummy content"
		tmpFileName, size, err := test.CreateTempFile(testContent)
		defer os.Remove(tmpFileName)
		So(err, ShouldBeNil)
		So(size, ShouldNotBeZeroValue)

		tmpFile, err := os.Open(tmpFileName)
		So(err, ShouldBeNil)

		f := &importer.File{
			Name:        "actual_file",
			ReadCloser:  tmpFile,
			SizeInBytes: size,
		}

		Convey("When split with default chunk size", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(importer.DefaultChunkSize, getConcatFunc(&content))

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

	Convey("Given a file from within a zip", t, func() {
		testContent := "test"
		tmpZipName, err := test.CreateTestZip(testContent)
		defer os.Remove(tmpZipName)
		So(err, ShouldBeNil)

		zipReader, err := zip.OpenReader(tmpZipName)
		defer zipReader.Close()
		So(err, ShouldBeNil)
		So(len(zipReader.File), ShouldEqual, 1)

		zipFile, err := zipReader.File[0].Open()
		size := int64(zipReader.File[0].UncompressedSize64)

		f := &importer.File{
			Name:        "zip_file",
			ReadCloser:  zipFile,
			SizeInBytes: size,
		}

		Convey("When split with default chunk size", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(importer.DefaultChunkSize, getConcatFunc(&content))

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
