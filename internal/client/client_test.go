package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchUserEvents(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/users/octocat/events"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Error("User-Agent header is empty")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"type":"WatchEvent","repo":{"name":"octocat/Hello-World"}}]`))
	}))
	defer srv.Close()

	events, err := NewWithBaseURL(srv.URL).FetchUserEvents(context.Background(), "octocat")
	if err != nil {
		t.Fatalf("FetchUserEvents() error = %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}

	if events[0].Type != "WatchEvent" || events[0].Repo.Name != "octocat/Hello-World" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestFetchUserEventsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := NewWithBaseURL(srv.URL).FetchUserEvents(context.Background(), "nobody")
	if err == nil {
		t.Fatal("expected an error for a 404 response")
	}
}

func TestFetchUserEventsRateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	_, err := NewWithBaseURL(srv.URL).FetchUserEvents(context.Background(), "octocat")
	if err == nil {
		t.Fatal("expected an error for a rate-limited response")
	}
}
