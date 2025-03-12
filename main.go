package main

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"os"
	"rwslinkman/kargo-promotion-check-ext-argo/internal"
	"time"
)

func main() {
	config, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("KPCEA started in %s mode \n", config.AuthMode)

	argoApiToken := config.ArgoApiToken // might be nil
	if config.AuthMode == internal.LoginMode {
		// ensure having an API Token
		argoApiClient := internal.NewArgoLoginClient()
		var apiToken, err = argoApiClient.GetApiToken(config.ArgoServer, config.ApiUsername, config.ApiPassword, config.AllowInsecure)
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
		Insecure:   true,
	}
	argoApiClient := apiclient.NewClientOrDie(&clientOpts)

	_, argoAppClient := argoApiClient.NewApplicationClientOrDie()
	appQuery := application.ApplicationQuery{Name: &config.ArgoAppName}

	ctx := context.Background()
	start := time.Now()
	success := false

	for {
		if time.Since(start) > config.PollTimeout {
			fmt.Println("Timeout reached while waiting for app to sync")
			break
		}

		argoApp, getErr := argoAppClient.Get(ctx, &appQuery)
		if getErr != nil {
			fmt.Printf("Failed to fetch App details: %v\n", getErr)
			continue
		}

		fmt.Println("Sync Status:", argoApp.Status.Sync.Status)
		fmt.Println("Sync Revision:", argoApp.Status.Sync.Revision)
		fmt.Println("Health Status:", argoApp.Status.Health.Status)

		if argoApp.Status.Sync.Status == "Synced" &&
			argoApp.Status.Health.Status == "Healthy" &&
			argoApp.Status.Sync.Revision == config.TargetRevision {
			fmt.Println("App is synced, healthy, and at the correct revision!")
			success = true
			break
		} else {
			fmt.Println("App is not in sync, retrying..")
		}

		// Success state not reached, try again after interval
		time.Sleep(config.PollInterval)
	}

	if success {
		fmt.Println(fmt.Sprintf("Argo App %s is currently synced\n", config.ArgoAppName))
		os.Exit(0)
	} else {
		fmt.Println(fmt.Sprintf("Argo App %s is currently not in sync\n", config.ArgoAppName))
		os.Exit(1)
	}
}
