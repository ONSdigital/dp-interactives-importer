package importer

import (
	"context"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	s3client "github.com/ONSdigital/dp-s3"
)

//go:generate moq -out mocks/s3.go -pkg mocks_importer . S3Interface

type S3Interface interface {
	UploadPart(ctx context.Context, req *s3client.UploadPartRequest, payload []byte) error
	CheckPartUploaded(ctx context.Context, req *s3client.UploadPartRequest) (bool, error)
	Checker(ctx context.Context, state *health.CheckState) error
}
