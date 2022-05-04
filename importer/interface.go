package importer

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
)

//go:generate moq -out mocks/s3.go -pkg mocks_importer . S3Interface
//go:generate moq -out mocks/upload_service_backend.go -pkg mocks_importer . UploadServiceBackend
//go:generate moq -out mocks/interactives_api.go -pkg mocks_importer . InteractivesAPIClient

type S3Interface interface {
	Get(key string) (io.ReadCloser, *int64, error)
	Checker(ctx context.Context, state *health.CheckState) error
}

type UploadServiceBackend interface {
	Upload(ctx context.Context, fileContent io.ReadCloser, metadata upload.Metadata) error
	Checker(ctx context.Context, state *health.CheckState) error
}

type InteractivesAPIClient interface {
	PatchInteractive(context.Context, string, string, string, interactives.PatchRequest) (interactives.Interactive, error)
	Checker(ctx context.Context, state *health.CheckState) error
}
