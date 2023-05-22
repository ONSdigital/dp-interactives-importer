package steps_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/schema"
	"github.com/ONSdigital/dp-interactives-importer/service"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^these events are consumed:$`, c.theseEventsAreConsumed)
	ctx.Step(`^"([^"]*)" interactives should be downloaded from s3 successfully$`, c.theseInteractivesAreDownloadedFromS3)
	ctx.Step(`^"([^"]*)" interactives should be uploaded via the upload service$`, c.interactivesShouldBeUploadedViaTheUploadService)
	ctx.Step(`^"([^"]*)" interactive should be successfully updated via the interactives API$`, c.interactiveShouldBeSuccessfullyUpdatedViaTheInteractivesAPI)
	ctx.Step(`^"([^"]*)" interactive should be updated as a failure via the interactives API$`, c.interactiveShouldBeUpdatedAsAFailureViaTheInteractivesAPI)
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
		if err := sendToConsumer(c.kafkaConsumer.Mock, e); err != nil {
			return err
		}
	}

	time.Sleep(300 * time.Millisecond)

	// kill application
	signals <- os.Interrupt

	return nil
}

func convertToEvents(table *godog.Table) ([]*importer.InteractivesUploaded, error) {
	assist := assistdog.NewDefault()
	events, err := assist.CreateSlice(&importer.InteractivesUploaded{}, table)
	if err != nil {
		return nil, err
	}
	return events.([]*importer.InteractivesUploaded), nil
}

func registerInterrupt() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return signals
}

func sendToConsumer(kafka kafka.IConsumerGroup, e *importer.InteractivesUploaded) error {
	bytes, err := schema.InteractivesUploadedEvent.Marshal(e)
	if err != nil {
		return err
	}

	msg, err := kafkatest.NewMessage(bytes, 0)
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}
	kafka.Channels().Upstream <- msg
	return nil
}

func (c *Component) theseInteractivesAreDownloadedFromS3(count int) error {
	assert.Equal(&c.ErrorFeature, count, len(c.S3Client.GetCalls()))
	return c.ErrorFeature.StepError()
}

func (c *Component) interactivesShouldBeUploadedViaTheUploadService(count int) error {
	assert.Equal(&c.ErrorFeature, count, len(c.UploadServiceBackend.UploadCalls()))
	return c.ErrorFeature.StepError()
}

func (c *Component) interactiveShouldBeSuccessfullyUpdatedViaTheInteractivesAPI(id string) error {
	assert.Equal(&c.ErrorFeature, 1, len(c.InteractivesAPI.PatchInteractiveCalls()))
	firstCall := c.InteractivesAPI.PatchInteractiveCalls()[0]
	dir := firstCall.PatchRequest.Interactive.Archive.UploadRootDirectory
	isUploadRootDirWithExpectedPrefix := strings.HasPrefix(dir, "interactives/")
	assert.True(&c.ErrorFeature, firstCall.PatchRequest.Interactive.Archive.ImportSuccessful)
	assert.Equal(&c.ErrorFeature, id, firstCall.S3)
	assert.Equal(&c.ErrorFeature, id, firstCall.PatchRequest.Interactive.ID)
	assert.True(&c.ErrorFeature, isUploadRootDirWithExpectedPrefix)
	return c.ErrorFeature.StepError()
}

func (c *Component) interactiveShouldBeUpdatedAsAFailureViaTheInteractivesAPI(id string) error {
	assert.Equal(&c.ErrorFeature, 1, len(c.InteractivesAPI.PatchInteractiveCalls()))
	firstCall := c.InteractivesAPI.PatchInteractiveCalls()[0]
	assert.Equal(&c.ErrorFeature, id, firstCall.S3)
	assert.Equal(&c.ErrorFeature, id, firstCall.PatchRequest.Interactive.ID)
	assert.NotEmpty(&c.ErrorFeature, firstCall.PatchRequest.Interactive.Archive.ImportMessage)
	assert.False(&c.ErrorFeature, firstCall.PatchRequest.Interactive.Archive.ImportSuccessful)
	return c.ErrorFeature.StepError()
}
