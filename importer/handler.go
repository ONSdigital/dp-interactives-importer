package importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/log.go/v2/log"
)

var (
	successful bool
)

func init() {
	successful = true
}

type InteractivesUploadedHandler struct {
	S3                    S3Interface
	UploadService         *UploadService
	InteractivesAPIClient InteractivesAPIClient
}

func (h *InteractivesUploadedHandler) Handle(ctx context.Context, event *InteractivesUploaded) error {
	logData := log.Data{"id": event.ID, "path": event.Path}

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
	var archiveFiles []*interactives.InteractiveFile
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

	// Update interactive via API todo user & service auth
	err = h.InteractivesAPIClient.PutInteractive(ctx, "", "",
		event.ID,
		interactives.InteractiveUpdate{
			ImportSuccessful: &successful,
			Interactive: interactives.Interactive{
				ID: event.ID,
				Archive: &interactives.InteractiveArchive{
					Name:  event.Path,
					Size:  *zipSize,
					Files: archiveFiles,
				},
			},
		})
	if err != nil {
		log.Error(ctx, "failed to update interactive", err, logData)
		return err
	}

	log.Info(ctx, "successfully processed", logData)

	return nil
}
