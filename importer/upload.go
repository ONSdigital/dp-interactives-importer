package importer

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
	"regexp"
	"strconv"
	"strings"
)

const (
	maxAttempts         = 10
	uploadRootDirectory = "interactives"
	DuplicateFileErr    = "already contains a file with this path"
	licenseName         = "Open Government Licence v3.0"
	licenseURL          = "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/"
)

var (
	versionRegEx = regexp.MustCompile("/version-(\\d+)/")
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
		IsPublishable: true,
		Title:         event.Title,
		FileSizeBytes: f.SizeInBytes,
		FileType:      f.MimeType,
		License:       licenseName,
		LicenseURL:    licenseURL,
	}

	version := 1
	for _, existing := range event.CurrentFiles {
		if strings.HasSuffix(existing, f.Name) {
			//file already saved so set base version to this +1
			re := versionRegEx.FindStringSubmatch(existing)
			if re != nil {
				version, _ = strconv.Atoi(re[1])
				version++
			}
		}
	}

	var attempts int
	for {
		metadata.Path, metadata.FileName = getPathAndFilename(f.Name, event.ID, version)
		err := s.backend.Upload(ctx, f.ReadCloser, metadata)
		if err == nil {
			break
		}

		if strings.Contains(err.Error(), DuplicateFileErr) {
			version++
			if attempts == maxAttempts {
				return "", fmt.Errorf("exhausted attempts to upload file %w", err)
			}
			attempts++
		} else {
			return "", err
		}
	}

	return fmt.Sprintf("%s/%s", metadata.Path, metadata.FileName), nil
}

//no leading slash: https://github.com/ONSdigital/dp-upload-service/blob/ecc6062e6fe5856385b5fafbe1105606c1a958ff/api/upload.go#L25
func getPathAndFilename(filename, id string, version int) (string, string) {
	return fmt.Sprintf("%s/%s/version-%d", uploadRootDirectory, id, version), filename
}
