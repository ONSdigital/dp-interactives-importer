package steps_test

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/health"
	component_test "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	mocks_importer "github.com/ONSdigital/dp-interactives-importer/importer/mocks"
	"github.com/ONSdigital/dp-interactives-importer/internal/client/uploadservice"
	"github.com/ONSdigital/dp-interactives-importer/internal/test"
	"github.com/ONSdigital/dp-interactives-importer/service"
	mocks_service "github.com/ONSdigital/dp-interactives-importer/service/mocks"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-kafka/v2/kafkatest"
	"io"
	"net/http"
	"os"
)

type Component struct {
	ErrorFeature         component_test.ErrorFeature
	serviceList          *service.ExternalServiceList
	KafkaConsumer        kafka.IConsumerGroup
	S3Client             *mocks_importer.S3InterfaceMock
	UploadServiceBackend *mocks_importer.UploadServiceBackendMock
	killChan             chan os.Signal
	errorChan            chan error
	testZipArchive       *os.File
}

func NewInteractivesImporterComponent() *Component {
	c := &Component{
		errorChan: make(chan error),
	}

	archiveName, _ := test.CreateTestZip("root.css", "root.html", "root.js")
	c.testZipArchive, _ = os.Open(archiveName)

	consumer := kafkatest.NewMessageConsumer(false)
	consumer.CheckerFunc = funcCheck
	c.KafkaConsumer = consumer

	c.S3Client = &mocks_importer.S3InterfaceMock{
		CheckerFunc: funcCheck,
		GetFunc: func(key string) (io.ReadCloser, *int64, error) {
			stat, _ := c.testZipArchive.Stat()
			size := stat.Size()
			return c.testZipArchive, &size, nil
		},
	}

	c.UploadServiceBackend = &mocks_importer.UploadServiceBackendMock{
		CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
			return nil
		},
		UploadFunc: func(_ context.Context, _ string, _ uploadservice.UploadJob) error {
			return nil
		},
	}

	initMock := &mocks_service.InitialiserMock{
		DoGetHTTPServerFunc:           DoGetHTTPServerOk,
		DoGetHealthCheckFunc:          DoGetHealthcheckOk,
		DoGetHealthClientFunc:         DoGetHealthClient,
		DoGetKafkaConsumerFunc:        DoGetConsumer(c),
		DoGetS3ClientFunc:             DoGetS3Client(c),
		DoGetUploadServiceBackendFunc: DoGetUploadServiceBackend(c),
	}

	c.serviceList = service.NewServiceList(initMock)

	return c
}

func (c *Component) Close() {
	_ = os.Remove(c.testZipArchive.Name())
}

func (c *Component) Reset() {

}

func DoGetConsumer(c *Component) func(context.Context, *config.Config) (kafka.IConsumerGroup, error) {
	return func(_ context.Context, _ *config.Config) (kafka.IConsumerGroup, error) {
		return c.KafkaConsumer, nil
	}
}

func DoGetHealthcheckOk(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
	return &mocks_service.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func DoGetHealthClient(name, url string) *health.Client {
	return &health.Client{}
}

func DoGetS3Client(c *Component) func(ctx context.Context, cfg *config.Config) (importer.S3Interface, error) {
	return func(_ context.Context, _ *config.Config) (importer.S3Interface, error) {
		return c.S3Client, nil
	}
}

func DoGetHTTPServerOk(bindAddr string, router http.Handler) service.HTTPServer {
	return &mocks_service.HTTPServerMock{
		ListenAndServeFunc: func() error {
			return nil
		},
	}
}

func DoGetUploadServiceBackend(c *Component) func(ctx context.Context, cfg *config.Config) (importer.UploadServiceBackend, error) {
	return func(_ context.Context, _ *config.Config) (importer.UploadServiceBackend, error) {
		return c.UploadServiceBackend, nil
	}
}

func funcClose(_ context.Context) error {
	return nil
}

func funcCheck(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}
