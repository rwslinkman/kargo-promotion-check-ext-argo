package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ArgoApiLoginInterface interface {
	GetApiToken(username string, password string) (string, error)
}

type ArgoLoginClient struct {
	client HTTPClient
}

func NewArgoLoginClient(client HTTPClient) *ArgoLoginClient {
	return &ArgoLoginClient{
		client: client,
	}
}

type LoginResponse struct {
	AuthToken string `json:"token"`
}

func (c *ArgoLoginClient) GetApiToken(argoServer string, apiUsername string, apiPassword string, allowInsecure bool) (string, error) {
	loginPostData := map[string]string{
		"username": apiUsername,
		"password": apiPassword,
	}
	loginJsonData, err := json.Marshal(loginPostData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return "", err
	}

	// Create HTTP POST request
	protocol := "https"
	if allowInsecure {
		protocol = "http"
	}
	argoLoginUrl := fmt.Sprintf("%s://%s/api/v1/session", protocol, argoServer)
	req, err := http.NewRequest("POST", argoLoginUrl, bytes.NewBuffer(loginJsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	if resp.StatusCode != 200 {
		// TODO: Handle errors from ArgoCD API e.g. {"error":"Invalid username or password","code":16,"message":"Invalid username or password"}
		fmt.Println(string(body))
		return "", fmt.Errorf("token request not accepted by ArgoCD (http %d)", resp.StatusCode)
	}

	// Map JSON response to struct
	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	if loginResp.AuthToken == "" {
		fmt.Println("Unable to get API token from ArgoCD")
		return "", fmt.Errorf("unable to get API token from ArgoCD")
	}

	return loginResp.AuthToken, nil
}
