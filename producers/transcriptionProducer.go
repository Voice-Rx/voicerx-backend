package producers

import(
	"fmt"
	"log"
	"strings"
	"encoding/json"
	"context"

	"voicerx-backend/utils"

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
)

type AudioInfo struct {
	Cred string `json:"cred"`
	S3Link string `json:s3link"`
}

type TranscriptInfo struct {
	Cred string `json:"cred"`
	S3Link string `json:s3link"`
}

func produceTranscript(transcribeClient *transcribe.Client, producer sarama.SyncProducer, cred string) error {
	
	notification := TranscriptInfo{
		Cred:    cred,
		S3Link:  "s3://voicerx-transcript-bucket/"+strings.TrimSpace(cred[0:12])+"/"+strings.Replace(cred, "$", "_", 2),
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	partition := int32(0)

	err = utils.WaitForTranscriptionCompletion(transcribeClient, strings.Replace(cred, "$", "_", 2))
	if err != nil{
		return fmt.Errorf("Transcription job couldn't be completed: %w", err)
	}

	log.Printf("Producing message")

	msg := &sarama.ProducerMessage{
		Topic: "transcripts",
		Value: sarama.StringEncoder(notificationJSON),
		Partition: partition,
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func transcribeAudio(cred string, s3Link string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	transcribeClient := transcribe.NewFromConfig(cfg)
	
	_, err = transcribeClient.StartMedicalTranscriptionJob(context.TODO(), &transcribe.StartMedicalTranscriptionJobInput{
		MedicalTranscriptionJobName: aws.String(strings.Replace(cred, "$", "_", 2)),
		LanguageCode:                "en-US",
		MediaFormat:                 "mp3", // Change based on your file format
		Media: &types.Media {
			MediaFileUri:            aws.String(s3Link),
		},
		OutputBucketName:            aws.String("voicerx-transcription-bucket"),
		OutputKey:                   aws.String(strings.TrimSpace(cred[0:12])+"/"+strings.Replace(cred, "$", "_", 2)),
		Specialty:                   "PRIMARYCARE",
		Type:                        "DICTATION",
	})

	if err != nil {
		return fmt.Errorf("failed to start transcription job: %w", err)
	}

	producer, err := SetupProducer()

	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()

	err = produceTranscript(transcribeClient, producer, cred)
	if err != nil{
		return fmt.Errorf("Couldn't produce transcription: %w", err)
	}

	return nil
}

func TranscribeAudioHandler(message []byte) error {
	var data AudioInfo
	err := json.Unmarshal([]byte(message), &data)
	if err != nil{
		return fmt.Errorf("Couldn't unmarshal JSON: %w", err)
	}

	err = transcribeAudio(data.Cred, data.S3Link)
	if err != nil{
		return fmt.Errorf("Couldn't transcribe or produce audio: %w", err)
	}
	
	return nil
}