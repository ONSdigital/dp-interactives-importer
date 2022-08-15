package importer

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
)

const (
	licenseName = "Open Government Licence v3.0"
	licenseURL  = "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/"
)

func NewUploadService(backend UploadServiceBackend) *UploadService {
	return &UploadService{
		backend: backend,
	}
}

//todo handling retries?

type UploadService struct {
	backend UploadServiceBackend
}

func (s *UploadService) SendFile(ctx context.Context, event *InteractivesUploaded, f *File) (string, error) {
	metadata := upload.Metadata{
		Path:          f.RootPath,
		IsPublishable: true,
		Title:         event.Title,
		FileSizeBytes: f.SizeInBytes,
		FileType:      f.MimeType,
		License:       licenseName,
		LicenseURL:    licenseURL,
		FileName:      f.Name,
		CollectionID:  &f.CollectionID,
	}

	err := s.backend.Upload(ctx, f.ReadCloser, metadata)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", metadata.Path, metadata.FileName), nil
}
