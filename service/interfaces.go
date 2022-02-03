package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	kafka "github.com/ONSdigital/dp-kafka/v2"
)

//go:generate moq -out mocks/initialiser.go -pkg mocks_service . Initialiser
//go:generate moq -out mocks/server.go -pkg mocks_service . HTTPServer
//go:generate moq -out mocks/healthcheck.go -pkg mocks_service . HealthChecker

// Initialiser defines the methods to initialise external services
type Initialiser interface {
	DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer
	DoGetKafkaProducer(ctx context.Context, cfg *config.Config) (kafka.IProducer, error)
	DoGetKafkaConsumer(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error)
	DoGetHealthClient(name, url string) *health.Client
	DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error)
	DoGetS3Client(ctx context.Context, cfg *config.Config) (importer.S3Interface, error)
}

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddCheck(name string, checker healthcheck.Checker) (err error)
}
