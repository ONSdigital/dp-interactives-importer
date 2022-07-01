package importer_test

import (
	"archive/zip"
	"bytes"
	"embed"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/internal/test"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	batchSize = 10

	//go:embed test/single-interactive.zip
	validZipFile embed.FS

	//go:embed test/dvc1774.zip
	invalidZipFile embed.FS
)

func TestArchive(t *testing.T) {

	Convey("Given an empty file", t, func() {
		archive, err := os.CreateTemp("", "test-zip_*.zip")
		So(err, ShouldBeNil)
		defer os.Remove(archive.Name())
		So(archive.Name(), ShouldEndWith, ".zip")
		_, err = io.Copy(archive, strings.NewReader(""))
		So(err, ShouldBeNil)

		Convey("Then there should an error returned when attempt to open", func() {
			err = importer.Process(batchSize, archive.Name(), importer.EmptyProcessor)
			So(err, ShouldBeError, zip.ErrFormat)
		})
	})

	Convey("Given a valid zip file", t, func() {
		archiveName, err := test.CreateTestZip("root.css", "root.html", "root.js", "index.html")
		defer os.Remove(archiveName)
		So(err, ShouldBeNil)
		So(archiveName, ShouldNotBeEmpty)

		Convey("Then open should run successfully", func() {

			var count uint64
			counter := func(uint64, string, *zip.File) error {
				atomic.AddUint64(&count, 1)
				return nil
			}

			err = importer.Process(batchSize, archiveName, counter)
			So(err, ShouldBeNil)

			Convey("And files in archive should be 4", func() {
				So(count, ShouldEqual, 4)
			})
		})
	})

	Convey("Given an actual valid zip file", t, func() {
		Convey("Then open should run successfully", func() {
			err := importer.Process(batchSize, "test/single-interactive.zip", importer.EmptyProcessor)
			So(err, ShouldBeNil)
		})
	})
}

func TestMimeType(t *testing.T) {

	Convey("Given a zip file with some interactives files", t, func() {

		raw, err := validZipFile.ReadFile("test/single-interactive.zip")
		So(err, ShouldBeNil)

		zipReader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
		So(err, ShouldBeNil)

		Convey("Then each file should match expected mimetype", func() {

			expectedMimeTypes := map[string]string{
				"index.html":                              "text/html; charset=utf-8",
				"config.json":                             "application/json",
				"exports.csv":                             "text/csv; charset=utf-8",
				"trade_exports.csv":                       "text/csv; charset=utf-8",
				"css/chosen.css":                          "text/css; charset=utf-8",
				"css/styles.css":                          "text/css; charset=utf-8",
				"js/base.js":                              "application/javascript",
				"js/chosen.jquery.js":                     "application/javascript",
				"js/DataStructures.Tree.js":               "application/javascript",
				"js/modernizr.custom.56904.js":            "application/javascript",
				"fonts/glyphicons-halflings-regular.woff": "font/woff",
				"atlas-tiles/geo/E00135356.geojson":       "application/geo+json",
			}

			var count int
			for _, f := range zipReader.File {
				m, ok := expectedMimeTypes[f.Name]
				if ok {
					mimeType, err := importer.MimeType(f)
					So(err, ShouldBeNil)
					So(mimeType, ShouldEqual, m)
					count++
				}
			}

			//minus directories (js/css/font/root)
			So(len(zipReader.File)-4, ShouldEqual, count)
		})
	})
}

func TestIsRegular(t *testing.T) {

	Convey("Given a regular file IsRegular should be true", t, func() {
		f := &zip.File{FileHeader: zip.FileHeader{Name: "regular"}}
		b := importer.IsRegular(f)
		So(b, ShouldBeTrue)
		f = &zip.File{FileHeader: zip.FileHeader{Name: "/dir1/dir2/regular"}}
		b = importer.IsRegular(f)
		So(b, ShouldBeTrue)
	})

	Convey("Given a hidden file IsRegular should be false", t, func() {
		f := &zip.File{FileHeader: zip.FileHeader{Name: ".hidden"}}
		b := importer.IsRegular(f)
		So(b, ShouldBeFalse)
		f = &zip.File{FileHeader: zip.FileHeader{Name: "/dir1/dir2/.hidden"}}
		b = importer.IsRegular(f)
		So(b, ShouldBeFalse)
		f = &zip.File{FileHeader: zip.FileHeader{Name: "/interactives/label-diKJI1pJ/__MACOSX/atlas-tiles/._index.html"}}
		b = importer.IsRegular(f)
		So(b, ShouldBeFalse)
	})

	Convey("Given a file from a MacOS compressed zip file IsRegular should be false", t, func() {
		f := &zip.File{FileHeader: zip.FileHeader{Name: "__MACOSX"}}
		b := importer.IsRegular(f)
		So(b, ShouldBeFalse)
		f = &zip.File{FileHeader: zip.FileHeader{Name: "/dir1/dir2/__MACOSX"}}
		b = importer.IsRegular(f)
		So(b, ShouldBeFalse)
		f = &zip.File{FileHeader: zip.FileHeader{Name: "/dir1/dir2/__MACOSX/subdir/name"}}
		b = importer.IsRegular(f)
		So(b, ShouldBeFalse)
	})

	Convey("Given a file from a Windows compressed zip file IsRegular should be false", t, func() {
		f := &zip.File{FileHeader: zip.FileHeader{Name: "Thumbs.db"}}
		b := importer.IsRegular(f)
		So(b, ShouldBeFalse)
		f = &zip.File{FileHeader: zip.FileHeader{Name: "/dir1/dir2/Thumbs.db"}}
		b = importer.IsRegular(f)
		So(b, ShouldBeFalse)
	})
}
