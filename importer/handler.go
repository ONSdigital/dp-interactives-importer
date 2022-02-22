package importer

import (
	"context"
	"github.com/ONSdigital/log.go/v2/log"
)

// VisualisationUploadedHandler ...
type VisualisationUploadedHandler struct {
	S3            S3Interface
	UploadService *UploadService
}

func (h *VisualisationUploadedHandler) Handle(ctx context.Context, event *VisualisationUploaded) error {
	logData := log.Data{"message_id": event.ID}
	log.Info(ctx, "event handler", logData)

	// Download zip file from s3
	//todo handle paths???? /my-dir/my-dir-again/file.css
	readCloser, size, err := h.S3.Get(event.Path)
	if err != nil {
		return err
	}
	file := &File{
		Name:        event.Path,
		ReadCloser:  readCloser,
		SizeInBytes: size,
	}

	// Parse/process

	// Upload visualisations
	err = h.UploadService.SendFile(ctx, file, "title", "collectionId", "licence", "licenceUrl")
	if err != nil {
		return err
	}

	// Respond to api  (kafka or rest ?)

	return nil
}
