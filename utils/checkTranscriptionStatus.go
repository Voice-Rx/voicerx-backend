package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
)

func WaitForTranscriptionCompletion(transcribeClient *transcribe.Client, jobName string) error {
	for {
		log.Printf("Checking if transcription is complete")
		// Get information about the transcription job
		resp, err := transcribeClient.GetMedicalTranscriptionJob(context.TODO(), &transcribe.GetMedicalTranscriptionJobInput{
			MedicalTranscriptionJobName: aws.String(jobName),
		})
		if err != nil {
			return fmt.Errorf("failed to get transcription job status: %w", err)
		}

		// Check the status of the transcription job
		switch resp.MedicalTranscriptionJob.TranscriptionJobStatus {
		case "COMPLETED":
			// Transcription is completed
			log.Printf("Transcription complete")
			return nil
		case "FAILED":
			// Transcription has failed
			return fmt.Errorf("transcription job failed")
		case "IN_PROGRESS":
			// Transcription is still in progress, wait and check again
			time.Sleep(5 * time.Second) // You can adjust the interval as needed
		case "QUEUED":
			// Transcription is still in progress, wait and check again
			time.Sleep(5 * time.Second) // You can adjust the interval as needed

		default:
			// Unknown status, consider handling appropriately
			return fmt.Errorf("unknown transcription job status: %s", resp.MedicalTranscriptionJob.TranscriptionJobStatus)
		}
	}
}