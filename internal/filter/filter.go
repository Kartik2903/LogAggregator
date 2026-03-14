package filter

// LAB 6 (Interface) + LAB 5 (Functions & Error Handling):
// Filter module — filters logs by level, keyword, or regex.
// Uses interfaces to make the filtering system extensible (pluggable filters).

import (
	"regexp"
	"strings"

	"logmerge/internal/models"
)

// ============================================================
// LAB 6: Filter Interface — extensible, pluggable filter system
// Any type implementing Match() can act as a log filter.
// ============================================================
type Filter interface {
	Match(event models.LogEvent) bool // LAB 6: interface method
}

// ============================================================
// LevelFilter — filters events by minimum log level
// LAB 5: Function with conditional logic
// LAB 6: Implements Filter interface
// ============================================================
type LevelFilter struct {
	MinLevel models.LogLevel // LAB 1: custom type field
}

// LAB 6: Method set — LevelFilter satisfies Filter interface
func (f LevelFilter) Match(event models.LogEvent) bool {
	// LAB 2: conditional — check if event level meets minimum threshold
	return event.Level.Enabled(f.MinLevel)
}

// ============================================================
// KeywordFilter — filters events whose message contains a keyword
// LAB 5: String manipulation functions
// LAB 6: Implements Filter interface
// ============================================================
type KeywordFilter struct {
	Keyword string // The keyword to search for (case-insensitive)
}

// LAB 6: Method set — KeywordFilter satisfies Filter interface
func (f KeywordFilter) Match(event models.LogEvent) bool {
	// LAB 5: strings.Contains for substring matching (case-insensitive)
	return strings.Contains(
		strings.ToLower(event.Message),
		strings.ToLower(f.Keyword),
	)
}

// ============================================================
// RegexFilter — filters events whose message matches a regex pattern
// LAB 5: Error handling in regex compilation
// LAB 6: Implements Filter interface
// ============================================================
type RegexFilter struct {
	Pattern *regexp.Regexp // LAB 7: pointer to compiled regex
}

// LAB 5: Constructor with error return for invalid regex
func NewRegexFilter(pattern string) (*RegexFilter, error) {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err // LAB 5: propagate error to caller
	}
	return &RegexFilter{Pattern: compiled}, nil // LAB 7: return pointer
}

// LAB 6: Method set — RegexFilter satisfies Filter interface
func (f *RegexFilter) Match(event models.LogEvent) bool {
	// LAB 5: use compiled regex to match message
	return f.Pattern.MatchString(event.Message)
}

// ============================================================
// CompositeFilter — chains multiple filters with AND logic
// LAB 3: Slice of interfaces
// LAB 6: Demonstrates polymorphism — any Filter impl works here
// ============================================================
type CompositeFilter struct {
	Filters []Filter // LAB 3: slice of interface values
}

// ChainFilters creates a composite filter that requires ALL sub-filters to match.
func ChainFilters(filters ...Filter) Filter {
	return &CompositeFilter{Filters: filters}
}

// LAB 6: CompositeFilter also satisfies the Filter interface (recursive composition)
func (cf *CompositeFilter) Match(event models.LogEvent) bool {
	// LAB 2: loop + conditional — all filters must pass
	for _, f := range cf.Filters {
		if !f.Match(event) {
			return false
		}
	}
	return true
}

// ============================================================
// ApplyFilter — utility function to filter a slice of events
// LAB 3: Slice operations (append, range)
// LAB 6: Accepts any Filter interface
// ============================================================
func ApplyFilter(events []models.LogEvent, f Filter) []models.LogEvent {
	// LAB 3: create result slice
	result := make([]models.LogEvent, 0)
	for _, event := range events {
		if f.Match(event) {
			result = append(result, event) // LAB 3: append
		}
	}
	return result
}
