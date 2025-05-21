package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"net/http"
	"os"
	"rwslinkman/kargo-promotion-check-ext-argo/internal"
	"strings"
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
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.AllowInsecure,
				},
			},
		}
		argoApiClient := internal.NewArgoLoginClient(client)
		var apiToken, err = argoApiClient.GetApiToken(config.ArgoServer, config.ApiUsername, config.ApiPassword, config.AllowInsecure)
		if err != nil {
			fmt.Println("Unable to get API token from ArgoCD: ", err)
			panic(err)
		}
		argoApiToken = apiToken
		fmt.Println("Successfully got a temporary API token from ArgoCD")
	}

	// Create API client with API token to interact with external Argo CD instance
	clientOpts := apiclient.ClientOptions{
		ServerAddr: config.ArgoServer,
		AuthToken:  argoApiToken,
		GRPCWeb:    true,
		Insecure:   config.AllowInsecure,
	}
	argoApiClient := apiclient.NewClientOrDie(&clientOpts)
	fmt.Println("ArgoCD API client created")

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

		fmt.Println("Fetching app details from ArgoCD...")
		argoApp, getErr := argoAppClient.Get(ctx, &appQuery)
		if getErr != nil {
			fmt.Printf("Failed to fetch App details: %v\n", getErr)
			fmt.Println("Retrying failed request")
			continue
		}

		fmt.Println("Sync Status:", argoApp.Status.Sync.Status)
		fmt.Println("Sync Revision:", argoApp.Status.Sync.Revision)
		fmt.Println("Health Status:", argoApp.Status.Health.Status)

		if argoApp.Status.Sync.Status == "Synced" && argoApp.Status.Health.Status == "Healthy" {
			if config.VerifyMode == internal.Exact {
				// Verify exact
				if argoApp.Status.Sync.Revision == config.TargetRevision {
					fmt.Println("App is synced, healthy, and at the expected target revision!")
					success = true
					break
				} else {
					fmt.Printf("App is synced, healthy, but not at expected revision. Expected %s but found %s \n", config.TargetRevision, argoApp.Status.Sync.Revision)
				}
			} else {
				// Fetch metadata for commit message
				revisionMetadata, fetchErr := argoAppClient.RevisionMetadata(ctx, &application.RevisionMetadataQuery{
					Name:     &config.ArgoAppName,
					Revision: &argoApp.Status.Sync.Revision,
				})
				if fetchErr != nil {
					fmt.Printf("Failed to get revision metadata: %v\n", fetchErr)
					panic(fetchErr)
				}

				fmt.Println("Synced Revision's Message: " + revisionMetadata.Message)
				match := strings.Contains(revisionMetadata.Message, config.SearchCommitMessage)
				if match {
					fmt.Println("App is synced, healthy, and commit message matches expectation!")
					success = true
					break
				} else {
					fmt.Println("App is synced, healthy, but commit message does not contain expected value")
				}
			}

		} else {
			fmt.Println("App is not in sync, retrying..")
		}

		// Success state not reached, try again after interval
		time.Sleep(config.PollInterval)
	}

	var exitCode = 1
	var exitMsgPart = " NOT"
	if success {
		exitCode = 0
		exitMsgPart = ""
	}
	fmt.Printf("Argo App '%s' is currently%s in expected state\n", config.ArgoAppName, exitMsgPart)
	fmt.Println("KPCEA completed")
	os.Exit(exitCode)
}
