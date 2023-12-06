package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"localevents/config"
	"net/http"
	"net/url"
	"strings"
)

// api ref: https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/api-reference/post-tweets
// TODO:
// 		add geo location to tweets
// 		add image to tweets
// 		maybe some kind of poll ?

type RequestBody struct {
	Text string `json:"text"`
}

type AccessTokenBody struct {
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	RedirectUri  string `json:"redirect_uri"`
	CodeVerifier string `json:"code_verifier"`
}

func GetAuthUrl(clientId string) (string, error) {
	// TODO: parameterize the url
	redirectUri := "http://localhost:8080/oauth_callback"
	// redirectUri := "https://nikoraisanen.com"
	url := "https://twitter.com/i/oauth2/authorize?response_type=code&client_id=" + clientId + "&redirect_uri=" + redirectUri + "&scope=tweet.read%20tweet.write%20users.read%20follows.read%20offline.access&state=state&code_challenge=challenge&code_challenge_method=plain"
	fmt.Printf("auth url: %s\n", url)

	return url, nil
}

func CreateTweet(accessToken string) error {
	url := "https://api.twitter.com/2/tweets"
	method := "POST"
	jsonBody, err := json.Marshal(&RequestBody{Text: "Hello World!"})
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("error creating new request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Printf("Response from POST Tweet endpoint: %s\n", bodyString)
	return nil
}

func GetAccessToken(code string, secretStore config.Secrets) (string, error) {
	apiUrl := "https://api.twitter.com/2/oauth2/token"
	method := "POST"

	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", secretStore.Integrations.Oauth2.ClientId)
	data.Set("redirect_uri", "https://nikoraisanen.com")
	data.Set("code_verifier", "challenge")

	req, err := http.NewRequest(method, apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating new request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// clientid:clientsecret base64 encoded
	bAuth := []byte(secretStore.Integrations.Oauth2.ClientId + ":" + secretStore.Integrations.Oauth2.ClientSecret)
	encoded := base64.StdEncoding.EncodeToString(bAuth)
	fmt.Printf("encoded: %s\n", encoded)
	req.Header.Set("Authorization", "Basic "+encoded)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	bodyString := string(respBody)
	fmt.Printf("response body: %s\n", bodyString)

	return "", nil
}

func FetchNewTwitterToken(secretStore config.Secrets) (string, error) {
	oauthCode, err := GetAuthUrl(secretStore.Integrations.Oauth2.ClientId)
	if err != nil {
		return "", fmt.Errorf("error getting link for oauth code: %w", err)
	}
	accessToken, _ := GetAccessToken(oauthCode, secretStore)
	return accessToken, err
}
