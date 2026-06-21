package mpp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DSanoussy/banter-engine/internal/model"
)

const baseURL = "https://api.mpp.football"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetStandings(challengeID string) ([]model.Standing, error) {
	endpoint, err := url.Parse(baseURL + "/challenge-standings/users-standings")
	if err != nil {
		return nil, fmt.Errorf("build standings URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("challengeId", challengeID)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create standings request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch standings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch standings: unexpected HTTP status %s", resp.Status)
	}

	var apiResponse standingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("decode standings response: %w", err)
	}

	standings := make([]model.Standing, 0, len(apiResponse.Standings))
	for _, standing := range apiResponse.Standings {
		standings = append(standings, model.Standing{
			UserID: standing.User.ID,
			Name:   standing.User.Username,
			Rank:   standing.Ranking.Rank,
			Points: standing.Ranking.Points,
		})
	}

	return standings, nil
}

type standingsResponse struct {
	Standings []standingDTO `json:"standings"`
}

type standingDTO struct {
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
	Ranking struct {
		Rank   int `json:"rank"`
		Points int `json:"points"`
	} `json:"ranking"`
}
