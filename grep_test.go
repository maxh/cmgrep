package main

import (
	"strings"
	"testing"
)

func TestCountMatchesCountsMatchingLines(t *testing.T) {
	input := "INFO startup\nERROR disk full\nERROR timeout\n"

	got, err := countMatches(strings.NewReader(input), "ERROR", false, false)
	if err != nil {
		t.Fatalf("countMatches returned error: %v", err)
	}
	if got != 2 {
		t.Errorf("countMatches = %d, want 2", got)
	}
}

func TestCountMatchesReturnsZeroForNoMatch(t *testing.T) {
	input := "INFO startup\nWARN retry\n"

	got, err := countMatches(strings.NewReader(input), "ERROR", false, false)
	if err != nil {
		t.Fatalf("countMatches returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("countMatches = %d, want 0", got)
	}
}

func TestCountMatchesIgnoresCaseWhenRequested(t *testing.T) {
	input := "ERROR first\nerror second\nInfo third\n"

	got, err := countMatches(strings.NewReader(input), "error", true, false)
	if err != nil {
		t.Fatalf("countMatches returned error: %v", err)
	}
	if got != 2 {
		t.Errorf("countMatches = %d, want 2", got)
	}
}

func TestCountMatchesUsesExtendedRegularExpressions(t *testing.T) {
	input := "ERROR disk full\nWARN retrying\nINFO ready\n"

	got, err := countMatches(
		strings.NewReader(input),
		"ERROR|WARN",
		false,
		true,
	)
	if err != nil {
		t.Fatalf("countMatches returned error: %v", err)
	}
	if got != 2 {
		t.Errorf("countMatches = %d, want 2", got)
	}
}

func TestCountMatchesRejectsInvalidRegularExpression(t *testing.T) {
	input := "ERROR disk full\n"

	got, err := countMatches(strings.NewReader(input), "[", false, false)
	if err == nil {
		t.Fatal("countMatches returned nil error for invalid regular expression")
	}
	if got != 0 {
		t.Errorf("countMatches = %d, want 0 when grep fails", got)
	}
}

func TestCountMatchesAcceptsPatternBeginningWithHyphen(t *testing.T) {
	input := "-n\nordinary line\nanother -n value\n"

	got, err := countMatches(strings.NewReader(input), "-n", false, false)
	if err != nil {
		t.Fatalf("countMatches returned error: %v", err)
	}
	if got != 2 {
		t.Errorf("countMatches = %d, want 2", got)
	}
}

func TestCountMatchesIsCaseSensitiveByDefault(t *testing.T) {
	input := "ERROR uppercase\nerror lowercase\n"

	got, err := countMatches(strings.NewReader(input), "error", false, false)
	if err != nil {
		t.Fatalf("countMatches returned error: %v", err)
	}
	if got != 1 {
		t.Errorf("countMatches = %d, want 1", got)
	}
}

func TestCountMatchesUsesBasicRegexByDefault(t *testing.T) {
	input := "ERROR disk full\nWARN retrying\nINFO ready\n"

	got, err := countMatches(
		strings.NewReader(input),
		"ERROR|WARN",
		false,
		false,
	)
	if err != nil {
		t.Fatalf("countMatches returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("countMatches = %d, want 0 without extended regex", got)
	}
}
