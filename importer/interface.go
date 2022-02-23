package importer

import (
	"context"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/internal/client/uploadservice"
	s3client "github.com/ONSdigital/dp-s3"
	"io"
)

//go:generate moq -out mocks/s3.go -pkg mocks_importer . S3Interface
//go:generate moq -out mocks/upload_service_backend.go -pkg mocks_importer . UploadServiceBackend

type S3Interface interface {
	Get(key string) (io.ReadCloser, *int64, error)
	//todo drop below
	UploadPart(ctx context.Context, req *s3client.UploadPartRequest, payload []byte) error
	CheckPartUploaded(ctx context.Context, req *s3client.UploadPartRequest) (bool, error)
	Checker(ctx context.Context, state *health.CheckState) error
}

type UploadServiceBackend interface {
	Upload(ctx context.Context, serviceToken string, job uploadservice.UploadJob) error
	Checker(ctx context.Context, state *health.CheckState) error
}
