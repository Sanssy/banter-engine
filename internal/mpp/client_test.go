package mpp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Sanssy/banter-engine/internal/forecasts"
	"github.com/Sanssy/banter-engine/internal/logging"
	"github.com/Sanssy/banter-engine/internal/matches"
	"github.com/Sanssy/banter-engine/internal/references"
)

func TestGetMatchesUsesChallengeCurrentGameWeek(t *testing.T) {
	client := NewClient("test-token", references.New(io.Discard))
	var logOutput bytes.Buffer
	client.logger = logging.New(&logOutput, "mpp")
	var requestedPaths []string
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requestedPaths = append(requestedPaths, req.URL.Path)
		if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("Authorization header = %q, want %q", got, "Bearer test-token")
		}

		var body string
		switch req.URL.Path {
		case "/challenge/challenge-1":
			body = `{"gameSettings":{"championshipId":8}}`
		case "/championship-calendar/8/nearest-game-weeks":
			body = `{
				"nearestGameWeeks": {
					"previousGameWeek": {"gameWeekNumber":1,"matchesIds":["club-match"]},
					"currentGameWeek": {
						"gameWeekNumber":2,
						"matchesIds":["world-cup-match"],
						"startDate":"2026-06-18T16:00:00Z",
						"endDate":"2026-06-24T02:00:00Z"
					},
					"nextGameWeek": {"gameWeekNumber":3,"matchesIds":["future-match"]}
				}
			}`
		case "/championship-match/summaries":
			if req.Method != http.MethodPost {
				t.Fatalf("request method = %q, want POST", req.Method)
			}
			var request matchSummariesRequest
			if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
				t.Fatalf("decode summaries request: %v", err)
			}
			if !reflect.DeepEqual(request.MatchIDs, []string{"world-cup-match"}) {
				t.Fatalf("matchesIds = %#v, want current game week IDs", request.MatchIDs)
			}
			body = `{
				"world-cup-match": {
					"matchId": "world-cup-match",
					"championshipId": 8,
					"gameWeekNumber": 2,
					"date": "2026-06-21T16:00:00Z",
					"period": "preMatch",
					"home": {"clubId": "mpp_championship_club_367", "score": 0},
					"away": {"clubId": "mpp_championship_club_522", "score": 0},
					"quotations": {"home": 125, "draw": 310, "away": 240},
					"stats": {"bets": {"home": 0.5, "draw": 0.2, "away": 0.3}}
				},
				"club-match": {
					"matchId": "club-match",
					"home": {"clubId": "mpp_championship_club_152"},
					"away": {"clubId": "mpp_championship_club_1430"}
				}
			}`
		case "/championship-clubs":
			body = `{
				"mpp_championship_club_367": {"lang": {"fr-FR": {"name": "Canada", "shortName": "CAN"}}},
				"mpp_championship_club_522": {"shortName": "Qatar"},
				"mpp_championship_club_152": {"shortName": "AS Saint-Étienne"},
				"mpp_championship_club_1430": {"shortName": "Guingamp"}
			}`
		default:
			t.Fatalf("unexpected request path %q", req.URL.Path)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	want := []matches.Match{
		{
			MatchID:  "world-cup-match",
			Date:     time.Date(2026, 6, 21, 16, 0, 0, 0, time.UTC),
			HomeTeam: "Canada",
			AwayTeam: "Qatar",
			Score:    matches.Score{},
			Status:   "preMatch",
			Quotations: matches.Quotations{
				Home: 125,
				Draw: 310,
				Away: 240,
			},
			PredictionStats: matches.PredictionStats{
				Home: 0.5,
				Draw: 0.2,
				Away: 0.3,
			},
		},
	}

	got, err := client.GetMatches("challenge-1")
	if err != nil {
		t.Fatalf("GetMatches() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetMatches() = %#v, want %#v", got, want)
	}
	if len(requestedPaths) != 4 {
		t.Fatalf("requested paths = %#v, want challenge-scoped retrieval only", requestedPaths)
	}
	for _, expected := range []string{
		"challenge_id=challenge-1",
		"route=GET /challenge/challenge-1",
		"championship_id=8",
		"route=GET /championship-calendar/8/nearest-game-weeks",
		"game_week=2",
		"match_ids_count=1",
		"match_ids_preview=[world-cup-match]",
		"route=POST /championship-match/summaries",
		"first_summary requested_match_id=world-cup-match summary_match_id=world-cup-match",
		"requested_match_id=world-cup-match summary_match_id=world-cup-match",
		"route=GET /championship-clubs",
		"club reference decoded clubs_count=4",
		"match_id=world-cup-match home_team=Canada away_team=Qatar date=2026-06-21T16:00:00Z",
	} {
		if !strings.Contains(logOutput.String(), expected) {
			t.Errorf("diagnostic log does not contain %q:\n%s", expected, logOutput.String())
		}
	}
}

func TestClubsResponseAcceptsWrappedReference(t *testing.T) {
	data := []byte(`{"championshipClubs":{"club-1":{"name":{"fr-FR":"France"}}}}`)
	var response clubsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got := clubDisplayName(response.ChampionshipClubs["club-1"]); got != "France" {
		t.Fatalf("clubDisplayName() = %q, want France", got)
	}
}

