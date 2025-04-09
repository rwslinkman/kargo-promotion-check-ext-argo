package internal

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockHTTPClient struct {
	Resp *http.Response
	Err  error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Resp, m.Err
}

func NewMockHTTPClient(statusCode int, responseBody string, err error) *MockHTTPClient {
	mockResp := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
	return &MockHTTPClient{
		Resp: mockResp,
		Err:  err,
	}
}

func TestArgoLoginClient_GetApiToken(t *testing.T) {
	mockClient := NewMockHTTPClient(200, `{"token": "mock-token"}`, nil)
	argoClient := NewArgoLoginClient(mockClient)

	token, err := argoClient.GetApiToken("myServer", "myUser", "myPass", false)

	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

func TestArgoLoginClient_GetApiToken_Insecure(t *testing.T) {
	mockClient := NewMockHTTPClient(200, `{"token": "mock-token"}`, nil)
	argoClient := NewArgoLoginClient(mockClient)

	token, err := argoClient.GetApiToken("myServer", "myUser", "myPass", true)

	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

func TestArgoLoginClient_GetApiToken_FailedHttpRequest(t *testing.T) {
	mockClient := NewMockHTTPClient(401, `{"error":"Invalid username or password","code":16,"message":"Invalid username or password"}`, nil)
	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("myServer", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "token request not accepted by ArgoCD (http 401)")
}

func TestArgoLoginClient_GetApiToken_EmptyToken(t *testing.T) {
	mockClient := NewMockHTTPClient(200, `{"token": ""}`, nil)

	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("myServer", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "unable to get API token from ArgoCD")
}
