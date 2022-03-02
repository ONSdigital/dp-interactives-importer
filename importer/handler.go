package importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-interactives-importer/config"
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

func (h *InteractivesUploadedHandler) Handle(ctx context.Context, event *InteractivesUploaded) (err error) {
	logData := log.Data{"id": event.ID, "path": event.Path}

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
				Metadata: make(map[string]string),
			},
		}
		if err != nil {
			logData["error"] = err.Error()
			update.ImportSuccessful = &failure
			update.Interactive.Metadata["error"] = err.Error()
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
		err = h.UploadService.SendFile(ctx, f, "title", "collectionId", "licence", "licenceUrl")
		if err != nil {
			log.Error(ctx, "failed to upload file", err, log.Data{"file": f})
			return err
		}
		archiveFiles = append(archiveFiles, &interactives.InteractiveFile{
			Name:     f.Name,
			Mimetype: f.MimeType,
			Size:     f.SizeInBytes,
		})
	}

	log.Info(ctx, "successfully processed", logData)

	return nil
}
