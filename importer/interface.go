package importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	s3client "github.com/ONSdigital/dp-s3"
	"io"
)

//go:generate moq -out mocks/s3.go -pkg mocks_importer . S3Interface
//go:generate moq -out mocks/upload_service_backend.go -pkg mocks_importer . UploadServiceBackend
//go:generate moq -out mocks/interactives_api.go -pkg mocks_importer . InteractivesAPIClient

type S3Interface interface {
	Get(key string) (io.ReadCloser, *int64, error)
	//todo drop below
	UploadPart(ctx context.Context, req *s3client.UploadPartRequest, payload []byte) error
	CheckPartUploaded(ctx context.Context, req *s3client.UploadPartRequest) (bool, error)
	Checker(ctx context.Context, state *health.CheckState) error
}

type UploadServiceBackend interface {
	Upload(ctx context.Context, fileContent io.ReadCloser, metadata upload.Metadata) error
	Checker(ctx context.Context, state *health.CheckState) error
}

type InteractivesAPIClient interface {
	PutInteractive(ctx context.Context, userAuthToken, serviceAuthToken, interactiveID string, update interactives.InteractiveUpdate) error
	Checker(ctx context.Context, state *health.CheckState) error
}
