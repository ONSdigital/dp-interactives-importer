package steps

import (
	"context"
	"net/http"
	"os"

	"github.com/ONSdigital/dp-api-clients-go/health"
	component_test "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	mocks_importer "github.com/ONSdigital/dp-interactives-importer/importer/mocks"
	"github.com/ONSdigital/dp-interactives-importer/service"
	mocks_service "github.com/ONSdigital/dp-interactives-importer/service/mocks"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-kafka/v2/kafkatest"
	s3client "github.com/ONSdigital/dp-s3"
)

type Component struct {
	ErrorFeature  component_test.ErrorFeature
	serviceList   *service.ExternalServiceList
	KafkaConsumer kafka.IConsumerGroup
	KafkaProducer kafka.IProducer
	S3Client      *mocks_importer.S3InterfaceMock
	killChan      chan os.Signal
	errorChan     chan error
}

func NewInteractivesImporterComponent() *Component {
	c := &Component{
		errorChan: make(chan error),
	}
	//kafka
	consumer := kafkatest.NewMessageConsumer(false)
	consumer.CheckerFunc = funcCheck
	c.KafkaConsumer = consumer
	channels := &kafka.ProducerChannels{
		Output: make(chan []byte),
	}
	c.KafkaProducer = &kafkatest.IProducerMock{
		ChannelsFunc: func() *kafka.ProducerChannels {
			return channels
		},
		CloseFunc:   funcClose,
		CheckerFunc: funcCheck,
	}
	//s3
	c.S3Client = &mocks_importer.S3InterfaceMock{
		CheckerFunc: funcCheck,
		UploadPartFunc: func(_ context.Context, _ *s3client.UploadPartRequest, _ []byte) error {
			return nil
		},
	}

	initMock := &mocks_service.InitialiserMock{
		DoGetHTTPServerFunc:    DoGetHTTPServerOk,
		DoGetHealthCheckFunc:   DoGetHealthcheckOk,
		DoGetHealthClientFunc:  DoGetHealthClient,
		DoGetKafkaConsumerFunc: DoGetConsumer(c),
		DoGetS3ClientFunc:      DoGetS3Client(c),
	}

	c.serviceList = service.NewServiceList(initMock)

	return c
}

func (c *Component) Close() {

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

func funcClose(_ context.Context) error {
	return nil
}

func funcCheck(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}
