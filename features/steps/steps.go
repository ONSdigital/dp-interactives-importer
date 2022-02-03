package steps

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/schema"
	"github.com/ONSdigital/dp-interactives-importer/service"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-kafka/v2/kafkatest"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^these events are consumed:$`, c.theseEventsAreConsumed)
	ctx.Step(`^"([^"]*)" interactives should be uploaded to s3 successfully$`, c.theseInteractivesAreUploadedToS3)
	ctx.Step(`^a message for "([^"]*)" with "([^"]*)" is produced$`, c.aMessageForWithIsProduced)
}

func (c *Component) theseEventsAreConsumed(table *godog.Table) error {
	events, err := convertToEvents(table)
	if err != nil {
		return err
	}

	signals := registerInterrupt()

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	// run application in separate goroutine
	go func() {
		_, _ = service.Run(context.TODO(), cfg, c.serviceList, "", "", "", c.errorChan)
	}()

	// consume extracted observations
	for _, e := range events {
		if err := sendToConsumer(c.KafkaConsumer, e); err != nil {
			return err
		}
	}

	time.Sleep(300 * time.Millisecond)

	// kill application
	signals <- os.Interrupt

	return nil
}

func (c *Component) theseInteractivesAreUploadedToS3(count int) error {
	assert.Equal(&c.ErrorFeature, count, len(c.S3Client.UploadPartCalls()))
	return c.ErrorFeature.StepError()
}

func (c *Component) aMessageForWithIsProduced(id, path string) error {
	message := <-c.KafkaProducer.Channels().Output

	var unmarshalledMessage importer.VisualisationUploaded
	err := schema.VisualisationUploadedEvent.Unmarshal(message, &unmarshalledMessage)
	if err != nil {
		return err
	}

	assert.Equal(&c.ErrorFeature, id, unmarshalledMessage.ID)
	assert.Equal(&c.ErrorFeature, path, unmarshalledMessage.Path)
	return c.ErrorFeature.StepError()
}

func convertToEvents(table *godog.Table) ([]*importer.VisualisationUploaded, error) {
	assist := assistdog.NewDefault()
	events, err := assist.CreateSlice(&importer.VisualisationUploaded{}, table)
	if err != nil {
		return nil, err
	}
	return events.([]*importer.VisualisationUploaded), nil
}

func registerInterrupt() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return signals
}

func sendToConsumer(kafka kafka.IConsumerGroup, e *importer.VisualisationUploaded) error {
	bytes, err := schema.VisualisationUploadedEvent.Marshal(e)
	if err != nil {
		return err
	}

	kafka.Channels().Upstream <- kafkatest.NewMessage(bytes, 0)
	return nil
}
