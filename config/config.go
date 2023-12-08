package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Secrets struct {
	Integrations Integration `json:"integrations"`
}

type Integration struct {
	TickerMaster    APIKeySecret              `json:"ticketmaster"`
	Twitter         APIKeySecret              `json:"twitter"`
	TwitterAccounts map[string]TwitterAccount `json:"twitterAccounts"`
	Oauth2          Oauth2                    `json:"oauth2"`
}

type APIKeySecret struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type TwitterAccount struct {
	BearerToken       string  `json:"bearerToken"`
	AccessToken       string  `json:"accessToken"`
	AccessTokenSecret string  `json:"accessTokenSecret"`
	RefreshToken      *string `json:"refreshToken"`
}

type Oauth2 struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func GetSecrets() (Secrets, error) {
	// open file
	file, err := os.Open("./secrets.json")
	if err != nil {
		return Secrets{}, fmt.Errorf("Error reading secrets file %w", err)
	}
	defer file.Close()
	// read file data

	bytes, err := io.ReadAll(file)
	// unmarshall json

	var allSecrets Secrets
	err = json.Unmarshal(bytes, &allSecrets)
	if err != nil {
		return allSecrets, fmt.Errorf("Error unmarshalling json into struct %w", err)
	}

	return allSecrets, nil
}

func UpdateTwitterCreds(twitterAccount string, bearer string, refresh string) error {
	fmt.Printf("New bearer token: %s\n", bearer)
	fmt.Printf("New refresh token: %s\n", refresh)
	secretStore, _ := GetSecrets()

	account, ok := secretStore.Integrations.TwitterAccounts[twitterAccount]
	if !ok {
		return fmt.Errorf("Error finding %s in secrets.json", twitterAccount)
	}
	account.BearerToken = bearer
	account.RefreshToken = &refresh

	secretStore.Integrations.TwitterAccounts[twitterAccount] = account

	// Marshal the updated secrets back to JSON
	updatedBytes, err := json.Marshal(secretStore)
	if err != nil {
		return fmt.Errorf("Error marshalling updated secrets: %w", err)
	}

	// Write the updated JSON back to the file
	err = os.WriteFile("./secrets.json", updatedBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing updated secrets to file: %w", err)
	}

	return nil
}
