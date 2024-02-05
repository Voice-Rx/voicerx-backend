package producers

import (
	"encoding/json"
	"fmt"
	"log"	

	"github.com/IBM/sarama"
	"voicerx-backend/models"
)

// ============== KAFKA RELATED FUNCTIONS ==============
func uploadLink(producer sarama.SyncProducer, cred string, link string) error {
	notification := models.AudioLink{
		Cred:    cred,
		S3Link:  link,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	partition := int32(0)

	msg := &sarama.ProducerMessage{
		Topic: "audio",
		Value: sarama.StringEncoder(notificationJSON),
		Partition: partition,
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func AddAudioKafka(cred string, s3Link string) error {
	producer, err := SetupProducer()

	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()

	err = uploadLink(producer, cred, s3Link)

	if err != nil {
		fmt.Errorf("Uploading link to Kafka failed: %w", err)
		return err
	}

	return nil
}