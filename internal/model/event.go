// Package model contains the domain types for GitHub activity events.
package model

import (
	"encoding/json"
	"time"
)

// EventType is the type of a GitHub activity event.
type EventType string

// Known GitHub event types.
const (
	EventTypePush                     EventType = "PushEvent"
	EventTypeWatch                    EventType = "WatchEvent"
	EventTypeFork                     EventType = "ForkEvent"
	EventTypePublic                   EventType = "PublicEvent"
	EventTypeCreate                   EventType = "CreateEvent"
	EventTypeDelete                   EventType = "DeleteEvent"
	EventTypeIssues                   EventType = "IssuesEvent"
	EventTypeIssueComment             EventType = "IssueCommentEvent"
	EventTypePullRequest              EventType = "PullRequestEvent"
	EventTypePullRequestReview        EventType = "PullRequestReviewEvent"
	EventTypePullRequestReviewComment EventType = "PullRequestReviewCommentEvent"
	EventTypeCommitComment            EventType = "CommitCommentEvent"
	EventTypeRelease                  EventType = "ReleaseEvent"
	EventTypeMember                   EventType = "MemberEvent"
	EventTypeGollum                   EventType = "GollumEvent"
)

// Event is a single public activity event of a user.
type Event struct {
	Type      EventType       `json:"type"`
	Repo      Repo            `json:"repo"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
