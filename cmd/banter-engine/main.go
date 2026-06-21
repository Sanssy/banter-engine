package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DSanoussy/banter-engine/internal/banter"
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
	for _, match := range matches {
		fmt.Printf("%s %d-%d %s (%s)\n", match.HomeTeam, match.Score.Home, match.Score.Away, match.AwayTeam, match.Status)

		forecasts, err := client.GetForecasts(challengeID, match)
		if err != nil {
			log.Fatal(err)
		}
		for _, forecast := range forecasts {
			fmt.Printf("  %s: %d-%d (%d points)\n", forecast.UserID, forecast.Prediction.Home, forecast.Prediction.Away, forecast.Points)
		}
	}

	for _, opportunity := range opportunities.Detect(previousStandings, standings) {
		fmt.Println(banter.Generate(opportunity))
	}

	if err := snapshot.SaveStandings(snapshotPath, standings); err != nil {
		log.Fatal(err)
	}
}
