package consumers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"voicerx-backend/producers"

	"github.com/IBM/sarama"
)

func consumeTranscript(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0 // specify appropriate Kafka version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup([]string{brokerAddress}, "transcriptgrp", config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Consume messages from the audio topic
	handler := &TranscriptMessageHandler{}
	stopChan := make(chan struct{})

	go func() {
		defer close(stopChan)
		for {
			select {
			case <-ctx.Done():
					return
			default:
				err := consumerGroup.Consume(ctx, []string{"transcripts"}, handler)
				if err != nil {
					log.Printf("Error consuming messages: %v", err)
				}
			}
		}
	}()

	// Wait for a signal to stop the consumer
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Println("Received termination signal. Stopping transcript consumer.")
	}
}

// AudioMessageHandler is your implementation of sarama.ConsumerGroupHandler
type TranscriptMessageHandler struct{}

func (h *TranscriptMessageHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Your setup logic
	return nil
}

func (h *TranscriptMessageHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// Your cleanup logic
	return nil
}

func (h *TranscriptMessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received transcript message %v\n", string(message.Value))
		session.MarkMessage(message, "")
		err := producers.ParseTranscriptHandler(message.Value)
		if err != nil{
			log.Printf("Error while transcibing %w", err)	
		}	
	}
	return nil
}