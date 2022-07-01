package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	UploadAPIURL               string        `envconfig:"UPLOAD_API_URL"`
	InteractivesAPIURL         string        `envconfig:"INTERACTIVES_API_URL"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN" json:"-"`
	AwsEndpoint                string        `envconfig:"AWS_ENDPOINT"`
	AwsRegion                  string        `envconfig:"AWS_REGION"`
	DownloadBucketName         string        `envconfig:"DOWNLOAD_BUCKET_NAME"`
	Brokers                    []string      `envconfig:"KAFKA_ADDR"`
	KafkaMaxBytes              int           `envconfig:"KAFKA_MAX_BYTES"`
	KafkaVersion               string        `envconfig:"KAFKA_VERSION"`
	KafkaSecProtocol           string        `envconfig:"KAFKA_SEC_PROTO"`
	KafkaSecCACerts            string        `envconfig:"KAFKA_SEC_CA_CERTS"`
	KafkaSecClientCert         string        `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	KafkaSecClientKey          string        `envconfig:"KAFKA_SEC_CLIENT_KEY" json:"-"`
	KafkaSecSkipVerify         bool          `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	InteractivesReadTopic      string        `envconfig:"INTERACTIVES_READ_TOPIC"`
	InteractivesGroup          string        `envconfig:"INTERACTIVES_GROUP"`
	KafkaConsumerWorkers       int           `envconfig:"KAFKA_CONSUMER_WORKERS"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	BatchSize                  int           `envconfig:"BATCH_SIZE"`
}

var cfg *Config

func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":27400",
		UploadAPIURL:               "http://localhost:25100",
		InteractivesAPIURL:         "http://localhost:27500",
		AwsRegion:                  "eu-west-1",
		DownloadBucketName:         "dp-interactives-file-uploads",
		Brokers:                    []string{"localhost:9093"},
		KafkaVersion:               "1.0.2",
		KafkaMaxBytes:              2000000,
		InteractivesReadTopic:      "interactives-import",
		KafkaConsumerWorkers:       1,
		InteractivesGroup:          "dp-interactives-importer",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		BatchSize:                  1000,
	}

	return cfg, envconfig.Process("", cfg)
}
