package utils

import(
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetTranscript(objectKey string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS configuration: %w", err)

	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Set up input parameters for S3 getObject
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String("voicerx-transcription-bucket"),
		Key:    aws.String(objectKey),
	}

	// Get the JSON file from S3
	getObjectOutput, err := s3Client.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		return "", fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer getObjectOutput.Body.Close()

	// Decode the JSON file
	// var jsonData YourJSONStruct
	// err = json.NewDecoder(getObjectOutput.Body).Decode(&jsonData)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to decode JSON:: %w", err)
	// }

	bodyContent, err := ioutil.ReadAll(getObjectOutput.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read object body: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(bodyContent), &result)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling transcript json:", err)
	}

	// Extract the "transcript" attribute
	transcripts := result["results"].(map[string]interface{})["transcripts"].([]interface{})
	transcriptObj := transcripts[0].(map[string]interface{})
	transcript := transcriptObj["transcript"].(string)

	// Use the transcript variable as needed
	return transcript, nil
}