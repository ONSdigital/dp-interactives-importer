package importer

import (
	"context"
	openapi "github.com/ONSdigital/dp-interactives-importer/internal/client/dp-upload-service/go"
	"os"
)

const (
	defaultChunkSize = 1024
)

type Upload struct {
	Title, CollectionId, Filename, Licence, LicenceUrl string
	TotalChunks, TotalSize                             int64
}

func NewApiUploadService(apiUrl string, chunkSize int64) (*ApiUploadService, error) {
	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{
		{URL: apiUrl},
	}

	c := chunkSize
	if c == 0 {
		c = defaultChunkSize
	}

	return &ApiUploadService{
		chunkSize: c,
		client:    openapi.NewAPIClient(cfg),
	}, nil
}

//todo handling retries?
//todo handle duplicates/replace? - path name convention

type ApiUploadService struct {
	client    *openapi.APIClient
	chunkSize int64
	Uploads   []Upload
}

func (s *ApiUploadService) Send(ctx context.Context, f *File, title, collectionId, filename, licence, licenceUrl string) error {
	uploadFileFunc := func(currentChunk, totalChunks, totalSize int, mimetype string, tmpFile *os.File) error {
		req := s.client.UploadFileAndProvideMetadataApi.V1UploadPost(ctx)
		req.Title(title)
		req.File(tmpFile)
		req.CollectionId(collectionId)
		req.Licence(licence)
		req.LicenceUrl(licenceUrl)
		req.IsPublishable(true) //todo isPublishable==true - assumes all files are publishable - confirm what this means (missing from swagger right now)
		req.ResumableFilename(filename)
		req.ResumableChunkNumber(int32(currentChunk))
		req.ResumableTotalChunks(int32(totalChunks))
		req.ResumableTotalSize(int32(totalSize))
		req.ResumableType(mimetype)
		_, err := s.client.UploadFileAndProvideMetadataApi.V1UploadPostExecute(req)
		if err != nil {
			return err
		}

		return nil
	}

	totalChunks, totalSize, err := f.SplitAndClose(s.chunkSize, uploadFileFunc)
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
		TotalSize:    totalSize,
	})

	return nil
}
