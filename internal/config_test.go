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
				os.Setenv(key, orig)
			} else {
				os.Unsetenv(key)
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
	assert.Equal(t, "target-revision", config.TargetRevision)
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
	})
	defer cleanup()

	config, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "argocd-server", config.ArgoServer)
	assert.Equal(t, "argo-app-name", config.ArgoAppName)
	assert.Equal(t, "target-revision", config.TargetRevision)
	assert.Equal(t, LoginMode, config.AuthMode)
	assert.Equal(t, "", config.ArgoApiToken)
	assert.Equal(t, "api-username", config.ApiUsername)
	assert.Equal(t, "api-password", config.ApiPassword)
	assert.Equal(t, 30*time.Second, config.PollTimeout)
	assert.Equal(t, 5*time.Second, config.PollInterval)
	assert.Equal(t, false, config.AllowInsecure)
}
