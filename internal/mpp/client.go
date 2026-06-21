package mpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/logging"
	"github.com/Sanssy/banter-engine/internal/matches"
	"github.com/Sanssy/banter-engine/internal/model"
)

const baseURL = "https://api.mpp.football"

type Client struct {
	token      string
	httpClient *http.Client
	logger     *logging.Logger
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logging.New(os.Stderr, "mpp"),
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

func (c *Client) GetMatches(challengeID string) ([]matches.Match, error) {
	c.logger.Info("match retrieval challenge_id=%s", challengeID)

	var challenge challengeDTO
	challengePath := "/challenge/" + url.PathEscape(challengeID)
	c.logMatchRoute(http.MethodGet, challengePath)
	if err := c.get(challengePath, nil, &challenge); err != nil {
		return nil, fmt.Errorf("fetch challenge: %w", err)
	}
	championshipID := challenge.GameSettings.ChampionshipID
	if championshipID <= 0 {
		return nil, fmt.Errorf("challenge %q has no championship", challengeID)
	}
	c.logger.Info("match retrieval championship_id=%d", championshipID)

	calendarPath := fmt.Sprintf("/championship-calendar/%d/nearest-game-weeks", championshipID)
	var calendar nearestGameWeeksResponse
	c.logMatchRoute(http.MethodGet, calendarPath)
	if err := c.get(calendarPath, nil, &calendar); err != nil {
		return nil, fmt.Errorf("fetch nearest game weeks: %w", err)
	}
	gameWeek, found := selectCurrentGameWeek(calendar.NearestGameWeeks, time.Now())
	if !found {
		return nil, fmt.Errorf("championship %d has no current game week", championshipID)
	}
	c.logger.Info(
		"match retrieval game_week=%d start=%s end=%s match_ids_count=%d match_ids_preview=%v",
		gameWeek.Number,
		gameWeek.Start.Format(time.RFC3339),
		gameWeek.End.Format(time.RFC3339),
		len(gameWeek.MatchIDs),
		previewStrings(gameWeek.MatchIDs, 10),
	)
	if len(gameWeek.MatchIDs) == 0 {
		return []matches.Match{}, nil
	}

	var apiMatches map[string]*matchDTO
	body := matchSummariesRequest{MatchIDs: gameWeek.MatchIDs}
	c.logMatchRoute(http.MethodPost, "/championship-match/summaries")
	if err := c.post("/championship-match/summaries", body, &apiMatches); err != nil {
		return nil, fmt.Errorf("fetch match summaries: %w", err)
	}
	firstRequestedID := gameWeek.MatchIDs[0]
	firstSummaryID := ""
	if firstSummary := apiMatches[firstRequestedID]; firstSummary != nil {
		firstSummaryID = firstSummary.MatchID
	}
	c.logger.Info(
		"match retrieval first_summary requested_match_id=%s summary_match_id=%s",
		firstRequestedID,
		firstSummaryID,
	)

	var apiClubs clubsResponse
	c.logMatchRoute(http.MethodGet, "/championship-clubs")
	if err := c.get("/championship-clubs", nil, &apiClubs); err != nil {
		return nil, fmt.Errorf("fetch clubs: %w", err)
	}

	result := make([]matches.Match, 0, len(gameWeek.MatchIDs))
	for index, id := range gameWeek.MatchIDs {
		match := apiMatches[id]
		if match == nil {
			if index < 5 {
				c.logger.Info("match retrieval summary requested_match_id=%s summary=nil", id)
			}
			continue
		}
		if index < 5 {
			c.logger.Info(
				"match retrieval summary requested_match_id=%s summary_match_id=%s",
				id,
				match.MatchID,
			)
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
	for index, match := range result {
		if index >= 5 {
			break
		}
		c.logger.Info(
			"match retrieval resolved index=%d match_id=%s home_team=%s away_team=%s date=%s",
			index,
			match.MatchID,
			match.HomeTeam,
			match.AwayTeam,
			match.Date.Format(time.RFC3339),
		)
	}

	return result, nil
}

func (c *Client) logMatchRoute(method, path string) {
	c.logger.Info("match retrieval route=%s %s", method, path)
}

func previewStrings(values []string, limit int) []string {
	if len(values) <= limit {
		return values
	}
	return values[:limit]
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
	return c.do(req, target)
}

func (c *Client) post(path string, body, target any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, target)
}

func (c *Client) do(req *http.Request, target any) error {
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

type challengeDTO struct {
	GameSettings struct {
		ChampionshipID int `json:"championshipId"`
	} `json:"gameSettings"`
}

type nearestGameWeeksResponse struct {
	NearestGameWeeks nearestGameWeeksDTO `json:"nearestGameWeeks"`
}

type nearestGameWeeksDTO struct {
	PreviousGameWeek *gameWeekDTO `json:"previousGameWeek"`
	CurrentGameWeek  *gameWeekDTO `json:"currentGameWeek"`
	NextGameWeek     *gameWeekDTO `json:"nextGameWeek"`
}

type gameWeekDTO struct {
	Number   int       `json:"gameWeekNumber"`
	MatchIDs []string  `json:"matchesIds"`
	Start    time.Time `json:"startDate"`
	End      time.Time `json:"endDate"`
}

type matchSummariesRequest struct {
	MatchIDs []string `json:"matchesIds"`
}

func selectCurrentGameWeek(nearest nearestGameWeeksDTO, now time.Time) (gameWeekDTO, bool) {
	if nearest.CurrentGameWeek != nil {
		return *nearest.CurrentGameWeek, true
	}
	for _, gameWeek := range []*gameWeekDTO{nearest.PreviousGameWeek, nearest.NextGameWeek} {
		if gameWeek != nil && !gameWeek.Start.IsZero() && !gameWeek.End.IsZero() &&
			!now.Before(gameWeek.Start) && !now.After(gameWeek.End) {
			return *gameWeek, true
		}
	}
	if nearest.NextGameWeek != nil {
		return *nearest.NextGameWeek, true
	}
	if nearest.PreviousGameWeek != nil {
		return *nearest.PreviousGameWeek, true
	}
	return gameWeekDTO{}, false
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
