package main

// Event apis to pull from https://github.com/public-apis-dev/public-apis#events
import (
	"encoding/json"
	"fmt"
	"io"
	"localevents/utils"
	"log"
	"os"
)

type Secret struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func getSecrets() (Secret, error) {
	// open file
	file, err := os.Open("./secrets.json")
	if err != nil {
		return Secret{}, fmt.Errorf("Error reading secrets file %w", err)
	}
	defer file.Close()
	// read file data

	bytes, err := io.ReadAll(file)
	// unmarshall json

	var ticketMasterSecret Secret
	err = json.Unmarshal(bytes, &ticketMasterSecret)
	if err != nil {
		return ticketMasterSecret, fmt.Errorf("Error unmarshalling json into struct %w", err)
	}

	return ticketMasterSecret, nil
}

func main() {
	ticketMasterSecret, err := getSecrets()
	if err != nil {
		log.Fatalf("Error occured fetching secrets %s", err)
	}

	resp, err := utils.GetEvents(ticketMasterSecret.Key)
	fmt.Printf("Response from GetEvents: %v", resp)
}
