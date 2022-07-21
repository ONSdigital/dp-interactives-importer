package importer_test

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	mocks_importer "github.com/ONSdigital/dp-interactives-importer/importer/mocks"
	"github.com/ONSdigital/log.go/v2/log"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	cfg = &config.Config{}
)

func TestJobFinish(t *testing.T) {
	anErr := errors.New("an error")
	logData := log.Data{}
	rootPath := "/root/path/"
	event := &importer.InteractivesUploaded{
		ID:   "1",
		Path: rootPath,
	}
	var mockInteractivesAPI *mocks_importer.InteractivesAPIClientMock

	Convey("Given a healthy interactives api", t, func() {
		mockInteractivesAPI = &mocks_importer.InteractivesAPIClientMock{
			PatchInteractiveFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, interactiveID string, req interactives.PatchRequest) (interactives.Interactive, error) {
				return interactives.Interactive{}, nil
			},
		}

		Convey("And a new upload job with a size and error passed", func() {
			var err error
			var zipSize int64

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				uploadJob := importer.NewJob(context.TODO(), cfg, mockInteractivesAPI)
				defer uploadJob.Finish(&logData, event, rootPath, &zipSize, &err)
				err = anErr
				wg.Done()
			}()
			wg.Wait()

			Convey("Then there should be an expected error when we add new files concurrently", func() {
				mockPatchReq := mockInteractivesAPI.PatchInteractiveCalls()[0].PatchRequest
				So(mockPatchReq.Interactive.Archive.ImportSuccessful, ShouldBeFalse)
				So(mockPatchReq.Interactive.Archive.ImportMessage, ShouldEqual, "an error")
			})
		})

	})
}
