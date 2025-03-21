package imagestore

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ImgurOpts struct {
	ClientID     string
	RefreshToken string
	ClientSecret string
	HttpClient   *http.Client
}

type imgurUploadRequest struct {
	Image       string `json:"image"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ImgurService struct {
	clientID     string
	refreshToken string
	clientSecret string
	accessToken  string

	httpClient *http.Client
}

type imgurUploadResponse struct {
	Status  uint16 `json:"status"`
	Success bool   `json:"success"`
	Data    struct {
		ID         string `json:"id"`
		DeleteHash string `json:"deletehash"`
		Link       string `json:"link"`
		Datetime   uint32 `json:"datetime"`
	} `json:"data"`
}

func NewImgurService(opts *ImgurOpts) (ImageStore, error) {
	service := &ImgurService{
		httpClient:   opts.HttpClient,
		clientID:     opts.ClientID,
		refreshToken: opts.RefreshToken,
		clientSecret: opts.ClientSecret,
	}

	accessToken, err := service.getAccessToken()

	if err != nil {
		log.Fatalln("Unable to create access token", err)
		return nil, err
	}

	service.accessToken = accessToken

	return service, nil
}

func (is *ImgurService) getAccessToken() (string, error) {
	authURL := "https://api.imgur.com/oauth2/token"

	formData := url.Values{}
	formData.Set("refresh_token", is.refreshToken)
	formData.Set("client_id", is.clientID)
	formData.Set("client_secret", is.clientSecret)
	formData.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set the appropriate header for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := is.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if token, ok := result["access_token"].(string); ok {
		log.Println("access token recieved")
		return token, nil
	}
	return "", fmt.Errorf("failed to get access token")
}

func (is *ImgurService) UploadImage(ctx context.Context, imageData string) (string, error) {
	// converting to base64 for test purposes
	file, err := os.Open(imageData)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read file content
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file %v", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(imageBytes)

	payload, err := json.Marshal(createUploadRequest(base64Str))
	if err != nil {
		return "", fmt.Errorf("failed to marshall upload request %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.imgur.com/3/image", bytes.NewBuffer(payload))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", is.accessToken))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "none", fmt.Errorf("failed to create upload request %w", err)
	}

	resp, err := is.httpClient.Do(request)

	if err != nil {
		return "none", err
	}
	defer resp.Body.Close()

	var result imgurUploadResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "none", fmt.Errorf("failed to marshal response, status: %v", resp.Status)
	}

	log.Printf("Image successfully uploaded: %s", resp.Status)
	return result.Data.Link, nil
}

func createUploadRequest(imageData string) imgurUploadRequest {
	return imgurUploadRequest{
		Type:        "base64",
		Title:       "Um",
		Description: "any",
		Image:       imageData,
	}

}