func TestSelectCurrentGameWeekUsesDateWhenCurrentIsOmitted(t *testing.T) {
	now := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
	previous := &gameWeekDTO{
		Number:   2,
		MatchIDs: []string{"world-cup-match"},
		Start:    now.Add(-24 * time.Hour),
		End:      now.Add(24 * time.Hour),
	}
	next := &gameWeekDTO{Number: 3, MatchIDs: []string{"future-match"}, Start: now.Add(48 * time.Hour), End: now.Add(72 * time.Hour)}

	got, found := selectCurrentGameWeek(nearestGameWeeksDTO{PreviousGameWeek: previous, NextGameWeek: next}, now)
	if !found || got.Number != 2 {
		t.Fatalf("selectCurrentGameWeek() = %+v, %v, want in-progress game week", got, found)
	}
}

func TestGetForecasts(t *testing.T) {
	resolver := references.New(io.Discard)
	resolver.RegisterUser("user-1", "LeDaveCoinCoin")
	client := NewClient("test-token", resolver)
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		wantPath := "/user-match-forecasts/contest/challenge-1/match/match-1"
		if req.URL.Path != wantPath {
			t.Fatalf("request path = %q, want %q", req.URL.Path, wantPath)
		}

		body := `{
			"user-1": {
				"homeScore": 2,
				"awayScore": 1,
				"points": {"base": 3, "exact": 2, "total": 5}
			}
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	match := matches.Match{MatchID: "match-1", Score: matches.Score{Home: 2, Away: 1}}
	want := []forecasts.Forecast{
		{
			UserID:     "user-1",
			UserName:   "LeDaveCoinCoin",
			MatchID:    "match-1",
			Prediction: matches.Score{Home: 2, Away: 1},
			Result:     matches.Score{Home: 2, Away: 1},
			Points:     5,
		},
	}

	got, err := client.GetForecasts("challenge-1", match)
	if err != nil {
		t.Fatalf("GetForecasts() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetForecasts() = %#v, want %#v", got, want)
	}
}

func TestGetMatchEvents(t *testing.T) {
	client := NewClient("test-token", references.New(io.Discard))
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/championship-match/match-1" {
			t.Fatalf("request path = %q, want %q", req.URL.Path, "/championship-match/match-1")
		}
		body := `{
			"eventsTimeline": [{
				"eventId": "event-1",
				"eventType": "goal",
				"goalType": "penalty",
				"time": "42'",
				"side": "home",
				"playerId": "player-1",
				"score": {"home": 1, "away": 0}
			}]
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	want := []matches.Event{
		{
			ID:       "event-1",
			Type:     "goal",
			Detail:   "penalty",
			Time:     "42'",
			Side:     "home",
			PlayerID: "player-1",
			Score:    matches.Score{Home: 1},
		},
	}
	got, err := client.GetMatchEvents("match-1")
	if err != nil {
		t.Fatalf("GetMatchEvents() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetMatchEvents() = %#v, want %#v", got, want)
	}
}

func TestGetMatchEventsDecodesStringScore(t *testing.T) {
	client := NewClient("test-token", references.New(io.Discard))
	var logOutput bytes.Buffer
	client.logger = logging.New(&logOutput, "mpp")
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body := `{
			"eventsTimeline": [{
				"eventId": "event-1",
				"eventType": "goal",
				"time": "90' +2",
				"side": "away",
				"score": "1 - 2"
			}]
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	events, err := client.GetMatchEvents("mpp_championship_match_2608265")
	if err != nil {
		t.Fatalf("GetMatchEvents() error = %v", err)
	}
	if len(events) != 1 || events[0].Score != (matches.Score{Home: 1, Away: 2}) {
		t.Fatalf("GetMatchEvents() = %#v, want string score decoded as 1-2", events)
	}
	for _, expected := range []string{
		"url=https://api.mpp.football/championship-match/mpp_championship_match_2608265",
		"body_preview=",
		`raw=\"1 - 2\"`,
	} {
		if !strings.Contains(logOutput.String(), expected) {
			t.Errorf("diagnostic log does not contain %q:\n%s", expected, logOutput.String())
		}
	}
}

func TestEventScoreDTOAcceptsSupportedFormats(t *testing.T) {
	tests := []struct {
		payload string
		want    eventScoreDTO
	}{
		{payload: `{"home":2,"away":1}`, want: eventScoreDTO{Home: 2, Away: 1}},
		{payload: `"2-1"`, want: eventScoreDTO{Home: 2, Away: 1}},
		{payload: `"2 : 1"`, want: eventScoreDTO{Home: 2, Away: 1}},
	}
	for _, test := range tests {
		var got eventScoreDTO
		if err := json.Unmarshal([]byte(test.payload), &got); err != nil {
			t.Fatalf("Unmarshal(%s) error = %v", test.payload, err)
		}
		if got.Home != test.want.Home || got.Away != test.want.Away {
			t.Errorf("Unmarshal(%s) = %d-%d, want %d-%d", test.payload, got.Home, got.Away, test.want.Home, test.want.Away)
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
