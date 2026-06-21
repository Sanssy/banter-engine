package mpp

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/DSanoussy/banter-engine/internal/matches"
)

func TestGetMatches(t *testing.T) {
	client := NewClient("test-token")
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("Authorization header = %q, want %q", got, "Bearer test-token")
		}

		var body string
		switch req.URL.Path {
		case "/championships-current-matches":
			body = `{
				"match-1": {
					"matchId": "match-1",
					"period": "preMatch",
					"home": {"clubId": "club-1", "score": 0},
					"away": {"clubId": "club-2", "score": 0},
					"quotations": {"home": 125, "draw": 310, "away": 240},
					"stats": {"bets": {"home": 0.5, "draw": 0.2, "away": 0.3}}
				},
				"match-without-data": null
			}`
		case "/championship-clubs":
			body = `{
				"championshipClubs": {
					"club-1": {"name": {"fr-FR": "France"}},
					"club-2": {"shortName": "Espagne"}
				}
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
			MatchID:  "match-1",
			HomeTeam: "France",
			AwayTeam: "Espagne",
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

	got, err := client.GetMatches()
	if err != nil {
		t.Fatalf("GetMatches() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetMatches() = %#v, want %#v", got, want)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
