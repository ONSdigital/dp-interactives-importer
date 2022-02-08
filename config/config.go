package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	AwsRegion                  string        `envconfig:"AWS_REGION"`
	UploadBucketName           string        `envconfig:"UPLOAD_BUCKET_NAME"`
	DownloadBucketName         string        `envconfig:"DOWNLOAD_BUCKET_NAME"`
	Brokers                    []string      `envconfig:"KAFKA_ADDR"`
	KafkaMaxBytes              int           `envconfig:"KAFKA_MAX_BYTES"`
	KafkaVersion               string        `envconfig:"KAFKA_VERSION"`
	KafkaSecProtocol           string        `envconfig:"KAFKA_SEC_PROTO"`
	KafkaSecCACerts            string        `envconfig:"KAFKA_SEC_CA_CERTS"`
	KafkaSecClientCert         string        `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	KafkaSecClientKey          string        `envconfig:"KAFKA_SEC_CLIENT_KEY"             json:"-"`
	KafkaSecSkipVerify         bool          `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	InteractivesReadTopic      string        `envconfig:"INTERACTIVES_READ_TOPIC"`
	InteractivesWriteTopic     string        `envconfig:"INTERACTIVES_WRITE_TOPIC"`
	InteractivesGroup          string        `envconfig:"INTERACTIVES_GROUP"`
	KafkaConsumerWorkers       int           `envconfig:"KAFKA_CONSUMER_WORKERS"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
}

var cfg *Config

func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:27400",
		AwsRegion:                  "eu-west-1",
		DownloadBucketName:         "dp-interactives-file-uploads",
		UploadBucketName:           "dp-interactives-vis-uploads",
		Brokers:                    []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		KafkaVersion:               "1.0.2",
		KafkaMaxBytes:              2000000,
		InteractivesReadTopic:      "dp-interactives-read-topic",
		InteractivesWriteTopic:     "dp-interactives-write-topic",
		KafkaConsumerWorkers:       1,
		InteractivesGroup:          "dp-interactives-importer",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
