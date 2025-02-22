package main

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"log"
	"os"
	"rwslinkman/kargo-promotion-check-ext-argo/internal"
)

func main() {
	argoServer, _ := os.LookupEnv("ARGOCD_SERVER")
	argoApiToken, _ := os.LookupEnv("ARGOCD_API_TOKEN")
	argoAppName, _ := os.LookupEnv("ARGOCD_APP_NAME")
	apiUsername, isLoginMode := os.LookupEnv("ARGOCD_API_USERNAME")
	apiPassword, _ := os.LookupEnv("ARGOCD_API_PASSWORD")

	if isLoginMode {
		argoApiClient := internal.NewArgoApiClient()
		var apiToken, err = argoApiClient.GetApiToken(argoServer, apiUsername, apiPassword)
		if err != nil {
			fmt.Println("Unable to get API token from ArgoCD: ", err)
		}
		fmt.Println("Response Body:", apiToken)
		argoApiToken = apiToken
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
