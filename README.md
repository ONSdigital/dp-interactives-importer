# dp-interactives-importer

Listens on a kafka topic for new interactives import events. When a new event is picked up it will:
- download zip file from an S3 bucket used for temporary storage between services
- unzip the file
- send each file to the dp-upload-service

## Getting started

* Start docker-compose environment here: https://github.com/ONSdigital/dp-interactives-compose: `docker-compose --env-file=start-backend.env`
* Run `make debug`
* `curl 'http://localhost:27400/health' | jq`
* Should see 200 with "status: OK"

## Dependencies

- Kafka messaging
- AWS S3
- Interactives API: https://github.com/ONSdigital/dp-interactives-api
- Upload Service API: https://github.com/ONSdigital/dp-upload-service

## Configuration

See [config.go](config/config.go) and https://github.com/kelseyhightower/envconfig

## License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

