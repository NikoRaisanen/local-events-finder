package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EventsResponse struct {
	Embedded struct {
		Events []Event `json:"events"`
	} `json:"_embedded"`
}

type Event struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	ID     string   `json:"id"`
	URL    string   `json:"url"`
	Images []Image  `json:"images"`
	Dates  DateInfo `json:"dates"`
}

type Image struct {
	Ratio    string `json:"ratio"`
	Url      string `json:"url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Fallback bool   `json:"fallback"`
}

type DateInfo struct {
	Start            StartInfo `json:"start"`
	SpanMultipleDays bool      `json:"spanMultipleDays"`
}

type StartInfo struct {
	DateTime string `json:"dateTime"`
}

func GetEvents(apiKey string) ([]Event, error) {
	var eventsUrl string = fmt.Sprintf("https://app.ticketmaster.com/discovery/v2/events.json?countryCode=US&apikey=%s", apiKey)
	req, err := http.NewRequest("GET", eventsUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating new req object %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error with api call %w", err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := io.ReadAll(resp.Body)
	// err = ioutil.WriteFile("./events.json", body, 0644)
	// if err != nil {
	// 	log.Fatalf("Error writing to file %w", err)
	// }
	var response EventsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// eventsJson, _ := json.Marshal(response.Embedded.Events)
	// err = ioutil.WriteFile("./response.json", eventsJson, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error writing to file %w", err)
	}
	return response.Embedded.Events, nil
}
