package service

import (
	"context"
	"github.com/ONSdigital/dp-interactives-importer/internal/client/uploadservice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	dphttp "github.com/ONSdigital/dp-net/http"
	dps3 "github.com/ONSdigital/dp-s3"
)

type ExternalServiceList struct {
	HealthCheck          bool
	KafkaConsumer        bool
	S3Client             bool
	UploadServiceBackend bool
	Init                 Initialiser
}

func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck:          false,
		KafkaConsumer:        false,
		S3Client:             false,
		UploadServiceBackend: false,
		Init:                 initialiser,
	}
}

type Init struct{}

// GetHTTPServer creates an http server
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetKafkaConsumer creates a Kafka consumer and sets the consumer flag to true
func (e *ExternalServiceList) GetKafkaConsumer(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
	consumer, err := e.Init.DoGetKafkaConsumer(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e.KafkaConsumer = true
	return consumer, nil
}

// GetS3Client creates a S3 client and sets the S3Client flag to true
func (e *ExternalServiceList) GetS3Client(ctx context.Context, cfg *config.Config) (importer.S3Interface, error) {
	s3, err := e.Init.DoGetS3Client(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e.S3Client = true
	return s3, nil
}

// GetUploadServiceBackend creates upload service backend and sets the UploadServiceBackend flag to true
func (e *ExternalServiceList) GetUploadServiceBackend(ctx context.Context, cfg *config.Config) (importer.UploadServiceBackend, error) {
	client, err := e.Init.DoGetUploadServiceBackend(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e.UploadServiceBackend = true
	return client, nil
}

// GetHealthClient returns a healthclient for the provided URL
func (e *ExternalServiceList) GetHealthClient(name, url string) *health.Client {
	return e.Init.DoGetHealthClient(name, url)
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// -- Implementations

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// DoGetKafkaConsumer returns a Kafka Consumer group
func (e *Init) DoGetKafkaConsumer(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
	kafkaOffset := kafka.OffsetOldest

	cConfig := &kafka.ConsumerGroupConfig{
		Offset:       &kafkaOffset,
		KafkaVersion: &cfg.KafkaVersion,
	}
	if cfg.KafkaSecProtocol == "TLS" {
		cConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.KafkaSecCACerts,
			cfg.KafkaSecClientCert,
			cfg.KafkaSecClientKey,
			cfg.KafkaSecSkipVerify,
		)
	}

	cgChannels := kafka.CreateConsumerGroupChannels(cfg.KafkaConsumerWorkers)

	return kafka.NewConsumerGroup(
		ctx,
		cfg.Brokers,
		cfg.InteractivesReadTopic,
		cfg.InteractivesGroup,
		cgChannels,
		cConfig,
	)
}

// DoGetS3Uploaded returns a S3Client
func (e *Init) DoGetS3Client(ctx context.Context, cfg *config.Config) (importer.S3Interface, error) {
	if cfg.AwsEndpoint != "" {
		//for local development only - set env var to initialise
		s, err := session.NewSession(&aws.Config{
			Endpoint:         aws.String(cfg.AwsEndpoint),
			Region:           aws.String(cfg.AwsRegion),
			S3ForcePathStyle: aws.Bool(true),
			Credentials:      credentials.NewStaticCredentials("n/a", "n/a", ""),
		})

		if err != nil {
			return nil, err
		}

		return dps3.NewClientWithSession(cfg.DownloadBucketName, s), nil
	}

	s3Client, err := dps3.NewClient(cfg.AwsRegion, cfg.DownloadBucketName)
	if err != nil {
		return nil, err
	}
	return s3Client, nil
}

// DoGetUploadServiceBackend returns an upload service backend
func (e *Init) DoGetUploadServiceBackend(ctx context.Context, cfg *config.Config) (importer.UploadServiceBackend, error) {
	uploadSvcBackend := uploadservice.New(cfg.ApiRouterUrl)

	//mocked - i got working e2e locally but had to make significant code changes in dp-upload-service
	//uploadSvcBackend := &mocks_importer.UploadServiceBackendMock{
	//	CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
	//		return state.Update(healthcheck.StatusOK, "mocked upload service backend healthy", 0)
	//	},
	//	UploadFunc: func(ctx context.Context, _ string, job uploadservice.UploadJob) error {
	//		logData := log.Data{"job": job}
	//		log.Info(ctx, "file uploaded", logData)
	//		return nil
	//	},
	//}

	return uploadSvcBackend, nil
}

// DoGetHealthClient creates a new Health Client for the provided name and url
func (e *Init) DoGetHealthClient(name, url string) *health.Client {
	return health.NewClient(name, url)
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}
