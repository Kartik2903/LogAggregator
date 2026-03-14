package output

// LAB 5 (Functions & Error Handling) + LAB 8 (JSON):
// Output module — color-coded terminal printing and file export.
// Demonstrates: fmt formatting, ANSI escape codes, file I/O, error handling.

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"logmerge/internal/models"
)

// ANSI color codes for terminal output
// LAB 1: Constants and string types
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m" // ERROR
	ColorYellow = "\033[33m" // WARN
	ColorCyan   = "\033[36m" // INFO
	ColorGray   = "\033[90m" // UNKNOWN / dim text
	ColorBold   = "\033[1m"
)

// ============================================================
// ColorForLevel returns the ANSI color code for a given log level.
// LAB 2: Switch statement for control flow
// ============================================================
func ColorForLevel(level models.LogLevel) string {
	switch level {
	case models.ERROR:
		return ColorRed
	case models.WARN:
		return ColorYellow
	case models.INFO:
		return ColorCyan
	default:
		return ColorGray
	}
}

// ============================================================
// PrintEvent prints a single LogEvent to the terminal with color coding.
// LAB 5: Function with formatted output
// ============================================================
func PrintEvent(event models.LogEvent) {
	color := ColorForLevel(event.Level)

	// LAB 5: fmt.Printf for formatted terminal output
	fmt.Printf("%s%s[%-5s]%s %s[%s]%s %s\n",
		ColorBold, color, event.Level, ColorReset, // colored level badge
		ColorGray, event.Timestamp.Format("2006-01-02 15:04:05"), ColorReset, // dimmed timestamp
		event.Message,
	)
}

// ============================================================
// PrintEventWithSource includes the source file name in the output.
// ============================================================
func PrintEventWithSource(event models.LogEvent) {
	color := ColorForLevel(event.Level)

	fmt.Printf("%s%s[%-5s]%s %s[%s]%s %s(%s)%s %s\n",
		ColorBold, color, event.Level, ColorReset,
		ColorGray, event.Timestamp.Format("2006-01-02 15:04:05"), ColorReset,
		ColorGray, event.Source, ColorReset,
		event.Message,
	)
}

// ============================================================
// PrintSeparator prints a visual separator line.
// ============================================================
func PrintSeparator() {
	fmt.Println(ColorGray + strings.Repeat("─", 70) + ColorReset)
}

// ============================================================
// PrintHeader prints a styled header for the log output.
// ============================================================
func PrintHeader() {
	fmt.Println()
	PrintSeparator()
	fmt.Printf("%s%s  LOG AGGREGATOR — LIVE STREAM  %s\n", ColorBold, ColorCyan, ColorReset)
	PrintSeparator()
	fmt.Println()
}

// ============================================================
// WriteToFile writes a slice of LogEvents to a plain-text log file.
// LAB 5: File I/O with error handling
// ============================================================
func WriteToFile(events []models.LogEvent, path string) error {
	// LAB 5: os.Create with error handling
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err) // LAB 5: error wrapping
	}
	defer file.Close()

	// LAB 2: Range loop over events
	for _, event := range events {
		line := fmt.Sprintf("[%s] [%-5s] [%s] %s\n",
			event.Timestamp.Format("2006-01-02 15:04:05"),
			event.Level,
			event.Source,
			event.Message,
		)
		// LAB 5: Write with error check
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("failed to write event: %w", err)
		}
	}

	return nil
}

// ============================================================
// WriteToJSON exports events as a JSON file.
// LAB 8: JSON marshalling
// ============================================================
func WriteToJSON(events []models.LogEvent, path string) error {
	// LAB 8: json.MarshalIndent for pretty JSON output
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	// LAB 5: os.WriteFile for atomic write
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}
