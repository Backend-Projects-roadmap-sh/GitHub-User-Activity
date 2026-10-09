package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github-activity/internal/client"
	"github-activity/internal/format"
)

func main() {
	// Extra arguments are ignored; only the first one is used.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: github-activity <username>")
		os.Exit(1)
	}

	username := strings.TrimSpace(os.Args[1])
	if username == "" {
		fmt.Fprintln(os.Stderr, "error: username must not be empty")
		os.Exit(1)
	}

	if err := run(context.Background(), username); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, username string) error {
	events, err := client.New().FetchUserEvents(ctx, username)
	if err != nil {
		return err
	}

	formatter := format.New(useColor())

	if len(events) == 0 {
		fmt.Printf("%s has no recent public activity.\n", username)
		return nil
	}

	for _, e := range events {
		fmt.Println(formatter.Event(e))
	}

	return nil
}

// useColor reports whether ANSI colors should be emitted: only when stdout is
// a terminal and the NO_COLOR convention is not set.
func useColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}
