package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"io"
	"log"
	"net/http"
	"os"
)

type LoginResponse struct {
	AuthToken string `json:"token"`
}

func main() {
	argoServer, _ := os.LookupEnv("ARGOCD_SERVER")
	argoApiToken, _ := os.LookupEnv("ARGOCD_API_TOKEN")
	argoAppName, _ := os.LookupEnv("ARGOCD_APP_NAME")
	apiUsername, isLoginMode := os.LookupEnv("ARGOCD_API_USERNAME")
	apiPassword, _ := os.LookupEnv("ARGOCD_API_PASSWORD")

	if isLoginMode {
		loginPostData := map[string]string{
			"username": apiUsername,
			"password": apiPassword,
		}
		loginJsonData, err := json.Marshal(loginPostData)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		// Create HTTP POST request
		argoLoginUrl := fmt.Sprintf("https://%s/api/v1/session", argoServer)
		req, err := http.NewRequest("POST", argoLoginUrl, bytes.NewBuffer(loginJsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		// Map JSON response to struct
		var loginResp LoginResponse
		err = json.Unmarshal(body, &loginResp)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		fmt.Println("Response Body:", loginResp.AuthToken)
		argoApiToken = loginResp.AuthToken
	}

	// Create API client options
	clientOpts := apiclient.ClientOptions{
		ServerAddr: argoServer,
		AuthToken:  argoApiToken,
		GRPCWeb:    true,
	}

	// Initialize the API client
	argoApiClient, err := apiclient.NewClient(&clientOpts)
	if err != nil {
		log.Fatalf("Failed to create Argo CD API client: %v", err)
		return
	}

	_, argoAppClient := argoApiClient.NewApplicationClientOrDie()
	appQuery := application.ApplicationQuery{Name: &argoAppName}

	var argoApp, getErr = argoAppClient.Get(context.Background(), &appQuery)
	if getErr != nil {
		log.Fatalf("Failed to fetch App details: %v", getErr)
	}

	fmt.Println("Sync Status:", argoApp.Status.Sync.Status)
	fmt.Println("Sync Revision:", argoApp.Status.Sync.Revision)
	fmt.Println("Health Status:", argoApp.Status.Health.Status)
}
