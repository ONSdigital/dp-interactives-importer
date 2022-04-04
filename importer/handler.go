package importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/schema"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
	"io"
)

var (
	successful, failure = true, false
)

type InteractivesUploadedHandler struct {
	Cfg                   *config.Config
	S3                    S3Interface
	UploadService         *UploadService
	InteractivesAPIClient InteractivesAPIClient
}

func (h *InteractivesUploadedHandler) Handle(ctx context.Context, workerID int, msg kafka.Message) error {
	logData := log.Data{"workerID": workerID}

	event, err := getAsEvent(ctx, msg)
	if err != nil {
		log.Error(ctx, "cannot unmarshal into an event", err, logData)
		return err
	}

	var readCloser io.ReadCloser
	var zipSize *int64
	var archiveFiles []*interactives.InteractiveFile

	// Defer an update via API - deferred so we always attempt an update!
	defer func() {
		update := interactives.InteractiveUpdate{
			Interactive: interactives.Interactive{
				ID: event.ID,
				Archive: &interactives.InteractiveArchive{
					Name: event.Path,
				},
			},
		}
		if err != nil {
			logData["error"] = err.Error()
			update.ImportSuccessful = &failure
			update.ImportMessage = err.Error()
		} else {
			update.ImportSuccessful = &successful
			update.Interactive.Archive.Size = *zipSize
			update.Interactive.Archive.Files = archiveFiles
		}
		// user token not valid - we auth user on api endpoints
		apiErr := h.InteractivesAPIClient.PutInteractive(ctx, "", h.Cfg.ServiceAuthToken, event.ID, update)
		if apiErr != nil {
			//todo what if this fails - retry?
			logData["apiError"] = apiErr.Error()
			log.Warn(ctx, "failed to update interactive", logData)
		}
	}()

	logData["id"] = event.ID
	logData["path"] = event.Path
	logData["title"] = event.Title
	logData["current_files"] = event.CurrentFiles

	// Download zip file from s3
	readCloser, zipSize, err = h.S3.Get(event.Path)
	if err != nil {
		log.Error(ctx, "cannot get zip from s3", err, logData)
		return err
	}
	logData["zip_size"] = zipSize

	// Open zip and validate contents
	archive := &Archive{Context: ctx, ReadCloser: readCloser}
	err = archive.OpenAndValidate()
	if err != nil {
		log.Error(ctx, "cannot open and validate zip", err, logData)
		return err
	}
	defer archive.Close()
	logData["num_files"] = len(archive.Files)

	// Upload each file in zip
	for _, f := range archive.Files {
		savedFilename, err := h.UploadService.SendFile(ctx, event, f)
		if err != nil {
			log.Error(ctx, "failed to upload file", err, log.Data{"file": f})
			return err
		}
		archiveFiles = append(archiveFiles, &interactives.InteractiveFile{
			Name:     savedFilename,
			Mimetype: f.MimeType,
			Size:     f.SizeInBytes,
		})
	}

	log.Info(ctx, "successfully processed", logData)

	return nil
}

// getAsEvent unmarshals the provided kafka message into an event and calls the handler.
func getAsEvent(ctx context.Context, message kafka.Message) (*InteractivesUploaded, error) {
	defer message.Commit()

	logData := log.Data{"message_offset": message.Offset()}

	var event InteractivesUploaded
	err := schema.InteractivesUploadedEvent.Unmarshal(message.GetData(), &event)
	if err != nil {
		log.Error(ctx, "failed to unmarshal event", err, logData)
		return nil, err
	}

	logData["event"] = event

	log.Info(ctx, "event received", logData)

	return &event, nil
}