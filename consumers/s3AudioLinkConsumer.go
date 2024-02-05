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

const brokerAddress = "192.168.0.7:9092"

func consumeAudioLink(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0 // specify appropriate Kafka version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup([]string{brokerAddress}, "audiogrp", config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Consume messages from the audio topic
	handler := &AudioMessageHandler{}
	stopChan := make(chan struct{})

	go func() {
		defer close(stopChan)
		for {
			select {
			case <-ctx.Done():
					return
			default:
				err := consumerGroup.Consume(ctx, []string{"audio"}, handler)
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
		log.Println("Received termination signal. Stopping audio consumer.")
	}
}

// AudioMessageHandler is your implementation of sarama.ConsumerGroupHandler
type AudioMessageHandler struct{}

func (h *AudioMessageHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Your setup logic
	return nil
}

func (h *AudioMessageHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// Your cleanup logic
	return nil
}

func (h *AudioMessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received audio message %v\n", string(message.Value))
		session.MarkMessage(message, "")
		err := producers.TranscribeAudioHandler(message.Value)
		if err != nil{
			log.Printf("Error while transcibing %w", err)	
		}	
	}
	return nil
}