package mpp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/DSanoussy/banter-engine/internal/forecasts"
	"github.com/DSanoussy/banter-engine/internal/matches"
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
	query := make(url.Values)
	query.Set("challengeId", challengeID)

	var apiResponse standingsResponse
	if err := c.get("/challenge-standings/users-standings", query, &apiResponse); err != nil {
		return nil, fmt.Errorf("fetch standings: %w", err)
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

func (c *Client) GetMatches() ([]matches.Match, error) {
	var apiMatches map[string]*matchDTO
	if err := c.get("/championships-current-matches", nil, &apiMatches); err != nil {
		return nil, fmt.Errorf("fetch matches: %w", err)
	}

	var apiClubs clubsResponse
	if err := c.get("/championship-clubs", nil, &apiClubs); err != nil {
		return nil, fmt.Errorf("fetch clubs: %w", err)
	}

	result := make([]matches.Match, 0, len(apiMatches))
	for id, match := range apiMatches {
		if match == nil {
			continue
		}

		matchID := match.MatchID
		if matchID == "" {
			matchID = id
		}

		result = append(result, matches.Match{
			MatchID:  matchID,
			Date:     match.Date,
			HomeTeam: clubName(apiClubs.ChampionshipClubs[match.Home.ClubID], match.Home.ClubID),
			AwayTeam: clubName(apiClubs.ChampionshipClubs[match.Away.ClubID], match.Away.ClubID),
			Score: matches.Score{
				Home: match.Home.Score,
				Away: match.Away.Score,
			},
			Status: match.Status,
			Quotations: matches.Quotations{
				Home: match.Quotations.Home,
				Draw: match.Quotations.Draw,
				Away: match.Quotations.Away,
			},
			PredictionStats: matches.PredictionStats{
				Home: match.Stats.Bets.Home,
				Draw: match.Stats.Bets.Draw,
				Away: match.Stats.Bets.Away,
			},
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Date.Equal(result[j].Date) {
			return result[i].MatchID < result[j].MatchID
		}
		return result[i].Date.Before(result[j].Date)
	})

	return result, nil
}

func (c *Client) GetForecasts(challengeID string, match matches.Match) ([]forecasts.Forecast, error) {
	path := fmt.Sprintf(
		"/user-match-forecasts/contest/%s/match/%s",
		url.PathEscape(challengeID),
		url.PathEscape(match.MatchID),
	)

	var apiForecasts map[string]forecastDTO
	if err := c.get(path, nil, &apiForecasts); err != nil {
		return nil, fmt.Errorf("fetch forecasts: %w", err)
	}

	result := make([]forecasts.Forecast, 0, len(apiForecasts))
	for userID, forecast := range apiForecasts {
		result = append(result, forecasts.Forecast{
			UserID:    userID,
			MatchID:   match.MatchID,
			MatchDate: match.Date,
			Prediction: matches.Score{
				Home: forecast.HomeScore,
				Away: forecast.AwayScore,
			},
			Result: match.Score,
			Points: forecast.Points.Total,
		})
	}

	return result, nil
}

func (c *Client) GetMatchEvents(matchID string) ([]matches.Event, error) {
	path := "/championship-match/" + url.PathEscape(matchID)
	var detail matchDetailDTO
	if err := c.get(path, nil, &detail); err != nil {
		return nil, fmt.Errorf("fetch match events: %w", err)
	}

	events := make([]matches.Event, 0, len(detail.EventsTimeline))
	for _, event := range detail.EventsTimeline {
		events = append(events, matches.Event{
			ID:       event.ID,
			Type:     event.Type,
			Detail:   eventDetail(event),
			Time:     event.Time,
			Side:     event.Side,
			PlayerID: event.PlayerID,
			Score: matches.Score{
				Home: event.Score.Home,
				Away: event.Score.Away,
			},
		})
	}
	return events, nil
}

func (c *Client) get(path string, query url.Values, target any) error {
	endpoint, err := url.Parse(baseURL + path)
	if err != nil {
		return fmt.Errorf("build URL: %w", err)
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
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

type matchDTO struct {
	MatchID string    `json:"matchId"`
	Date    time.Time `json:"date"`
	Home    struct {
		ClubID string `json:"clubId"`
		Score  int    `json:"score"`
	} `json:"home"`
	Away struct {
		ClubID string `json:"clubId"`
		Score  int    `json:"score"`
	} `json:"away"`
	Status     string        `json:"period"`
	Quotations outcomeValues `json:"quotations"`
	Stats      struct {
		Bets predictionStatsDTO `json:"bets"`
	} `json:"stats"`
}

type outcomeValues struct {
	Home int `json:"home"`
	Draw int `json:"draw"`
	Away int `json:"away"`
}

type predictionStatsDTO struct {
	Home float64 `json:"home"`
	Draw float64 `json:"draw"`
	Away float64 `json:"away"`
}

type forecastDTO struct {
	HomeScore int `json:"homeScore"`
	AwayScore int `json:"awayScore"`
	Points    struct {
		Total int `json:"total"`
	} `json:"points"`
}

type matchDetailDTO struct {
	EventsTimeline []matchEventDTO `json:"eventsTimeline"`
}

type matchEventDTO struct {
	ID          string `json:"eventId"`
	Type        string `json:"eventType"`
	GoalType    string `json:"goalType"`
	BookingType string `json:"bookingType"`
	Time        string `json:"time"`
	Side        string `json:"side"`
	PlayerID    string `json:"playerId"`
	Score       struct {
		Home int `json:"home"`
		Away int `json:"away"`
	} `json:"score"`
}

func eventDetail(event matchEventDTO) string {
	if event.GoalType != "" {
		return event.GoalType
	}
	return event.BookingType
}

type clubsResponse struct {
	ChampionshipClubs map[string]clubDTO `json:"championshipClubs"`
}

type clubDTO struct {
	Name      map[string]string `json:"name"`
	ShortName string            `json:"shortName"`
}

func clubName(club clubDTO, fallback string) string {
	if name := club.Name["fr-FR"]; name != "" {
		return name
	}
	if club.ShortName != "" {
		return club.ShortName
	}
	return fallback
}
