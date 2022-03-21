package importer_test

import (
	"archive/zip"
	"bytes"
	"context"
	"embed"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/internal/test"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	//go:embed test/single-interactive.zip
	zipFile embed.FS
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
			a := &importer.Archive{Context: context.TODO(), ReadCloser: archive}
			err = a.OpenAndValidate()
			So(err, ShouldBeError, zip.ErrFormat)
		})
	})

	Convey("Given an invalid zip file (no index.html)", t, func() {
		archiveName, err := test.CreateTestZip("root.css", "root.html", "root.js")
		defer os.Remove(archiveName)
		So(err, ShouldBeNil)
		So(archiveName, ShouldNotBeEmpty)

		Convey("Then open should run successfully", func() {
			archive, err := os.Open(archiveName)
			So(err, ShouldBeNil)

			a := &importer.Archive{Context: context.TODO(), ReadCloser: archive}
			err = a.OpenAndValidate()
			So(err, ShouldEqual, importer.ErrNoIndexHtml)
		})
	})

	Convey("Given an invalid zip file (multiple index.html)", t, func() {
		archiveName, err := test.CreateTestZip("root.css", "root.html", "root.js", "index.html", "test/index.html")
		defer os.Remove(archiveName)
		So(err, ShouldBeNil)
		So(archiveName, ShouldNotBeEmpty)

		Convey("Then open should run successfully", func() {
			archive, err := os.Open(archiveName)
			So(err, ShouldBeNil)

			a := &importer.Archive{Context: context.TODO(), ReadCloser: archive}
			err = a.OpenAndValidate()
			So(err, ShouldEqual, importer.ErrMoreThanOneIndexHtml)
		})
	})

	Convey("Given a valid zip file", t, func() {
		archiveName, err := test.CreateTestZip("root.css", "root.html", "root.js", "index.html")
		defer os.Remove(archiveName)
		So(err, ShouldBeNil)
		So(archiveName, ShouldNotBeEmpty)

		Convey("Then open should run successfully", func() {
			archive, err := os.Open(archiveName)
			So(err, ShouldBeNil)

			a := &importer.Archive{Context: context.TODO(), ReadCloser: archive}
			err = a.OpenAndValidate()
			So(err, ShouldBeNil)

			Convey("And files in archive should be 4", func() {
				So(len(a.Files), ShouldEqual, 4)
			})
		})
	})
}

func TestMimeType(t *testing.T) {

	Convey("Given a zip file with some interactives files", t, func() {

		raw, err := zipFile.ReadFile("test/single-interactive.zip")
		So(err, ShouldBeNil)

		zipReader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
		So(err, ShouldBeNil)

		Convey("Then each file should match expected mimetype", func() {

			expectedMimeTypes := map[string]string{
				"single-interactive/index.html":                              "text/html; charset=utf-8",
				"single-interactive/config.json":                             "application/json",
				"single-interactive/exports.csv":                             "text/csv; charset=utf-8",
				"single-interactive/trade_exports.csv":                       "text/csv; charset=utf-8",
				"single-interactive/css/chosen.css":                          "text/css; charset=utf-8",
				"single-interactive/css/styles.css":                          "text/css; charset=utf-8",
				"single-interactive/js/base.js":                              "application/javascript",
				"single-interactive/js/chosen.jquery.js":                     "application/javascript",
				"single-interactive/js/DataStructures.Tree.js":               "application/javascript",
				"single-interactive/js/modernizr.custom.56904.js":            "application/javascript",
				"single-interactive/fonts/glyphicons-halflings-regular.woff": "font/woff",
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
	})

	Convey("Given a file from a MacOS compressed zip file IsRegular should be false", t, func() {
		f := &zip.File{FileHeader: zip.FileHeader{Name: "__MACOSX"}}
		b := importer.IsRegular(f)
		So(b, ShouldBeFalse)
		f = &zip.File{FileHeader: zip.FileHeader{Name: "/dir1/dir2/__MACOSX"}}
		b = importer.IsRegular(f)
		So(b, ShouldBeFalse)
	})
}
