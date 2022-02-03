package importer

import (
	"context"

	kafka "github.com/ONSdigital/dp-kafka/v2"
)

// VisualisationUploadedHandler ...
type VisualisationUploadedHandler struct {
	S3UploadBucket string
	S3Interface    S3Interface
	Producer       kafka.IProducer
}

func (h *VisualisationUploadedHandler) Handle(ctx context.Context, event *VisualisationUploaded) (err error) {
	// 1. Download zip file from s3
	// 2. Parse/process
	// 3. Upload visualisations
	// 4. Respond to api  (kafka or rest ?)
	return nil
}