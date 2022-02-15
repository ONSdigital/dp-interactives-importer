# dp-interactives-importer

Reads interactive event (visualisations), parse the content and push visualisations to S3 bucket.

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                           | Description                                           |
|------------------------------|-----------------------------------|-------------------------------------------------------|
| BIND_ADDR                    | 27500                             | The host and port to bind to
| AWS_REGION                   | eu-west-1                         | The AWS region
| UPLOAD_BUCKET_NAME           | dp-interactives-visual-uploads    | Name of the S3 bucket to upload the processed visualisations
| DOWNLOAD_BUCKET_NAME         | dp-interactives-file-uploads      | Name of the S3 bucket to fetch uploaded visualisations
| KAFKA_ADDR                   | `localhost:9092`                  | The address of Kafka brokers (comma-separated values)
| KAFKA_VERSION                | `1.0.2`                           | The version of Kafka
| KAFKA_MAX_BYTES              | 2000000                           | Maximum number of bytes in a kafka message
| KAFKA_SEC_PROTO              | _unset_                           | if set to `TLS`, kafka connections will use TLS [1]
| KAFKA_SEC_CLIENT_KEY         | _unset_                           | PEM for the client key [1]
| KAFKA_SEC_CLIENT_CERT        | _unset_                           | PEM for the client certificate [1]
| KAFKA_SEC_CA_CERTS           | _unset_                           | CA cert chain for the server cert [1]
| KAFKA_SEC_SKIP_VERIFY        | false                             | ignores server certificate issues if `true` [1]
| KAFKA_CONSUMER_WORKERS       | 1                                 | The maximum number of parallel kafka consumers
| INTERACTIVES_GROUP           | dp-interactives-importer          | The consumer group this application uses

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
