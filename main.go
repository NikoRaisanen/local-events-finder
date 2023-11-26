package main

// Event apis to pull from https://github.com/public-apis-dev/public-apis#events
// TODO:
// 1. If multiple times available for event at same venue, time, an date, combine them into one event in end result
// 2. Set up twitter api to tweet out events
// 3. Implement pagination in GetEvents func... Can I get the num of pages up front and then create a goroutine for each page?
// 4. Long term - Set up db to store events
// 5. Set up project structure according to golang standards

// MVP:
// 1. Get Reno events for next week
// 2. Tweet out events for next week
import (
	"encoding/json"
	"fmt"
	"io"
	"localevents/utils"
	"log"
	"os"
	"time"
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

func summarizeEvents(dailyEvents map[string][]utils.Event) {
	for date, events := range dailyEvents {
		fmt.Printf("Events for %s", date)
		for _, event := range events {
			fmt.Printf("\n\t%s", event.Name)
		}
		fmt.Println()
	}
}

func aggregateDuplicates(events []utils.Event) (map[string][]utils.Event, error) {
	dailyEvents := make(map[string][]utils.Event)
	for _, event := range events {
		parsedTime, err := time.Parse(time.RFC3339, event.Dates.Start.DateTime)
		if err != nil {
			return nil, fmt.Errorf("Error parsing time %w", err)
		}
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			return nil, fmt.Errorf("Error loading location %w", err)
		}
		localTime := parsedTime.In(loc)
		truncatedTime := localTime.Format("2006-01-02")
		if _, ok := dailyEvents[truncatedTime]; ok {
			dailyEvents[truncatedTime] = append(dailyEvents[truncatedTime], event)
		} else {
			dailyEvents[truncatedTime] = []utils.Event{event}
		}
	}
	return dailyEvents, nil
}

func main() {
	ticketMasterSecret, err := getSecrets()
	if err != nil {
		log.Fatalf("Error occured fetching secrets %s", err)
	}

	timeNow := time.Now().Format("2006-01-02T15:04:05Z")
	endTime := time.Now().AddDate(0, 0, 7).Format("2006-01-02T15:04:05Z")
	postalCode := "89501"
	resp, err := utils.GetEvents(ticketMasterSecret.Key, postalCode, timeNow, endTime)
	condensedEvents, err := aggregateDuplicates(resp)
	summarizeEvents(condensedEvents)
}
