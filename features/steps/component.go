package steps_test

import (
	"archive/zip"
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
	component_test "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	mocks_importer "github.com/ONSdigital/dp-interactives-importer/importer/mocks"
	"github.com/ONSdigital/dp-interactives-importer/service"
	mocks_service "github.com/ONSdigital/dp-interactives-importer/service/mocks"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/pkg/errors"
)

var (
	//go:embed test/*
	testZips           embed.FS
	zero               = int64(0)
	ComponentTestGroup = "component-test" // kafka group name for the component test consumer
)

type Component struct {
	ErrorFeature         component_test.ErrorFeature
	serviceList          *service.ExternalServiceList
	KafkaConsumer        *kafkatest.Consumer
	S3Client             *mocks_importer.S3InterfaceMock
	UploadServiceBackend *mocks_importer.UploadServiceBackendMock
	InteractivesAPI      *mocks_importer.InteractivesAPIClientMock
	killChan             chan os.Signal
	errorChan            chan error
}

func NewInteractivesImporterComponent() (*Component, error) {
	c := &Component{
		errorChan: make(chan error),
	}

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	ctx := context.Background()

	kafkaOffset := kafka.OffsetOldest
	if c.KafkaConsumer, err = kafkatest.NewConsumer(
		ctx,
		&kafka.ConsumerGroupConfig{
			BrokerAddrs:  cfg.Brokers,
			Topic:        cfg.InteractivesReadTopic,
			GroupName:    ComponentTestGroup,
			KafkaVersion: &cfg.KafkaVersion,
			Offset:       &kafkaOffset,
		},
		nil,
	); err != nil {
		return nil, fmt.Errorf("error creating kafka consumer: %w", err)
	}
	c.KafkaConsumer.Mock.CheckerFunc = funcCheck
	c.KafkaConsumer.Mock.RegisterHandlerFunc = func(ctx context.Context, h kafka.Handler) error {
		go func() {
			for {
				select {
				case message, ok := <-c.KafkaConsumer.Mock.Channels().Upstream:
					if !ok {
						return
					}
					err := h(context.TODO(), 1, message)
					if err != nil {
						return
					}
					message.Release()
				case <-c.KafkaConsumer.Mock.Channels().Closer:
					return
				}
			}
		}()
		return nil
	}
	c.KafkaConsumer.Mock.StartFunc = func() error { return nil }

	raw, err := testZips.ReadFile("test/test_zips.zip")
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	zipFilesForTest := make(map[string]*zip.File)
	for _, f := range zipReader.File {
		zipFilesForTest[f.Name] = f
	}

	c.S3Client = &mocks_importer.S3InterfaceMock{
		CheckerFunc: funcCheck,
		GetFunc: func(key string) (io.ReadCloser, *int64, error) {
			zipFile, ok := zipFilesForTest[key]
			if !ok {
				return nil, &zero, errors.Errorf("does not exist")
			}
			rc, e := zipFile.Open()
			size := int64(zipFile.UncompressedSize64)
			return rc, &size, e
		},
	}

	c.UploadServiceBackend = &mocks_importer.UploadServiceBackendMock{
		UploadFunc: func(context.Context, io.ReadCloser, upload.Metadata) error {
			return nil
		},
	}

	c.InteractivesAPI = &mocks_importer.InteractivesAPIClientMock{
		PatchInteractiveFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, interactiveID string, req interactives.PatchRequest) (interactives.Interactive, error) {
			return interactives.Interactive{}, nil
		},
	}

	initMock := &mocks_service.InitialiserMock{
		DoGetHTTPServerFunc:            DoGetHTTPServerOk,
		DoGetHealthCheckFunc:           DoGetHealthcheckOk,
		DoGetHealthClientFunc:          DoGetHealthClient,
		DoGetKafkaConsumerFunc:         DoGetConsumer(c),
		DoGetS3ClientFunc:              DoGetS3Client(c),
		DoGetUploadServiceBackendFunc:  DoGetUploadServiceBackend(c),
		DoGetInteractivesAPIClientFunc: DoGetInteractivesAPIClient(c),
	}

	c.serviceList = service.NewServiceList(initMock)

	return c, nil
}

func (c *Component) Close() {

}

func (c *Component) Reset() {

}

func DoGetConsumer(c *Component) func(context.Context, *config.Config) (kafka.IConsumerGroup, error) {
	return func(_ context.Context, _ *config.Config) (kafka.IConsumerGroup, error) {
		return c.KafkaConsumer.Mock, nil
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

func DoGetInteractivesAPIClient(c *Component) func(ctx context.Context, cfg *config.Config) (importer.InteractivesAPIClient, error) {
	return func(_ context.Context, _ *config.Config) (importer.InteractivesAPIClient, error) {
		return c.InteractivesAPI, nil
	}
}

func funcClose(_ context.Context) error {
	return nil
}

func funcCheck(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}
