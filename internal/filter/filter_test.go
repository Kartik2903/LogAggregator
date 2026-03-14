package filter

// LAB 8: Unit tests for the Filter module
// Demonstrates: table-driven tests, interface-based testing, test helpers

import (
	"testing"
	"time"

	"logmerge/internal/models"
)

// LAB 8: Helper function to create test events
func makeEvent(level models.LogLevel, message string) models.LogEvent {
	return models.LogEvent{
		Timestamp: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC),
		Level:     level,
		Message:   message,
		Source:    "test.log",
	}
}

// LAB 8: Table-driven test for LevelFilter
func TestLevelFilter(t *testing.T) {
	tests := []struct {
		name     string
		minLevel models.LogLevel
		event    models.LogEvent
		want     bool
	}{
		{"INFO passes INFO filter", models.INFO, makeEvent(models.INFO, "hello"), true},
		{"WARN passes INFO filter", models.INFO, makeEvent(models.WARN, "hello"), true},
		{"ERROR passes INFO filter", models.INFO, makeEvent(models.ERROR, "hello"), true},
		{"INFO blocked by WARN filter", models.WARN, makeEvent(models.INFO, "hello"), false},
		{"WARN passes WARN filter", models.WARN, makeEvent(models.WARN, "hello"), true},
		{"ERROR passes WARN filter", models.WARN, makeEvent(models.ERROR, "hello"), true},
		{"INFO blocked by ERROR filter", models.ERROR, makeEvent(models.INFO, "hello"), false},
		{"WARN blocked by ERROR filter", models.ERROR, makeEvent(models.WARN, "hello"), false},
		{"ERROR passes ERROR filter", models.ERROR, makeEvent(models.ERROR, "hello"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := LevelFilter{MinLevel: tt.minLevel}
			got := f.Match(tt.event)
			if got != tt.want {
				t.Errorf("LevelFilter(%v).Match(%v) = %v, want %v",
					tt.minLevel, tt.event.Level, got, tt.want)
			}
		})
	}
}

// LAB 8: Table-driven test for KeywordFilter
func TestKeywordFilter(t *testing.T) {
	tests := []struct {
		name    string
		keyword string
		message string
		want    bool
	}{
		{"keyword found", "timeout", "Database connection timeout", true},
		{"keyword case insensitive", "TIMEOUT", "Database connection timeout", true},
		{"keyword not found", "banana", "Database connection timeout", false},
		{"empty keyword matches all", "", "any message", true},
		{"partial match", "connect", "Database connection timeout", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := KeywordFilter{Keyword: tt.keyword}
			event := makeEvent(models.INFO, tt.message)
			got := f.Match(event)
			if got != tt.want {
				t.Errorf("KeywordFilter(%q).Match(%q) = %v, want %v",
					tt.keyword, tt.message, got, tt.want)
			}
		})
	}
}

// LAB 8: Test for RegexFilter
func TestRegexFilter(t *testing.T) {
	// LAB 5: Error handling — test invalid regex
	_, err := NewRegexFilter("[invalid")
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}

	// Valid regex tests
	f, err := NewRegexFilter(`timeout|failed`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{"matches timeout", "connection timeout occurred", true},
		{"matches failed", "operation failed", true},
		{"no match", "all systems operational", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := makeEvent(models.ERROR, tt.message)
			got := f.Match(event)
			if got != tt.want {
				t.Errorf("RegexFilter.Match(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}

// LAB 8: Test for CompositeFilter (ChainFilters)
func TestCompositeFilter(t *testing.T) {
	level := LevelFilter{MinLevel: models.WARN}
	keyword := KeywordFilter{Keyword: "timeout"}
	combined := ChainFilters(level, keyword)

	tests := []struct {
		name    string
		level   models.LogLevel
		message string
		want    bool
	}{
		{"WARN + keyword match", models.WARN, "connection timeout", true},
		{"ERROR + keyword match", models.ERROR, "request timeout", true},
		{"INFO blocked by level", models.INFO, "connection timeout", false},
		{"WARN but no keyword", models.WARN, "all good", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := makeEvent(tt.level, tt.message)
			got := combined.Match(event)
			if got != tt.want {
				t.Errorf("ChainFilters.Match(%v, %q) = %v, want %v",
					tt.level, tt.message, got, tt.want)
			}
		})
	}
}

// LAB 8: Test for ApplyFilter utility function
func TestApplyFilter(t *testing.T) {
	events := []models.LogEvent{
		makeEvent(models.INFO, "info message"),
		makeEvent(models.WARN, "warning message"),
		makeEvent(models.ERROR, "error message"),
		makeEvent(models.INFO, "another info"),
	}

	f := LevelFilter{MinLevel: models.WARN}
	result := ApplyFilter(events, f)

	if len(result) != 2 {
		t.Errorf("ApplyFilter: expected 2 events, got %d", len(result))
	}

	for _, event := range result {
		if event.Level < models.WARN {
			t.Errorf("ApplyFilter: unexpected level %v in results", event.Level)
		}
	}
}
