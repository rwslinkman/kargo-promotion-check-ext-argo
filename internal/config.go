package internal

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type AuthMode string

const (
	LoginMode AuthMode = "LOGIN"
	TokenMode AuthMode = "TOKEN"
)

type Config struct {
	ArgoServer     string
	ArgoApiToken   string
	ArgoAppName    string
	ApiUsername    string
	ApiPassword    string
	AuthMode       AuthMode
	TargetRevision string
	PollTimeout    time.Duration
	PollInterval   time.Duration
	AllowInsecure  bool
}

// LoadConfig reads environment variables and initializes the configuration
func LoadConfig() (*Config, error) {
	argoServer, hasServer := os.LookupEnv("ARGOCD_SERVER")
	argoAppName, hasAppName := os.LookupEnv("ARGOCD_APP_NAME")
	targetRevision, hasTargetRevision := os.LookupEnv("KPCEA_TARGET_REVISION")
	// Ensure mandatory fields are present
	if !hasServer || !hasAppName || argoServer == "" || argoAppName == "" || !hasTargetRevision || targetRevision == "" {
		return nil, fmt.Errorf("KPCEA_TARGET_REVISION, ARGOCD_SERVER and ARGOCD_APP_NAME must be set")
	}

	argoApiToken, hasToken := os.LookupEnv("ARGOCD_API_TOKEN")
	apiUsername, hasUsername := os.LookupEnv("ARGOCD_API_USERNAME")
	apiPassword, hasPassword := os.LookupEnv("ARGOCD_API_PASSWORD")
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

	// Other (optional) configuration
	timeout, hasTimeout := os.LookupEnv("KPCEA_TIMEOUT")
	if !hasTimeout {
		timeout = "30"
	}
	timeoutSeconds, timeoutConfigErr := strconv.Atoi(timeout)
	if timeoutConfigErr != nil {
		return nil, fmt.Errorf("provided KPCEA_TIMEOUT must be a number")
	}
	interval, hasInterval := os.LookupEnv("KPCEA_INTERVAL")
	if !hasInterval {
		interval = "5"
	}
	intervalSeconds, intervalConfigErr := strconv.Atoi(interval)
	if intervalConfigErr != nil {
		return nil, fmt.Errorf("provided KPCEA_INTERVAL must be a number")
	}
	allowInsecure, hasInsecure := os.LookupEnv("KPCEA_INSECURE")
	if !hasInsecure {
		allowInsecure = "false"
	}

	// Return configuration struct
	return &Config{
		ArgoServer:     argoServer,
		ArgoApiToken:   argoApiToken,
		ArgoAppName:    argoAppName,
		ApiUsername:    apiUsername,
		ApiPassword:    apiPassword,
		AuthMode:       authMode,
		TargetRevision: targetRevision,
		PollTimeout:    time.Duration(timeoutSeconds) * time.Second,
		PollInterval:   time.Duration(intervalSeconds) * time.Second,
		AllowInsecure:  allowInsecure == "true",
	}, nil
}
