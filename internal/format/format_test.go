package format

import (
	"encoding/json"
	"testing"

	"github-activity/internal/model"
)

func TestEventPlain(t *testing.T) {
	f := New(false)
	tests := []struct {
		name string
		ev   model.Event
		want string
	}{
		{
			name: "push without commit count uses branch",
			ev:   model.Event{Type: model.EventTypePush, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"ref":"refs/heads/main"}`)},
			want: "Pushed to main in user/repo",
		},
		{
			name: "push with commit count",
			ev:   model.Event{Type: model.EventTypePush, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"ref":"refs/heads/main","size":3}`)},
			want: "Pushed 3 commits to main in user/repo",
		},
		{
			name: "single commit uses singular",
			ev:   model.Event{Type: model.EventTypePush, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"ref":"refs/heads/main","size":1}`)},
			want: "Pushed 1 commit to main in user/repo",
		},
		{
			name: "watch",
			ev:   model.Event{Type: model.EventTypeWatch, Repo: model.Repo{Name: "user/repo"}},
			want: "Starred user/repo",
		},
		{
			name: "fork",
			ev:   model.Event{Type: model.EventTypeFork, Repo: model.Repo{Name: "user/repo"}},
			want: "Forked user/repo",
		},
		{
			name: "create branch",
			ev:   model.Event{Type: model.EventTypeCreate, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"ref_type":"branch","ref":"feature"}`)},
			want: "Created branch feature in user/repo",
		},
		{
			name: "create repository",
			ev:   model.Event{Type: model.EventTypeCreate, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"ref_type":"repository"}`)},
			want: "Created repository user/repo",
		},
		{
			name: "issues opened",
			ev:   model.Event{Type: model.EventTypeIssues, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"action":"opened","issue":{"number":7}}`)},
			want: "Opened issue #7 in user/repo",
		},
		{
			name: "merged pull request",
			ev:   model.Event{Type: model.EventTypePullRequest, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"action":"closed","number":12,"pull_request":{"merged":true}}`)},
			want: "Merged pull request #12 in user/repo",
		},
		{
			name: "opened pull request",
			ev:   model.Event{Type: model.EventTypePullRequest, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"action":"opened","number":12,"pull_request":{"merged":false}}`)},
			want: "Opened pull request #12 in user/repo",
		},
		{
			name: "release",
			ev:   model.Event{Type: model.EventTypeRelease, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"action":"published","release":{"tag_name":"v1.0.0"}}`)},
			want: "Published release v1.0.0 in user/repo",
		},
		{
			name: "unknown event",
			ev:   model.Event{Type: "SomeNewEvent", Repo: model.Repo{Name: "user/repo"}},
			want: "Performed SomeNewEvent in user/repo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := f.Event(tc.ev); got != tc.want {
				t.Errorf("Event() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestEventColor(t *testing.T) {
	f := New(true)
	ev := model.Event{Type: model.EventTypePush, Repo: model.Repo{Name: "user/repo"}, Payload: json.RawMessage(`{"ref":"refs/heads/main"}`)}

	want := colorGreen + "Pushed" + colorReset + " to main in " + colorYellow + "user/repo" + colorReset
	if got := f.Event(ev); got != want {
		t.Errorf("Event() = %q, want %q", got, want)
	}
}
