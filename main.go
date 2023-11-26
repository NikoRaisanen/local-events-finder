package main

// Event apis to pull from https://github.com/public-apis-dev/public-apis#events
// TODO:
// 1. If multiple times available for event at same venue, time, an date, combine them into one event in end result
// 2. Set up twitter api to tweet out events
// 3. Implement pagination in GetEvents func... Can I get the num of pages up front and then create a goroutine for each page?
// 4. Long term - Set up db to store events
// 5. Set up project structure according to golang standards
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

func summarizeEvents(events []utils.Event) {
	for _, event := range events {
		fmt.Println(event.Name)
	}
}

func main() {
	ticketMasterSecret, err := getSecrets()
	if err != nil {
		log.Fatalf("Error occured fetching secrets %s", err)
	}

	resp, err := utils.GetEvents(ticketMasterSecret.Key)
	// fmt.Printf("Response from GetEvents: %v", resp)
	fmt.Printf("Number of events: %d", len(resp))
	summarizeEvents(resp)
}
