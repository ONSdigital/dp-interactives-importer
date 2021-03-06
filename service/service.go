package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the dp-upload-service API
type Service struct {
	config        *config.Config
	serviceList   *ExternalServiceList
	healthCheck   HealthChecker
	kafkaConsumer kafka.IConsumerGroup
}

func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Info(ctx, "running service")

	r := mux.NewRouter()
	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	// Setup downstream dependencies
	consumer, err := serviceList.GetKafkaConsumer(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka consumer", err)
		return nil, err
	}

	s3Client, err := serviceList.GetS3Client(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise S3 client for uploaded bucket", err)
		return nil, err
	}

	uploadServiceBackend, err := serviceList.GetUploadServiceBackend(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise upload service", err)
		return nil, err
	}
	uploadService := importer.NewUploadService(uploadServiceBackend)

	interactivesAPIClient, err := serviceList.GetInteractivesAPIClient(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise interactives api", err)
		return nil, err
	}

	// Event Handler for Kafka Consumer
	handler := &importer.InteractivesUploadedHandler{
		Cfg:                   cfg,
		S3:                    s3Client,
		UploadService:         uploadService,
		InteractivesAPIClient: interactivesAPIClient,
	}
	err = consumer.RegisterHandler(ctx, handler.Handle)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka consumer", err)
		return nil, err
	}
	err = consumer.Start()
	if err != nil {
		log.Fatal(ctx, "failed to start kafka consumer", err)
		return nil, err
	}

	//heathcheck - start
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}
	err = registerCheckers(ctx, cfg, hc, consumer, s3Client, uploadServiceBackend, interactivesAPIClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").Methods(http.MethodGet).HandlerFunc(hc.Handler)
	hc.Start(ctx)
	//healthcheck - end

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		config:        cfg,
		serviceList:   serviceList,
		healthCheck:   hc,
		kafkaConsumer: consumer,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var gracefulShutdown bool

	go func() {
		defer cancel()
		var hasShutdownError bool

		// stop healthcheck first, as it depends on everything else
		if svc.serviceList.HealthCheck {
			svc.healthCheck.Stop()
		}

		if svc.serviceList.KafkaConsumer {
			if err := svc.kafkaConsumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing Kafka consumer", err)
				hasShutdownError = true
			}
		}

		if !hasShutdownError {
			gracefulShutdown = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	if !gracefulShutdown {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context,
	cfg *config.Config,
	hc HealthChecker,
	consumer kafka.IConsumerGroup,
	s3 importer.S3Interface,
	uploadServiceBackend importer.UploadServiceBackend,
	interactivesAPIClient importer.InteractivesAPIClient) (err error) {

	hasErrors := false

	if err = hc.AddCheck("Kafka consumer", consumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for kafka consumer", err, log.Data{"group": cfg.InteractivesGroup, "topic": cfg.InteractivesReadTopic})
	}

	if err = hc.AddCheck("S3 bucket", s3.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for s3", err)
	}

	if err = hc.AddCheck("Upload API", uploadServiceBackend.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add Upload API health checker", err)
	}

	if err = hc.AddCheck("Interactives API", interactivesAPIClient.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add Interactives API health checker", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
