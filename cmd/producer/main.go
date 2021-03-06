package main

import (
	"context"
	"os"
	"time"

	"github.com/ONSdigital/dp-interactives-importer/config"
	"github.com/ONSdigital/dp-interactives-importer/importer"
	"github.com/ONSdigital/dp-interactives-importer/schema"
	kafka "github.com/ONSdigital/dp-kafka/v3"
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

	pConfig := &kafka.ProducerConfig{
		BrokerAddrs:     cfg.Brokers,               // compulsory
		Topic:           cfg.InteractivesReadTopic, // compulsory
		KafkaVersion:    &cfg.KafkaVersion,
		MaxMessageBytes: &cfg.KafkaMaxBytes,
	}
	producer, err := kafka.NewProducer(ctx, pConfig)
	if err != nil {
		log.Fatal(ctx, "failed to create kafka producer", err)
		os.Exit(1)
	}

	//aws --endpoint-url=http://localhost:4566 s3 cp ~/Desktop/ovpn_configs.zip s3://testing/
	//browser at to see ls -> http://localhost:4566/testing

	//https://docs.aws.amazon.com/AmazonS3/latest/userguide/VirtualHosting.html
	//https://bucket-name.s3.Region.amazonaws.com/key-name

	event1 := importer.InteractivesUploaded{
		ID:   "52bd5e13-8dda-4593-bfe3-d4999bf3cd51", //manually created via Postman
		Path: "f5XNzqLK76cMwldF835lkCuKO34=/single-interactive.zip",
	}
	sendEvent(producer, event1)

	time.Sleep(5 * time.Second)
	producer.Close(context.TODO())
}

func sendEvent(producer *kafka.Producer, v importer.InteractivesUploaded) {
	bytes, err := schema.InteractivesUploadedEvent.Marshal(v)
	if err != nil {
		panic(err)
	}
	producer.Channels().Output <- bytes
}
