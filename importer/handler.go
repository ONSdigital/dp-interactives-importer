package importer

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/schema"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"io"
	"os"
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

	var zipSize int64
	//no leading slash: https://github.com/ONSdigital/dp-upload-service/blob/ecc6062e6fe5856385b5fafbe1105606c1a958ff/api/upload.go#L25
	randomString := gonanoid.Must(16)
	uploadRootPath := fmt.Sprintf("%s/%s/%s", "interactives", event.ID, randomString)

	uploadJob := NewJob(ctx, h.Cfg, h.InteractivesAPIClient)
	defer uploadJob.Finish(&logData, event, uploadRootPath, &zipSize, &err) // defer finish() so we always attempt!

	logData["id"] = event.ID
	logData["path"] = event.Path
	logData["title"] = event.Title

	log.Info(ctx, "download zip file from s3", logData)
	readCloser, size, err := h.S3.Get(event.Path)
	if err != nil {
		log.Error(ctx, "cannot get zip from s3", err, logData)
		return err
	}
	zipSize = *size
	tmpZip, err := os.CreateTemp("", "s3-zip_*.zip")
	if err != nil {
		return err
	}
	if _, err = io.Copy(tmpZip, readCloser); err != nil {
		return err
	}
	tmpZip.Close()
	defer os.Remove(tmpZip.Name())
	logData["zip_size"] = zipSize

	log.Info(ctx, "validate zip", logData)
	counterFunc := func(count uint64, mimetype string, zip *zip.File) error {
		if count%1000 == 0 {
			log.Info(ctx, "processed 1000 files", logData)
		}
		return nil
	}
	err = Process(h.Cfg.BatchSize, tmpZip.Name(), counterFunc)
	if err != nil {
		log.Error(ctx, "cannot validate zip", err, logData)
		return err
	}

	// Upload each file in zip
	log.Info(ctx, "start upload of zip files", logData)
	uploadFunc := func(count uint64, mimetype string, zip *zip.File) error {
		if count%1000 == 0 {
			log.Info(ctx, "processed 1000 files", logData)
		}

		rc, err := zip.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		size := int64(zip.UncompressedSize64)
		file := &File{
			Context:     ctx,
			Name:        zip.Name,
			ReadCloser:  rc,
			SizeInBytes: size,
			MimeType:    mimetype,
		}

		savedFileName, err := h.UploadService.SendFile(ctx, event, file, uploadRootPath)
		if err != nil {
			return err
		}

		err = uploadJob.Add(event, &interactives.ArchiveFile{
			Name:     savedFileName,
			Size:     size,
			Mimetype: mimetype,
			URI:      zip.Name, //this could be rendered from http://domain/interactives/uri
		})
		return err
	}
	err = Process(h.Cfg.BatchSize, tmpZip.Name(), uploadFunc)
	if err != nil {
		log.Error(ctx, "cannot process zip", err, logData)
		return err
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
