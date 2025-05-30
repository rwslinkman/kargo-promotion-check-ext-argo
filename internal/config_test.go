package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func setEnvVars(t *testing.T, envs map[string]string) func() {
	t.Helper()

	originals := make(map[string]string)
	for key, val := range envs {
		// Save original value (if any) for full cleanup later
		originals[key] = os.Getenv(key)
		err := os.Setenv(key, val)
		if err != nil {
			t.Fatalf("failed to set env var %s: %v", key, err)
		}
	}

	return func() {
		for key := range envs {
			if orig, ok := originals[key]; ok && orig != "" {
				_ = os.Setenv(key, orig)
			} else {
				_ = os.Unsetenv(key)
			}
		}
	}
}

func TestLoadConfig_MinimalValidEnvVars_TokenMode(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, Exact, config.VerifyMode)
	assert.Equal(t, "target-revision", config.TargetRevision)
	assert.Equal(t, "", config.SearchCommitMessage)
	assert.Equal(t, TokenMode, config.AuthMode)
	assert.Equal(t, "api-token", config.ArgoApiToken)
	assert.Equal(t, "", config.ApiUsername)
	assert.Equal(t, "", config.ApiPassword)
	assert.Equal(t, 30*time.Second, config.PollTimeout)
	assert.Equal(t, 5*time.Second, config.PollInterval)
	assert.Equal(t, false, config.AllowInsecure)
}

func TestLoadConfig_MinimalValidEnvVars_LoginMode(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_USERNAME":   "api-username",
		"ARGOCD_API_PASSWORD":   "api-password",
		"KPCEA_TIMEOUT":         "20",
		"KPCEA_INTERVAL":        "3",
		"KPCEA_INSECURE":        "true",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, Exact, config.VerifyMode)
	assert.Equal(t, "target-revision", config.TargetRevision)
	assert.Equal(t, "", config.SearchCommitMessage)
	assert.Equal(t, LoginMode, config.AuthMode)
	assert.Equal(t, "", config.ArgoApiToken)
	assert.Equal(t, "api-username", config.ApiUsername)
	assert.Equal(t, "api-password", config.ApiPassword)
	assert.Equal(t, 20*time.Second, config.PollTimeout)
	assert.Equal(t, 3*time.Second, config.PollInterval)
	assert.Equal(t, true, config.AllowInsecure)
}

func TestLoadConfig_AllValidEnvVars_TokenMode(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
		"KPCEA_TIMEOUT":         "20",
		"KPCEA_INTERVAL":        "3",
		"KPCEA_INSECURE":        "true",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, Exact, config.VerifyMode)
	assert.Equal(t, "target-revision", config.TargetRevision)
	assert.Equal(t, "", config.SearchCommitMessage)
	assert.Equal(t, TokenMode, config.AuthMode)
	assert.Equal(t, "api-token", config.ArgoApiToken)
	assert.Equal(t, "", config.ApiUsername)
	assert.Equal(t, "", config.ApiPassword)
	assert.Equal(t, 20*time.Second, config.PollTimeout)
	assert.Equal(t, 3*time.Second, config.PollInterval)
	assert.Equal(t, true, config.AllowInsecure)
}

func TestLoadConfig_AllValidEnvVars_LoginMode(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_USERNAME":   "api-username",
		"ARGOCD_API_PASSWORD":   "api-password",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, Exact, config.VerifyMode)
	assert.Equal(t, "target-revision", config.TargetRevision)
	assert.Equal(t, "", config.SearchCommitMessage)
	assert.Equal(t, LoginMode, config.AuthMode)
	assert.Equal(t, "", config.ArgoApiToken)
	assert.Equal(t, "api-username", config.ApiUsername)
	assert.Equal(t, "api-password", config.ApiPassword)
	assert.Equal(t, 30*time.Second, config.PollTimeout)
	assert.Equal(t, 5*time.Second, config.PollInterval)
	assert.Equal(t, false, config.AllowInsecure)
}

func TestLoadConfig_MissingArgoServerProperty(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "ARGOCD_SERVER and ARGOCD_APP_NAME must be set", err.Error())
}

func TestLoadConfig_EmptyArgoServerProperty(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "ARGOCD_SERVER and ARGOCD_APP_NAME must be set", err.Error())
}

func TestLoadConfig_MissingArgoAppNameProperty(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "ARGOCD_SERVER and ARGOCD_APP_NAME must be set", err.Error())
}

func TestLoadConfig_EmptyArgoAppNameProperty(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "ARGOCD_SERVER and ARGOCD_APP_NAME must be set", err.Error())
}

func TestLoadConfig_MissingTargetRevisionProperty(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":    "argocd-server",
		"ARGOCD_APP_NAME":  "argo-app-name",
		"ARGOCD_API_TOKEN": "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "KPCEA_TARGET_REVISION must be set for verification mode EXACT", err.Error())
}

func TestLoadConfig_EmptyTargetRevisionProperty(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "KPCEA_TARGET_REVISION must be set for verification mode EXACT", err.Error())
}

