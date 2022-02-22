package importer_test

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	mocks_importer "github.com/ONSdigital/dp-interactives-importer/importer/mocks"
	"github.com/ONSdigital/dp-interactives-importer/internal/client/uploadservice"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	downstreamErr = errors.New("downstream upload error")
	chunkSize     = int64(1024)
	checkerOkFunc = func(context.Context, *healthcheck.CheckState) error {
		return nil
	}
	uploadOkFunc = func(context.Context, string, uploadservice.UploadJob) error {
		return nil
	}
	uploadErrFunc = func(context.Context, string, uploadservice.UploadJob) error {
		return downstreamErr
	}
)

func TestUploadService(t *testing.T) {

	Convey("Given a valid file", t, func() {
		testContent := "this is some dummy content"
		size := int64(len(testContent))
		f := &importer.File{
			Name:        "testing.css",
			ReadCloser:  ioutil.NopCloser(strings.NewReader(testContent)),
			SizeInBytes: &size,
		}

		Convey("And a healthy upload service backend", func() {
			mockBackend := &mocks_importer.UploadServiceBackendMock{
				CheckerFunc: checkerOkFunc,
				UploadFunc:  uploadOkFunc,
			}
			svc := importer.NewUploadService(mockBackend, chunkSize)

			Convey("Then there should be no error and 1 upload saved when we send the file", func() {
				err := svc.SendFile(context.TODO(), f, "title", "collectionId", "licence", "licenceUrl")

				So(err, ShouldBeNil)
				So(len(svc.Uploads), ShouldEqual, 1)
			})
		})

		Convey("And an upload service backend that fails on upload", func() {
			mockBackend := &mocks_importer.UploadServiceBackendMock{
				CheckerFunc: checkerOkFunc,
				UploadFunc:  uploadErrFunc,
			}
			svc := importer.NewUploadService(mockBackend, chunkSize)

			Convey("Then there should be an expected error when we send the file", func() {
				err := svc.SendFile(context.TODO(), f, "title", "collectionId", "licence", "licenceUrl")

				So(err, ShouldNotBeNil)
				So(err, ShouldBeError, downstreamErr)
			})
		})
	})
}
