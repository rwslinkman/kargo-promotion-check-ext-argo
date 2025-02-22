package main

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"log"
	"os"
)

func main() {
	argoServer, _ := os.LookupEnv("ARGOCD_SERVER")
	argoApiToken, _ := os.LookupEnv("ARGOCD_API_TOKEN")
	argoAppName, _ := os.LookupEnv("ARGOCD_APP_NAME")
	// TODO; define modes LOGIN (username/password) or TOKEN (created using cli)

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
