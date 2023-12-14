package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func CreateTweet(accessToken string, text string) error {
	url := "https://api.twitter.com/2/tweets"
	method := "POST"
	jsonBody, err := json.Marshal(&RequestBody{Text: text})
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
