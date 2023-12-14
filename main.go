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

	// "flag"
	"fmt"
	"localevents/config"
	"localevents/utils"
	"log"
	"net/http"
	"time"
)

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

func twitterAuth() {
	http.HandleFunc("/start_oauth", utils.StartOauth)
	http.HandleFunc("/oauth_callback", utils.CallbackHandler)
	// Start the HTTP server
	log.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	// Code paths:
	// Try to use access token -> if it fails, go through oauth workflow

	// get refresh token
	secretStore, err := config.GetSecrets()
	if err != nil {
		log.Fatalf("Error occured fetching secrets %s", err)
	}
	// refreshToken := secretStore.Integrations.TwitterAccounts["RenoLocalEvents"].RefreshToken
	// exchange refresh token for access token
	// twitterToken, err := utils.RefreshAccessToken(secretStore, *refreshToken)

	timeNow := time.Now().Format("2006-01-02T15:04:05Z")
	endTime := time.Now().AddDate(0, 0, 7).Format("2006-01-02T15:04:05Z")
	postalCode := "89501"
	resp, err := utils.GetEvents(secretStore.Integrations.TickerMaster.Key, postalCode, timeNow, endTime)
	condensedEvents, err := aggregateDuplicates(resp)
	summarizeEvents(condensedEvents)

	// utils.CreateTweet(twitterToken, "How is everyone doing?")
}
