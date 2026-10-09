// Package format renders GitHub activity events as human-readable text.
package format

import (
	"encoding/json"
	"fmt"
	"strings"

	"github-activity/internal/model"
)

// ANSI color escapes.
const (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

// Formatter renders events, optionally with ANSI colors.
type Formatter struct {
	color bool
}

// New returns a Formatter. When color is true, action keywords are rendered
// green and repository names yellow.
func New(color bool) *Formatter {
	return &Formatter{color: color}
}

// Event renders a single event as a human-readable sentence.
func (f *Formatter) Event(e model.Event) string {
	repoName := f.repo(e.Repo.Name)

	switch e.Type {
	case model.EventTypePush:
		var p struct {
			Ref     string     `json:"ref"`
			Size    int        `json:"size"`
			Commits []struct{} `json:"commits"`
		}
		_ = json.Unmarshal(e.Payload, &p)

		branch := strings.TrimPrefix(p.Ref, "refs/heads/")
		n := len(p.Commits)
		if n == 0 {
			n = p.Size
		}

		// The public events feed omits commit counts, so fall back to the branch.
		if n > 0 {
			return fmt.Sprintf("%s %d %s to %s in %s", f.action("Pushed"), n, plural(n, "commit", "commits"), branch, repoName)
		}

		if branch != "" {
			return fmt.Sprintf("%s to %s in %s", f.action("Pushed"), branch, repoName)
		}

		return fmt.Sprintf("%s to %s", f.action("Pushed"), repoName)

	case model.EventTypeWatch:
		return fmt.Sprintf("%s %s", f.action("Starred"), repoName)

	case model.EventTypeFork:
		return fmt.Sprintf("%s %s", f.action("Forked"), repoName)

	case model.EventTypePublic:
		return fmt.Sprintf("%s %s public", f.action("Made"), repoName)

	case model.EventTypeCreate:
		var p struct {
			RefType string `json:"ref_type"`
			Ref     string `json:"ref"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		if p.RefType == "repository" || p.RefType == "" {
			return fmt.Sprintf("%s repository %s", f.action("Created"), repoName)
		}

		return fmt.Sprintf("%s %s %s in %s", f.action("Created"), p.RefType, p.Ref, repoName)

	case model.EventTypeDelete:
		var p struct {
			RefType string `json:"ref_type"`
			Ref     string `json:"ref"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s %s %s in %s", f.action("Deleted"), p.RefType, p.Ref, repoName)

	case model.EventTypeIssues:
		var p struct {
			Action string `json:"action"`
			Issue  struct {
				Number int `json:"number"`
			} `json:"issue"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s issue #%d in %s", f.action(capitalize(p.Action)), p.Issue.Number, repoName)

	case model.EventTypeIssueComment:
		var p struct {
			Action string `json:"action"`
			Issue  struct {
				Number int `json:"number"`
			} `json:"issue"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s a comment on issue #%d in %s", f.action(capitalize(p.Action)), p.Issue.Number, repoName)

	case model.EventTypePullRequest:
		var p struct {
			Action      string `json:"action"`
			Number      int    `json:"number"`
			PullRequest struct {
				Merged bool `json:"merged"`
			} `json:"pull_request"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		if p.Action == "closed" && p.PullRequest.Merged {
			return fmt.Sprintf("%s pull request #%d in %s", f.action("Merged"), p.Number, repoName)
		}

		return fmt.Sprintf("%s pull request #%d in %s", f.action(capitalize(p.Action)), p.Number, repoName)

	case model.EventTypePullRequestReview:
		var p struct {
			Action      string `json:"action"`
			PullRequest struct {
				Number int `json:"number"`
			} `json:"pull_request"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s a review on pull request #%d in %s", f.action(capitalize(p.Action)), p.PullRequest.Number, repoName)

	case model.EventTypePullRequestReviewComment:
		var p struct {
			PullRequest struct {
				Number int `json:"number"`
			} `json:"pull_request"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s on pull request #%d in %s", f.action("Commented"), p.PullRequest.Number, repoName)

	case model.EventTypeCommitComment:
		return fmt.Sprintf("%s on a commit in %s", f.action("Commented"), repoName)

	case model.EventTypeRelease:
		var p struct {
			Action  string `json:"action"`
			Release struct {
				TagName string `json:"tag_name"`
			} `json:"release"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s release %s in %s", f.action(capitalize(p.Action)), p.Release.TagName, repoName)

	case model.EventTypeMember:
		var p struct {
			Action string `json:"action"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s a collaborator to %s", f.action(capitalize(p.Action)), repoName)

	case model.EventTypeGollum:
		return fmt.Sprintf("%s the wiki in %s", f.action("Updated"), repoName)

	default:
		return fmt.Sprintf("%s %s in %s", f.action("Performed"), e.Type, repoName)
	}
}

// action styles an action keyword in green.
func (f *Formatter) action(s string) string {
	if !f.color {
		return s
	}

	return colorGreen + s + colorReset
}

// repo styles a repository name in yellow.
func (f *Formatter) repo(s string) string {
	if !f.color {
		return s
	}

	return colorYellow + s + colorReset
}

func plural(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}

	return plural
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