func TestLoadConfig_MissingCredentialsOrTokenProperties(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "ARGOCD_API_USERNAME and ARGOCD_API_PASSWORD must be set for LOGIN mode", err.Error())
}

func TestLoadConfig_TakeTokenModeOverLoginMode(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_USERNAME":   "api-username",
		"ARGOCD_API_PASSWORD":   "api-password",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, "target-revision", config.TargetRevision)
	assert.Equal(t, TokenMode, config.AuthMode)
	assert.Equal(t, Exact, config.VerifyMode)
}

func TestLoadConfig_InvalidTimeoutValue(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
		"KPCEA_TIMEOUT":         "twenty",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "provided KPCEA_TIMEOUT must be a number", err.Error())
}

func TestLoadConfig_InvalidIntervalValue(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
		"KPCEA_INTERVAL":        "twenty",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "provided KPCEA_INTERVAL must be a number", err.Error())
}

func TestLoadConfig_InvalidInsecureValueDoesNotMakeItSecure(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
		"KPCEA_INSECURE":        "yes",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, false, config.AllowInsecure)
}

func TestLoadConfig_ProvidedSearchCommitMsgInsteadOfTargetRevision(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":           "argocd-server",
		"ARGOCD_APP_NAME":         "argo-app-name",
		"KPCEA_SEARCH_COMMIT_MSG": "search-param",
		"ARGOCD_API_TOKEN":        "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "KPCEA_TARGET_REVISION must be set for verification mode EXACT", err.Error())
}

func TestLoadConfig_InvalidVerifyModeValueDefaultsToExact(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_VERIFY_MODE":     "invalid-value",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, Exact, config.VerifyMode)
	assert.Equal(t, "target-revision", config.TargetRevision)
	assert.Equal(t, "", config.SearchCommitMessage)
	assert.Equal(t, TokenMode, config.AuthMode)
	assert.Equal(t, "api-token", config.ArgoApiToken)
	assert.Equal(t, "", config.ApiUsername)
	assert.Equal(t, "", config.ApiPassword)
	assert.Equal(t, 30*time.Second, config.PollTimeout)
	assert.Equal(t, 5*time.Second, config.PollInterval)
	assert.Equal(t, false, config.AllowInsecure)
}

func TestLoadConfig_VerifyModeSearchSelectedButNoParameterProvided(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":     "argocd-server",
		"ARGOCD_APP_NAME":   "argo-app-name",
		"KPCEA_VERIFY_MODE": "SEARCH_COMMIT_MSG",
		"ARGOCD_API_TOKEN":  "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "KPCEA_SEARCH_COMMIT_MSG must be set for verification mode SEARCH_COMMIT_MSG", err.Error())
}

func TestLoadConfig_VerifyModeSearchSelectedButEmptyParameterProvided(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":           "argocd-server",
		"ARGOCD_APP_NAME":         "argo-app-name",
		"KPCEA_VERIFY_MODE":       "SEARCH_COMMIT_MSG",
		"KPCEA_SEARCH_COMMIT_MSG": "",
		"ARGOCD_API_TOKEN":        "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "KPCEA_SEARCH_COMMIT_MSG must be set for verification mode SEARCH_COMMIT_MSG", err.Error())
}

func TestLoadConfig_VerifyModeSearchSelectedButTargetRevisionProvided(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":         "argocd-server",
		"ARGOCD_APP_NAME":       "argo-app-name",
		"KPCEA_VERIFY_MODE":     "SEARCH_COMMIT_MSG",
		"KPCEA_TARGET_REVISION": "target-revision",
		"ARGOCD_API_TOKEN":      "api-token",
	})
	defer cleanup()

	_, err := LoadConfig()

	assert.Error(t, err)
	assert.Equal(t, "KPCEA_SEARCH_COMMIT_MSG must be set for verification mode SEARCH_COMMIT_MSG", err.Error())
}

func TestLoadConfig_MinimalValidEnvVars_SearchCommitMsgMode(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{
		"ARGOCD_SERVER":           "argocd-server",
		"ARGOCD_APP_NAME":         "argo-app-name",
		"KPCEA_VERIFY_MODE":       "SEARCH_COMMIT_MSG",
		"KPCEA_SEARCH_COMMIT_MSG": "search-param",
		"ARGOCD_API_TOKEN":        "api-token",
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, SearchCommitMessage, config.VerifyMode)
	assert.Equal(t, "", config.TargetRevision)
	assert.Equal(t, "search-param", config.SearchCommitMessage)
	assert.Equal(t, TokenMode, config.AuthMode)
	assert.Equal(t, "api-token", config.ArgoApiToken)
	assert.Equal(t, "", config.ApiUsername)
	assert.Equal(t, "", config.ApiPassword)
	assert.Equal(t, 30*time.Second, config.PollTimeout)
	assert.Equal(t, 5*time.Second, config.PollInterval)
	assert.Equal(t, false, config.AllowInsecure)
}
