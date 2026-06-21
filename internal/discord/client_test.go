package discord

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSendRetriesOnce(t *testing.T) {
	attempts := 0
	client := NewClient("https://discord.example/webhook")
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return response(http.StatusInternalServerError), nil
		}

		var payload struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload.Content != "test message" {
			t.Fatalf("content = %q, want %q", payload.Content, "test message")
		}
		return response(http.StatusNoContent), nil
	})

	if err := client.Send("test message"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestSendReturnsErrorAfterRetry(t *testing.T) {
	attempts := 0
	client := NewClient("https://discord.example/webhook")
	client.httpClient.Transport = roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		attempts++
		return response(http.StatusInternalServerError), nil
	})

	if err := client.Send("test message"); err == nil {
		t.Fatal("Send() error = nil, want an error")
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func response(statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}
}
