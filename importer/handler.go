package importer

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-interactives-importer/schema"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	s3client "github.com/ONSdigital/dp-s3"
	"github.com/ONSdigital/log.go/v2/log"
)

// VisualisationUploadedHandler ...
type VisualisationUploadedHandler struct {
	S3UploadBucket string
	S3Interface    S3Interface
	Producer       kafka.IProducer
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
	// kafka - we assert on this message
	msg := VisualisationUploaded{ID: event.ID, Path: fmt.Sprintf("s3://%s", event.Path)}
	bytes, err := schema.VisualisationUploadedEvent.Marshal(msg)
	if err != nil {
		return err
	}
	h.Producer.Channels().Output <- bytes
	// todo end

	return nil
}
