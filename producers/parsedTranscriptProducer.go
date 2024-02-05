package producers

import(
	"fmt"
	"encoding/json"
	"context"
	"strings"

	"voicerx-backend/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehendmedical"
)

// func produceParsedTranscript(producer sarama.SyncProducer, cred string) error {
	
// 	notification := TranscriptInfo{
// 		Cred:    cred,
// 		S3Link:  "s3://voicerx-transcript-bucket/"+strings.TrimSpace(cred[0:11])+"/"+cred+".json",
// 	}

// 	notificationJSON, err := json.Marshal(notification)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal notification: %w", err)
// 	}

// 	partition := int32(0)

// 	msg := &sarama.ProducerMessage{
// 		Topic: "transcripts",
// 		Value: sarama.StringEncoder(notificationJSON),
// 		Partition: partition,
// 	}

// 	_, _, err = producer.SendMessage(msg)
// 	return err
// }

func parseTranscript(cred string, s3link string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	comprehendClient := comprehendmedical.NewFromConfig(cfg)

	transcript, err := utils.GetTranscript(strings.TrimSpace(cred[0:12])+"/"+strings.Replace(cred, "$", "_", 2))
	if err != nil{
		return fmt.Errorf("failed to load transcript from s3: %w", err)
	}

	
	doc := &comprehendmedical.DetectEntitiesInput{
		Text: aws.String(transcript),
	}

	resp, err := comprehendClient.DetectEntities(context.TODO(), doc)
	if err != nil {
		return fmt.Errorf("failed to detect entities with Comprehend Medical: %w", err)
	}

	for _, entity := range resp.Entities {
		fmt.Printf("Entity: %s, Category: %s\n", entity.Text, entity.Category)
	}

	// fmt.printf(string(resp))

	return nil
}

func ParseTranscriptHandler(message []byte) error {	
	var data TranscriptInfo
	err := json.Unmarshal([]byte(message), &data)
	if err != nil{
		return fmt.Errorf("Couldn't unmarshal JSON: %w", err)
	}

	err = parseTranscript(data.Cred, data.S3Link)
	if err != nil{
		return fmt.Errorf("Couldn't parse transcribed audio: %w", err)
	}

	// err = produceTranscript(producer, data.Cred)
	// if err != nil{
	// 	return fmt.Errorf("Couldn't produce transcription: %w", err)
	// }

	return nil
}