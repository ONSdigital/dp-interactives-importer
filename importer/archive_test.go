package importer_test

import (
	"archive/zip"
	"context"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/internal/test"
	"io"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArchive(t *testing.T) {

	Convey("Given an empty file", t, func() {
		archive, err := os.CreateTemp("", "test-zip_*.zip")
		So(err, ShouldBeNil)
		So(archive.Name(), ShouldEndWith, ".zip")
		_, err = io.Copy(archive, strings.NewReader(""))
		So(err, ShouldBeNil)

		Convey("Then there should an error returned when attempt to open", func() {
			a := &importer.Archive{Context: context.TODO(), ReadCloser: archive}
			err = a.OpenAndValidate()
			So(err, ShouldBeError, zip.ErrFormat)
		})
	})

	Convey("Given a valid zip file", t, func() {
		archiveName, err := test.CreateTestZip("root.css", "root.html", "root.js")
		So(err, ShouldBeNil)
		So(archiveName, ShouldNotBeEmpty)

		Convey("Then open should run successfully", func() {
			archive, err := os.Open(archiveName)
			So(err, ShouldBeNil)

			a := &importer.Archive{Context: context.TODO(), ReadCloser: archive}
			err = a.OpenAndValidate()
			So(err, ShouldBeNil)

			Convey("And files in archive should be 3", func() {
				So(len(a.Files), ShouldEqual, 3)
			})
		})
	})
}
