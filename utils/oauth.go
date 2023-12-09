package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"localevents/config"
	"net/http"
	"net/url"
	"strings"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	secretStore, err := config.GetSecrets()
	if err != nil {
		fmt.Errorf("error getting secrets: %w", err)
	}
	fmt.Printf("secrets: %v\n", secretStore)
	// Extract the code from query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	tokens, err := getAccessToken(code, secretStore)
	if err != nil {
		fmt.Errorf("error getting access token: %w", err)
	}
	// hardcode twitter user for now
	saveAccessToken("RenoLocalEvents", tokens)
}

func saveAccessToken(accout string, tokens AccessTokenResponse) error {
	err := config.UpdateTwitterCreds(accout, tokens.AccessToken, tokens.RefreshToken)
	if err != nil {
		return fmt.Errorf("error updating twitter bearer: %w", err)
	}
	return nil
}

// TODO: generalize this function so it handles first time oauth and refresh token
func getAccessToken(code string, secretStore config.Secrets) (AccessTokenResponse, error) {
	apiUrl := "https://api.twitter.com/2/oauth2/token"
	method := "POST"

	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", secretStore.Integrations.Oauth2.ClientId)
	data.Set("redirect_uri", "http://localhost:8080/oauth_callback")
	data.Set("code_verifier", "challenge")

	req, err := http.NewRequest(method, apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return AccessTokenResponse{}, fmt.Errorf("error creating new request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// clientid:clientsecret base64 encoded
	bAuth := []byte(secretStore.Integrations.Oauth2.ClientId + ":" + secretStore.Integrations.Oauth2.ClientSecret)
	encoded := base64.StdEncoding.EncodeToString(bAuth)
	req.Header.Set("Authorization", "Basic "+encoded)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AccessTokenResponse{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	var response AccessTokenResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return AccessTokenResponse{}, fmt.Errorf("error unmarshalling access token response: %w", err)
	}

	return response, nil
}

func refreshAccessToken(secretStore config.Secrets, refreshToken string) (string, error) {

	apiUrl := "https://api.twitter.com/2/oauth2/token"
	method := "POST"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest(method, apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating new request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// clientid:clientsecret base64 encoded
	bAuth := []byte(secretStore.Integrations.Oauth2.ClientId + ":" + secretStore.Integrations.Oauth2.ClientSecret)
	encoded := base64.StdEncoding.EncodeToString(bAuth)
	req.Header.Set("Authorization", "Basic "+encoded)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	var response AccessTokenResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return AccessTokenResponse{}, fmt.Errorf("error unmarshalling access token response: %w", err)
	}
}

func getAuthUrl(clientId string) (string, error) {
	redirectUri := "http://localhost:8080/oauth_callback"
	url := "https://twitter.com/i/oauth2/authorize?response_type=code&client_id=" + clientId + "&redirect_uri=" + redirectUri + "&scope=tweet.read%20tweet.write%20users.read%20follows.read%20offline.access&state=state&code_challenge=challenge&code_challenge_method=plain"

	return url, nil
}

func StartOauth(w http.ResponseWriter, r *http.Request) {
	secretStore, _ := config.GetSecrets()
	url, _ := getAuthUrl(secretStore.Integrations.Oauth2.ClientId)
	http.Redirect(w, r, url, http.StatusFound)
}
