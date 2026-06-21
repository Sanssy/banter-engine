package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DSanoussy/banter-engine/internal/banter"
	"github.com/DSanoussy/banter-engine/internal/discord"
	"github.com/DSanoussy/banter-engine/internal/forecasts"
	matchmodel "github.com/DSanoussy/banter-engine/internal/matches"
	"github.com/DSanoussy/banter-engine/internal/mpp"
	"github.com/DSanoussy/banter-engine/internal/opportunities"
	"github.com/DSanoussy/banter-engine/internal/rivalries"
	"github.com/DSanoussy/banter-engine/internal/snapshot"
)

const challengeID = "mpp_challenge_UDKDDH27"
const snapshotPath = "data/standings.json"
const matchesSnapshotPath = "data/matches.json"
const rivalriesSnapshotPath = "data/rivalries.json"
const runInterval = 5 * time.Minute

func main() {
	token := os.Getenv("MPP_TOKEN")
	if token == "" {
		log.Fatal("MPP_TOKEN environment variable is required")
	}
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("DISCORD_WEBHOOK_URL environment variable is required")
	}

	client := mpp.NewClient(token)
	discordClient := discord.NewClient(webhookURL)
	ticker := time.NewTicker(runInterval)
	defer ticker.Stop()

	for {
		if err := run(client, discordClient); err != nil {
			log.Printf("banter engine run failed: %v", err)
		}
		<-ticker.C
	}
}

func run(client *mpp.Client, discordClient *discord.Client) error {
	previousStandings, err := snapshot.LoadStandings(snapshotPath)
	if err != nil {
		return err
	}
	previousMatches, err := snapshot.LoadMatches(matchesSnapshotPath)
	if err != nil {
		return err
	}
	rivalryState, err := snapshot.LoadRivalries(rivalriesSnapshotPath)
	if err != nil {
		return err
	}

	standings, err := client.GetStandings(challengeID)
	if err != nil {
		return err
	}

	for _, standing := range standings {
		fmt.Printf("%d. %-12s %d\n", standing.Rank, standing.Name, standing.Points)
	}

	matches, err := client.GetMatches()
	if err != nil {
		return err
	}
	previousMatchesByID := make(map[string]matchmodel.Match, len(previousMatches))
	for _, match := range previousMatches {
		previousMatchesByID[match.MatchID] = match
	}
	for i := range matches {
		isLive := matches[i].Status != "" && matches[i].Status != "preMatch" && matches[i].Status != "fullTime"
		previousMatch, existed := previousMatchesByID[matches[i].MatchID]
		justEnded := existed && matches[i].Status == "fullTime" && previousMatch.Status != "fullTime"
		if !isLive && !justEnded {
			continue
		}
		events, err := client.GetMatchEvents(matches[i].MatchID)
		if err != nil {
			return err
		}
		matches[i].Events = events
	}
	var forecastHistory []forecasts.Forecast
	var messages []string
	updatedRivalries, rivalryOpportunities := rivalries.Update(standings, rivalryState)
	for _, opportunity := range rivalryOpportunities {
		message := banter.Generate(opportunity)
		fmt.Println(message)
		messages = append(messages, message)
	}
	for _, opportunity := range opportunities.DetectLiveUpdates(previousMatches, matches) {
		message := banter.Generate(opportunity)
		fmt.Println(message)
		messages = append(messages, message)
	}
	for _, match := range matches {
		fmt.Printf("%s %d-%d %s (%s)\n", match.HomeTeam, match.Score.Home, match.Score.Away, match.AwayTeam, match.Status)

		forecasts, err := client.GetForecasts(challengeID, match)
		if err != nil {
			return err
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
		if previousMatch, ok := previousMatchesByID[match.MatchID]; ok {
			for _, opportunity := range opportunities.DetectHeartbreaks(previousMatch, match, forecasts) {
				message := banter.Generate(opportunity)
				fmt.Println(message)
				messages = append(messages, message)
			}
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

	for _, message := range messages {
		if err := discordClient.Send(message); err != nil {
			return err
		}
	}

	if err := snapshot.SaveStandings(snapshotPath, standings); err != nil {
		return err
	}
	if err := snapshot.SaveMatches(matchesSnapshotPath, matches); err != nil {
		return err
	}
	if err := snapshot.SaveRivalries(rivalriesSnapshotPath, updatedRivalries); err != nil {
		return err
	}
	return nil
}
