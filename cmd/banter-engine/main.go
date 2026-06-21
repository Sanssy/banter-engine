package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DSanoussy/banter-engine/internal/mpp"
)

const challengeID = "mpp_challenge_UDKDDH27"

func main() {
	token := os.Getenv("MPP_TOKEN")
	if token == "" {
		log.Fatal("MPP_TOKEN environment variable is required")
	}

	client := mpp.NewClient(token)
	standings, err := client.GetStandings(challengeID)
	if err != nil {
		log.Fatal(err)
	}

	for _, standing := range standings {
		fmt.Printf("%d. %-12s %d\n", standing.Rank, standing.Name, standing.Points)
	}
}
