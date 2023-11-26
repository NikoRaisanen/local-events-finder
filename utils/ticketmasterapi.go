package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type EventsResponse struct {
	Embedded struct {
		Events []Event `json:"events"`
	} `json:"_embedded"`
	PageInfo PageInfo `json:"page"`
}

type PageInfo struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Number        int `json:"number"`
}

type Event struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ID   string `json:"id"`
	URL  string `json:"url"`
	// Images []Image  `json:"images"`
	Dates DateInfo `json:"dates"`
}

// type Image struct {
// 	Ratio    string `json:"ratio"`
// 	Url      string `json:"url"`
// 	Width    int    `json:"width"`
// 	Height   int    `json:"height"`
// 	Fallback bool   `json:"fallback"`
// }

type DateInfo struct {
	Start            StartInfo `json:"start"`
	SpanMultipleDays bool      `json:"spanMultipleDays"`
}

type StartInfo struct {
	DateTime string `json:"dateTime"`
}

func GetEvents(apiKey string, postalCode string, startTime string, endTime string) ([]Event, error) {
	// var allEvents []Event
	var eventsUrl string = fmt.Sprintf("https://app.ticketmaster.com/discovery/v2/events.json?postalCode=%s&startDateTime=%s&endDateTime=%s&page=0&size=20&apikey=%s", postalCode, startTime, endTime, apiKey)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error reading response body %w", err)
	}
	fmt.Printf("Response body: %s", body)
	err = ioutil.WriteFile("./tmp.json", body, 0644)
	var response EventsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	err = ioutil.WriteFile("./response.json", body, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error writing to file %w", err)
	}
	return response.Embedded.Events, nil
}
