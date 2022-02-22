package importer

import (
	"context"
	s3client "github.com/ONSdigital/dp-s3"
	"github.com/ONSdigital/log.go/v2/log"
)

// VisualisationUploadedHandler ...
type VisualisationUploadedHandler struct {
	S3UploadBucket string
	S3Interface    S3Interface
}

func (h *VisualisationUploadedHandler) Handle(ctx context.Context, event *VisualisationUploaded) error {
	logData := log.Data{"message_id": event.ID}
	log.Info(ctx, "event handler", logData)

	// 1. Download zip file from s3
	// 2. Parse/process
	// 3. Upload visualisations
	// 4. Respond to api  (kafka or rest ?)

	// todo added below impl just to get basic test passing for component-testing framework
	// s3 - mocked, we check assert against num of calls
	err := h.S3Interface.UploadPart(ctx, &s3client.UploadPartRequest{}, nil)
	if err != nil {
		return err
	}
	// todo end

	return nil
}
