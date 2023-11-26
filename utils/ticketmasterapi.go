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

// Times are in UTC
type DateInfo struct {
	Start            StartInfo `json:"start"`
	SpanMultipleDays bool      `json:"spanMultipleDays"`
}

type StartInfo struct {
	DateTime string `json:"dateTime"`
}

func GetEvents(apiKey string, postalCode string, startTime string, endTime string) ([]Event, error) {
	var allEvents []Event
	var page int = 0

	for {
		eventsUrl := fmt.Sprintf("https://app.ticketmaster.com/discovery/v2/events.json?postalCode=%s&startDateTime=%s&endDateTime=%s&page=%d&size=20&apikey=%s", postalCode, startTime, endTime, page, apiKey)
		req, err := http.NewRequest("GET", eventsUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating new request object: %w", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error with API call: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var response EventsResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}

		allEvents = append(allEvents, response.Embedded.Events...)

		if response.PageInfo.Number >= response.PageInfo.TotalPages-1 {
			break
		}

		page++
	}

	return allEvents, nil
}
