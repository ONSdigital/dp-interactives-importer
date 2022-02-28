package importer

import (
	"context"
	"github.com/ONSdigital/log.go/v2/log"
)

type InteractivesUploadedHandler struct {
	S3            S3Interface
	UploadService *UploadService
}

func (h *InteractivesUploadedHandler) Handle(ctx context.Context, event *InteractivesUploaded) error {
	logData := log.Data{"message_id": event.ID, "path": event.Path}

	// Download zip file from s3
	//todo handle paths???? /my-dir/my-dir-again/file.css
	readCloser, zipSize, err := h.S3.Get(event.Path)
	if err != nil {
		log.Error(ctx, "cannot get zip from s3", err, logData)
		return err
	}
	logData["zip_size"] = zipSize

	// todo Sanity check - do we need to check for 0 size here?

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
	}

	// Respond to api  (kafka or rest ?)

	log.Info(ctx, "successfully processed", logData)

	return nil
}
