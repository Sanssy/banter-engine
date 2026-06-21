package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DSanoussy/banter-engine/internal/banter"
	"github.com/DSanoussy/banter-engine/internal/discord"
	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/mpp"
	"github.com/DSanoussy/banter-engine/internal/opportunities"
	"github.com/DSanoussy/banter-engine/internal/snapshot"
)

const challengeID = "mpp_challenge_UDKDDH27"
const snapshotPath = "data/standings.json"

func main() {
	token := os.Getenv("MPP_TOKEN")
	if token == "" {
		log.Fatal("MPP_TOKEN environment variable is required")
	}
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("DISCORD_WEBHOOK_URL environment variable is required")
	}

	previousStandings, err := snapshot.LoadStandings(snapshotPath)
	if err != nil {
		log.Fatal(err)
	}

	client := mpp.NewClient(token)
	standings, err := client.GetStandings(challengeID)
	if err != nil {
		log.Fatal(err)
	}

	for _, standing := range standings {
		fmt.Printf("%d. %-12s %d\n", standing.Rank, standing.Name, standing.Points)
	}

	matches, err := client.GetMatches()
	if err != nil {
		log.Fatal(err)
	}
	var forecastHistory []forecasts.Forecast
	var messages []string
	for _, match := range matches {
		fmt.Printf("%s %d-%d %s (%s)\n", match.HomeTeam, match.Score.Home, match.Score.Away, match.AwayTeam, match.Status)

		forecasts, err := client.GetForecasts(challengeID, match)
		if err != nil {
			log.Fatal(err)
		}
		for _, forecast := range forecasts {
			fmt.Printf("  %s: %d-%d (%d points)\n", forecast.UserID, forecast.Prediction.Home, forecast.Prediction.Away, forecast.Points)
		}
		if match.Status == "fullTime" {
			forecastHistory = append(forecastHistory, forecasts...)
		}
		for _, opportunity := range opportunities.DetectSurprises(match, forecasts) {
			message := banter.Generate(opportunity)
			fmt.Println(message)
			messages = append(messages, message)
		}
	}
	for _, opportunity := range opportunities.DetectStreaks(forecastHistory) {
		message := banter.Generate(opportunity)
		fmt.Println(message)
		messages = append(messages, message)
	}

	for _, opportunity := range opportunities.Detect(previousStandings, standings) {
		message := banter.Generate(opportunity)
		fmt.Println(message)
		messages = append(messages, message)
	}

	discordClient := discord.NewClient(webhookURL)
	for _, message := range messages {
		if err := discordClient.Send(message); err != nil {
			log.Fatal(err)
		}
	}

	if err := snapshot.SaveStandings(snapshotPath, standings); err != nil {
		log.Fatal(err)
	}
}
