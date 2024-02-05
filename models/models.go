package models

// Define the Audio struct
type AudioLink struct {
    Cred   string    `json:"cred"`
    S3Link string    `json:"s3link"`
}

type HttpResp struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}