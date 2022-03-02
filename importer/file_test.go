package importer_test

import (
	"archive/zip"
	"bytes"
	"embed"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var (
	//go:embed test/styles.css
	testFile embed.FS
)

func getConcatFunc(t *testing.T, content *[]byte) func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error {
	return func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error {
		if totalChunks <= 0 {
			t.Errorf("cannot send <=0 total chunks: %d", totalChunks)
		}
		if currentChunk <= 0 {
			t.Errorf("cannot send <=0 current chunk: %d", currentChunk)
		}
		if totalSize <= 0 {
			t.Errorf("cannot send <=0 totalSize: %d", totalSize)
		}

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
			totalChunks, err := f.SplitAndClose(1024, getConcatFunc(t, &content))

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
			totalChunks, err := f.SplitAndClose(size/4, getConcatFunc(t, &content))

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
			totalChunks, err := f.SplitAndClose(importer.DefaultChunkSize, getConcatFunc(t, &content))

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
		raw, err := testFile.ReadFile("test/styles.css")
		So(err, ShouldBeNil)

		f := &importer.File{
			Name:        "test/styles.css",
			ReadCloser:  io.NopCloser(bytes.NewReader(raw)),
			SizeInBytes: int64(len(raw)),
		}

		Convey("When split with default chunk size", func() {
			var content []byte
			totalChunks, err := f.SplitAndClose(importer.DefaultChunkSize, getConcatFunc(t, &content))

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the content should not be empty with 1 chunk processed", func() {
				So(string(content), ShouldStartWith, "@import url(\"//fonts.googleapis.com")
				So(string(content), ShouldEndWith, "color: #222222;\n}\n")
				So(totalChunks, ShouldEqual, 1)
			})

			Convey("And the file is closed", func() {
				So(f.Closed, ShouldBeTrue)
			})
		})
	})

	Convey("Given an actual sample zip", t, func() {
		raw, err := zipFile.ReadFile("test/single-interactive.zip")
		So(err, ShouldBeNil)

		zipReader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
		So(err, ShouldBeNil)

		var count, totalChunks int
		Convey("When we split and close each file within", func() {
			for _, z := range zipReader.File {
				if z.Mode().IsRegular() {
					count++
					size := int64(z.UncompressedSize64)
					rc, err := z.Open()
					So(err, ShouldBeNil)

					f := &importer.File{
						Name:        z.Name,
						ReadCloser:  rc,
						SizeInBytes: size,
					}

					var content []byte
					tc, err := f.SplitAndClose(importer.DefaultChunkSize, getConcatFunc(t, &content))
					totalChunks = totalChunks + int(tc)

					So(err, ShouldBeNil)
					So(content, ShouldNotBeEmpty)
				}
			}

			Convey("Then total chunks should equal number of files processed", func() {
				So(totalChunks, ShouldEqual, count)
			})
		})
	})
}
