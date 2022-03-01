package importer

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-interactives-importer/internal/client/uploadservice"
	"os"
)

const (
	uploadRootDirectory = "interactives"
)

type Upload struct {
	Title, CollectionId, Filename, Licence, LicenceUrl string
	TotalChunks, TotalSize                             int64
}

func NewUploadService(backend UploadServiceBackend, chunkSize int64) *UploadService {
	c := chunkSize
	if c <= 0 {
		c = DefaultChunkSize
	}

	return &UploadService{
		chunkSize: c,
		backend:   backend,
	}
}

//todo handling retries?
//todo handle duplicates/replace? - path name convention

type UploadService struct {
	backend   UploadServiceBackend
	chunkSize int64
	Uploads   []Upload //todo get rid of this when integrated into dp-upload-service
}

func (s *UploadService) SendFile(ctx context.Context, f *File, title, collectionId, licence, licenceUrl string) error {
	filename := getUploadFilename(f.Name, collectionId, 1)
	uploadFileFunc := func(currentChunk, totalChunks, totalSize int32, mimetype string, tmpFile *os.File) error {
		req := uploadservice.UploadJob{
			ResumableFilename:    filename,
			IsPublishable:        true, //todo isPublishable==true - assumes all files are publishable - confirm what this means (missing from swagger right now)
			CollectionId:         collectionId,
			Title:                title,
			ResumableTotalSize:   totalSize,
			ResumableType:        mimetype,
			Licence:              licence,
			LicenceUrl:           licenceUrl,
			ResumableChunkNumber: currentChunk,
			ResumableTotalChunks: totalChunks,
			File:                 tmpFile,
		}

		err := s.backend.Upload(ctx, "", req)
		if err != nil {
			return err
		}

		return nil
	}

	totalChunks, err := f.SplitAndClose(s.chunkSize, uploadFileFunc)
	if err != nil {
		return err
	}

	s.Uploads = append(s.Uploads, Upload{
		Title:        title,
		CollectionId: collectionId,
		Filename:     filename,
		Licence:      licence,
		LicenceUrl:   licenceUrl,
		TotalChunks:  totalChunks,
		TotalSize:    f.SizeInBytes,
	})

	return nil
}

//no leading slash: https://github.com/ONSdigital/dp-upload-service/blob/ecc6062e6fe5856385b5fafbe1105606c1a958ff/api/upload.go#L25
func getUploadFilename(filename, collectionId string, version int) string {
	return fmt.Sprintf("%s/%s/version-%d/%s", uploadRootDirectory, collectionId, version, filename)
}
