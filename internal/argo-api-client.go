package internal

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ArgoApiLoginInterface interface {
	GetApiToken(username string, password string) (string, error)
}

type ArgoLoginClient struct{}

func NewArgoLoginClient() *ArgoLoginClient {
	return &ArgoLoginClient{}
}

type LoginResponse struct {
	AuthToken string `json:"token"`
}

func (c *ArgoLoginClient) GetApiToken(argoServer string, apiUsername string, apiPassword string, allowInsecure bool) (string, error) {
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
	protocol := "https"
	if allowInsecure {
		protocol = "http"
	}
	argoLoginUrl := fmt.Sprintf("%s://%s/api/v1/session", protocol, argoServer)
	req, err := http.NewRequest("POST", argoLoginUrl, bytes.NewBuffer(loginJsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: allowInsecure, // Skip TLS verification if true
		},
	}
	client := &http.Client{
		Transport: tr,
	}
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
