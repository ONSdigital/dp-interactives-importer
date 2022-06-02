package service

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-api-clients-go/v2/upload"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	dphttp "github.com/ONSdigital/dp-net/http"
	dps3 "github.com/ONSdigital/dp-s3"
)

type ExternalServiceList struct {
	HealthCheck          bool
	KafkaConsumer        bool
	S3Client             bool
	UploadServiceBackend bool
	InteractivesApi      bool
	Init                 Initialiser
}

func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck:          false,
		KafkaConsumer:        false,
		S3Client:             false,
		UploadServiceBackend: false,
		InteractivesApi:      false,
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

// GetInteractivesAPIClient creates an interactives api client and sets the InteractivesApi flag to true
func (e *ExternalServiceList) GetInteractivesAPIClient(ctx context.Context, cfg *config.Config) (importer.InteractivesAPIClient, error) {
	client, err := e.Init.DoGetInteractivesAPIClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e.InteractivesApi = true
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

	cgConfig := &kafka.ConsumerGroupConfig{
		Offset:       &kafkaOffset,
		BrokerAddrs:  cfg.Brokers,               // compulsory
		Topic:        cfg.InteractivesReadTopic, // compulsory
		GroupName:    cfg.InteractivesGroup,     // compulsory
		KafkaVersion: &cfg.KafkaVersion,
		NumWorkers:   &cfg.KafkaConsumerWorkers,
	}
	if cfg.KafkaSecProtocol == "TLS" {
		cgConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.KafkaSecCACerts,
			cfg.KafkaSecClientCert,
			cfg.KafkaSecClientKey,
			cfg.KafkaSecSkipVerify,
		)
	}

	return kafka.NewConsumerGroup(ctx, cgConfig)
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
	//apiClient := &mocks_importer.UploadServiceBackendMock{
	//	UploadFunc: func(context.Context, io.ReadCloser, upload.Metadata) error {
	//		return nil
	//	},
	//	CheckerFunc: func(_ context.Context, _ *healthcheck.CheckState) error {
	//		return nil
	//	},
	//}

	apiClient := upload.NewAPIClient(cfg.UploadAPIURL, cfg.ServiceAuthToken)
	return apiClient, nil
}

// DoGetInteractivesApiClient returns an interactives api client
func (e *Init) DoGetInteractivesAPIClient(ctx context.Context, cfg *config.Config) (importer.InteractivesAPIClient, error) {
	apiClient := interactives.NewAPIClient(cfg.InteractivesAPIURL, "v1")
	return apiClient, nil
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
