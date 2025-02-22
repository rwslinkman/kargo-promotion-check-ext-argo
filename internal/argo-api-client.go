package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ArgoApiClientInterface interface {
	GetApiToken(username string, password string) (string, error)
}

type ArgoApiClient struct{}

func NewArgoApiClient() *ArgoApiClient {
	return &ArgoApiClient{}
}

type LoginResponse struct {
	AuthToken string `json:"token"`
}

func (c *ArgoApiClient) GetApiToken(argoServer string, apiUsername string, apiPassword string) (string, error) {
	loginPostData := map[string]string{
		"username": apiUsername,
		"password": apiPassword,
	}
	loginJsonData, err := json.Marshal(loginPostData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return "", err
	}

	// Create HTTP POST request
	argoLoginUrl := fmt.Sprintf("https://%s/api/v1/session", argoServer)
	req, err := http.NewRequest("POST", argoLoginUrl, bytes.NewBuffer(loginJsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Map JSON response to struct
	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	return loginResp.AuthToken, nil
}
