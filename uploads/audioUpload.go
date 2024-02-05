package uploads

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"context"
	"time"
	"encoding/json"
	"strings"

	"voicerx-backend/producers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func uploadToS3(file *multipart.FileHeader, customFilename string) (string, error) {
	// Create an S3 client using AWS credentials
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// Open the uploaded file
	srcFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}

	defer srcFile.Close()

	// Generate the object key (filename) for S3
	objectKey := customFilename + filepath.Ext(file.Filename)

	// Upload the file to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("voicerx-audio-bucket"),
		Key:    aws.String(objectKey),
		Body:   srcFile,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate the S3 link
	s3Link := fmt.Sprintf("s3://voicerx-audio-bucket/%s", objectKey)

	return s3Link, nil
}

func HandleAudioUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for the entire request
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	formData := r.MultipartForm
	now := time.Now()
	patientId := formData.Value["patientId"][0]
	doctorId := formData.Value["doctorId"][0]

	// Get the file from the form data
	file, fileHeader, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Unable to get audio file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	patientId = strings.Replace(patientId, "/", "-", 2)
	doctorId = strings.Replace(doctorId, "/", "-", 2)

	// Configure the filename
	customFilename := patientId + "$" +doctorId + "$" + now.Format("2006-01-02")

	// Upload the file to S3 and get the S3 link
	s3Link, err := uploadToS3(fileHeader, customFilename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload file to S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the S3 link
	apiResp := struct {
		Status string `json:"status"`
		PatientId string `json:"patientId"`
		DoctorId string `json:"doctorId"`
		Date string `json:"date"`
		S3Link string `json:"s3link"`
	}{
		Status: "Success",
		PatientId: patientId,
		DoctorId: doctorId,
		Date: now.Format("2006-01-02"),
		S3Link: s3Link,
	}
	
	err = producers.AddAudioKafka(customFilename, s3Link)
	if(err != nil){
		http.Error(w, "Could not upload audio S3 link to Kafka", http.StatusInternalServerError);
		return;
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResp)
}