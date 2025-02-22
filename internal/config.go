package internal

import (
	"fmt"
	"os"
)

type AuthMode string

const (
	LoginMode AuthMode = "LOGIN"
	TokenMode AuthMode = "TOKEN"
)

type Config struct {
	ArgoServer   string
	ArgoApiToken string
	ArgoAppName  string
	ApiUsername  string
	ApiPassword  string
	AuthMode     AuthMode
}

// LoadConfig reads environment variables and initializes the configuration
func LoadConfig() (*Config, error) {
	argoServer, hasServer := os.LookupEnv("ARGOCD_SERVER")
	argoApiToken, hasToken := os.LookupEnv("ARGOCD_API_TOKEN")
	argoAppName, hasAppName := os.LookupEnv("ARGOCD_APP_NAME")
	apiUsername, hasUsername := os.LookupEnv("ARGOCD_API_USERNAME")
	apiPassword, hasPassword := os.LookupEnv("ARGOCD_API_PASSWORD")

	// Ensure mandatory fields are present
	if !hasServer || !hasAppName || argoServer == "" || argoAppName == "" {
		return nil, fmt.Errorf("ARGOCD_SERVER and ARGOCD_APP_NAME must be set")
	}

	// Determine authentication mode
	var authMode AuthMode
	if hasToken && argoApiToken != "" {
		authMode = TokenMode
	} else {
		if !hasUsername || !hasPassword || apiUsername == "" || apiPassword == "" {
			return nil, fmt.Errorf("ARGOCD_API_USERNAME and ARGOCD_API_PASSWORD must be set for LOGIN mode")
		}
		authMode = LoginMode
	}

	// Return configuration struct
	return &Config{
		ArgoServer:   argoServer,
		ArgoApiToken: argoApiToken,
		ArgoAppName:  argoAppName,
		ApiUsername:  apiUsername,
		ApiPassword:  apiPassword,
		AuthMode:     authMode,
	}, nil
}
