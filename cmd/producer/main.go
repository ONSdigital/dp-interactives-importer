package main

import (
	"context"
	"os"
	"time"

	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/schema"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/log.go/v2/log"
)

func main() {
	ctx := context.Background()
	var brokers []string
	brokers = append(brokers, "localhost:9092")

	cfg, err := config.Get()
	if err != nil {
		log.Fatal(ctx, "failed to retrieve configuration", err)
		os.Exit(1)
	}

	pChannels := kafka.CreateProducerChannels()
	pConfig := &kafka.ProducerConfig{
		MaxMessageBytes: &cfg.KafkaMaxBytes,
		KafkaVersion:    &cfg.KafkaVersion,
	}

	producer, err := kafka.NewProducer(ctx, brokers, cfg.InteractivesReadTopic, pChannels, pConfig)
	if err != nil {
		log.Fatal(ctx, "failed to create kafka producer", err)
		os.Exit(1)
	}

	event1 := importer.VisualisationUploaded{Path: "1_path_to_a_interactive"}
	sendEvent(producer, event1)
	event2 := importer.VisualisationUploaded{Path: "2_path_to_a_interactive"}
	sendEvent(producer, event2)

	time.Sleep(5 * time.Second)
	producer.Close(context.TODO())
}

func sendEvent(producer *kafka.Producer, v importer.VisualisationUploaded) {
	bytes, err := schema.VisualisationUploadedEvent.Marshal(v)
	if err != nil {
		panic(err)
	}
	producer.Channels().Output <- bytes
}
