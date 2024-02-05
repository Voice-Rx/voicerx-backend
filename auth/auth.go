package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var refreshToken string

type CognitoErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	godotenv.Load("../.env")
	var requestBody struct {
		AuthorizationCode string `json:"authorizationCode"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for tokens
	tokenEndpoint := envVar("COGNITO_DOMAIN") + "/oauth2/token"
	tokenData := url.Values{}
	tokenData.Set("grant_type", "authorization_code")
	tokenData.Set("code", requestBody.AuthorizationCode)
	tokenData.Set("client_id", envVar("COGNITO_CLIENT_ID"))
	tokenData.Set("redirect_uri", "https://localhost:3000/dashboard")

	// Create an HTTP request with the necessary data
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(tokenData.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		http.Error(w, "Request creation failed", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authHead := "Basic " + envVar("AUTHORIZATION")
	// Add the Authorization header
	req.Header.Set("Authorization", authHead)

	// Perform the HTTP request
	tokenResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Token request failed: %v\n", err)
		http.Error(w, "Token request failed", http.StatusInternalServerError)
		return
	}
	defer tokenResponse.Body.Close()

	if tokenResponse.StatusCode == http.StatusOK {
		var tokenPayload map[string]interface{}
		err = json.NewDecoder(tokenResponse.Body).Decode(&tokenPayload)
		if err != nil {
			log.Printf("Failed to decode token response: %v\n", err)
			http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
			return
		}
		accessToken := tokenPayload["access_token"].(string)
		refreshToken := tokenPayload["refresh_token"].(string)
		idToken := tokenPayload["id_token"].(string)
		expirationTime := time.Now().Add(3600 * time.Second).Format(time.RFC3339)

		responsePayload := struct {
			AccessToken string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			IdToken string `json:"idToken"`
			ExpirationTime string `json:"expirationTime"`
		}{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			IdToken: idToken,
			ExpirationTime: expirationTime,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(tokenResponse.StatusCode)
		json.NewEncoder(w).Encode(responsePayload)

	} else {

		var errorResponse CognitoErrorResponse
		err := json.NewDecoder(tokenResponse.Body).Decode(&errorResponse)
		if err != nil {
			log.Printf("Failed to decode error response: %v\n", err)
			http.Error(w, "Failed to decode error response", http.StatusInternalServerError)
			return
		}

		// Return the Cognito error response to the frontend
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(tokenResponse.StatusCode)
		json.NewEncoder(w).Encode(errorResponse)
	}
}

// func GetToken(w http.ResponseWriter, r *http.Request){

// }

// func VerifyToken(){

// }

// func TokenIsValid(){

// }

func envVar(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
