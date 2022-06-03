package importer_test

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	mocks_importer "github.com/ONSdigital/dp-interactives-importer/importer/mocks"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	downstreamErr              = errors.New("downstream upload error")
	downstreamAlreadyExistsErr = errors.New(importer.DuplicateFileErr)
)

func TestUploadService(t *testing.T) {

	Convey("Given a valid file", t, func() {

		testContent := "this is some dummy content"
		filename := "root/dir/testing.css"
		size := int64(len(testContent))
		f := &importer.File{
			Name:        filename,
			ReadCloser:  ioutil.NopCloser(strings.NewReader(testContent)),
			SizeInBytes: size,
		}

		Convey("And a healthy upload service backend", func() {
			mockBackend := &mocks_importer.UploadServiceBackendMock{
				UploadFunc: func(context.Context, io.ReadCloser, upload.Metadata) error {
					return nil
				},
			}
			svc := importer.NewUploadService(mockBackend)

			Convey("Then there should be no error when we send the file", func() {
				f, err := svc.SendFile(context.TODO(), getTestEvent(filename), f)

				So(err, ShouldBeNil)
				So(f, ShouldEqual, "interactives/collection_id/id/version-1/root/dir/testing.css")
			})

			Convey("Then there should be no error when we send the file and the event has some existing files", func() {
				f, err := svc.SendFile(context.TODO(), getTestEvent(filename, "/interactives/id/version-2/root/dir/testing.css"), f)

				So(err, ShouldBeNil)
				So(f, ShouldEqual, "interactives/collection_id/id/version-3/root/dir/testing.css")
				So(len(mockBackend.UploadCalls()), ShouldEqual, 1)
			})

			Convey("Then there should be no error when we send the file and the event has some existing files without versioned path", func() {
				f, err := svc.SendFile(context.TODO(), getTestEvent(filename, "/interactives/id/root/dir/testing.css"), f)

				So(err, ShouldBeNil)
				So(f, ShouldEqual, "interactives/collection_id/id/version-1/root/dir/testing.css")
				So(len(mockBackend.UploadCalls()), ShouldEqual, 1)
			})
		})

		Convey("And an upload service backend that fails on upload", func() {
			mockBackend := &mocks_importer.UploadServiceBackendMock{
				UploadFunc: func(context.Context, io.ReadCloser, upload.Metadata) error {
					return downstreamErr
				},
			}
			svc := importer.NewUploadService(mockBackend)

			Convey("Then there should be an expected error when we send the file", func() {
				f, err := svc.SendFile(context.TODO(), getTestEvent(filename), f)

				So(err, ShouldNotBeNil)
				So(f, ShouldBeEmpty)
				So(err, ShouldBeError, downstreamErr)
			})
		})

		Convey("And an upload service backend that fails on upload because already exists", func() {
			mockBackend := &mocks_importer.UploadServiceBackendMock{
				UploadFunc: func(_ context.Context, _ io.ReadCloser, u upload.Metadata) error {
					if strings.Contains(u.Path, "version-3") {
						return nil
					}
					return downstreamAlreadyExistsErr
				},
			}
			svc := importer.NewUploadService(mockBackend)

			Convey("Then there should be no error after attempting with valid version", func() {
				f, err := svc.SendFile(context.TODO(), getTestEvent(filename), f)

				So(err, ShouldBeNil)
				So(f, ShouldEqual, "interactives/collection_id/id/version-3/root/dir/testing.css")
				So(len(mockBackend.UploadCalls()), ShouldEqual, 3)
			})
		})
	})
}

func getTestEvent(filename string, in ...string) *importer.InteractivesUploaded {
	var existing []string
	existing = append(existing, in...)

	return &importer.InteractivesUploaded{
		CollectionID: "collection_id",
		ID:           "id",
		Path:         filename,
		Title:        "title",
		CurrentFiles: existing,
	}
}
