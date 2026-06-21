package main

import (
	"fmt"
	"log"
	"os"

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

	for _, opportunity := range opportunities.DetectRankingOvertakes(previousStandings, standings) {
		fmt.Printf("%s: %s overtook %s\n", opportunity.Type, opportunity.Actor, opportunity.Target)
	}

	if err := snapshot.SaveStandings(snapshotPath, standings); err != nil {
		log.Fatal(err)
	}
}
