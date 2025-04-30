package internal

import (
	"fmt"
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

func NewMockHTTPClient(statusCode int, hasResponse bool, responseBody string, err error) *MockHTTPClient {
	response := strings.NewReader(responseBody)
	if !hasResponse {
		response = nil
	}
	mockResp := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(response),
	}
	return &MockHTTPClient{
		Resp: mockResp,
		Err:  err,
	}
}

func TestArgoLoginClient_GetApiToken(t *testing.T) {
	mockClient := NewMockHTTPClient(200, true, `{"token": "mock-token"}`, nil)
	argoClient := NewArgoLoginClient(mockClient)

	token, err := argoClient.GetApiToken("myServer", "myUser", "myPass", false)

	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

func TestArgoLoginClient_GetApiToken_Insecure(t *testing.T) {
	mockClient := NewMockHTTPClient(200, true, `{"token": "mock-token"}`, nil)
	argoClient := NewArgoLoginClient(mockClient)

	token, err := argoClient.GetApiToken("myServer", "myUser", "myPass", true)

	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

func TestArgoLoginClient_GetApiToken_FailedHttpRequest(t *testing.T) {
	mockClient := NewMockHTTPClient(401, true, `{"error":"Invalid username or password","code":16,"message":"Invalid username or password"}`, nil)
	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("myServer", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, "token request not accepted by ArgoCD (http 401)", err.Error())
}

func TestArgoLoginClient_GetApiToken_EmptyToken(t *testing.T) {
	mockClient := NewMockHTTPClient(200, true, `{"token": ""}`, nil)

	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("myServer", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, "unable to get API token from ArgoCD", err.Error())
}

func TestArgoLoginClient_GetApiToken_UnableToCreateRequest(t *testing.T) {
	mockClient := NewMockHTTPClient(0, false, "", nil)
	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("\t", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, `parse "http://\t/api/v1/session": net/url: invalid control character in URL`, err.Error())
}

func TestArgoLoginClient_GetApiToken_HttpClientError(t *testing.T) {
	mockClient := NewMockHTTPClient(0, false, "", fmt.Errorf("testError"))
	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("argoServer", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, "testError", err.Error())
}

func TestArgoLoginClient_GetApiToken_ResponseBodyIsNotJson(t *testing.T) {
	mockClient := NewMockHTTPClient(200, true, "", nil)
	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("argoServer", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, "unexpected end of JSON input", err.Error())
}

func TestArgoLoginClient_GetApiToken_ResponseBodyHasNoToken(t *testing.T) {
	mockClient := NewMockHTTPClient(200, true, `{"notToken": ""}`, nil)
	argoClient := NewArgoLoginClient(mockClient)

	_, err := argoClient.GetApiToken("argoServer", "myUser", "myPass", true)

	assert.Error(t, err)
	assert.Equal(t, "unable to get API token from ArgoCD", err.Error())
}

// TODO: Add test for error in io.ReadAll
