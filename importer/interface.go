package importer

import (
	"context"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	s3client "github.com/ONSdigital/dp-s3"
	"io"
)

//go:generate moq -out mock/s3.go -pkg mock . S3Interface
//go:generate moq -out mock/upload_service.go -pkg mock . UploadService

type S3Interface interface {
	GetFromS3URL(rawURL string, style s3client.URLStyle) (io.ReadCloser, *int64, error)
	//todo drop below
	UploadPart(ctx context.Context, req *s3client.UploadPartRequest, payload []byte) error
	CheckPartUploaded(ctx context.Context, req *s3client.UploadPartRequest) (bool, error)
	Checker(ctx context.Context, state *health.CheckState) error
}

type UploadService interface {
	Send(file *File, title, collectionId, filename, licence, licenceUrl string) error
	Close()
}
