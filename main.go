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

	config, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}

	argoApiToken := config.ArgoApiToken // might be nil
	if config.AuthMode == internal.LoginMode {
		// ensure having an API Token
		argoApiClient := internal.NewArgoLoginClient()
		var apiToken, err = argoApiClient.GetApiToken(config.ArgoServer, config.ApiUsername, config.ApiPassword)
		if err != nil {
			fmt.Println("Unable to get API token from ArgoCD: ", err)
		}
		argoApiToken = apiToken
	}

	// Create API client with API token to interact with external Argo CD instance
	clientOpts := apiclient.ClientOptions{
		ServerAddr: config.ArgoServer,
		AuthToken:  argoApiToken,
		GRPCWeb:    true,
	}
	argoApiClient, err := apiclient.NewClient(&clientOpts)
	if err != nil {
		log.Fatalf("Failed to create Argo CD API client: %v", err)
		return
	}

	_, argoAppClient := argoApiClient.NewApplicationClientOrDie()
	appQuery := application.ApplicationQuery{Name: &config.ArgoAppName}

	var argoApp, getErr = argoAppClient.Get(context.Background(), &appQuery)
	if getErr != nil {
		log.Fatalf("Failed to fetch App details: %v", getErr)
	}

	fmt.Println("Sync Status:", argoApp.Status.Sync.Status)
	fmt.Println("Sync Revision:", argoApp.Status.Sync.Revision)
	fmt.Println("Health Status:", argoApp.Status.Health.Status)

	if argoApp.Status.Sync.Status == "Synced" {
		fmt.Println(fmt.Sprintf("Argo App %s is currently synced\n", config.ArgoAppName))
		os.Exit(0)
	} else {
		fmt.Println(fmt.Sprintf("Argo App %s is currently not in sync\n", config.ArgoAppName))
		os.Exit(1)
	}
}
