package config

// LAB 8: JSON Marshal/Unmarshal — Configuration handling and JSON export.
// Demonstrates: struct tags, json.Marshal, json.Unmarshal, file I/O.

import (
	"encoding/json"
	"fmt"
	"os"

	"logmerge/internal/models"
)

// ============================================================
// AppConfig holds application configuration loaded from JSON.
// LAB 8: Struct with JSON tags for serialization/deserialization.
// LAB 4: Struct definition with multiple field types.
// ============================================================
type AppConfig struct {
	Files      []string `json:"files"`                 // LAB 8: JSON array
	Level      string   `json:"level"`                 // Minimum log level
	Match      string   `json:"match,omitempty"`       // LAB 8: omitempty tag
	OutputFile string   `json:"output_file,omitempty"` // Where to save filtered logs
	ShowStats  bool     `json:"show_stats"`            // Show statistics at end
	ColorMode  bool     `json:"color_mode"`            // Enable/disable color output
}

// ============================================================
// DefaultConfig returns a sane default configuration.
// LAB 4: Struct literal initialization.
// ============================================================
func DefaultConfig() AppConfig {
	return AppConfig{
		Files:     []string{},
		Level:     "INFO",
		ColorMode: true,
	}
}

// ============================================================
// LoadConfig reads a JSON config file and unmarshals it into AppConfig.
// LAB 8: json.Unmarshal demonstration.
// LAB 5: Error handling with descriptive messages.
// ============================================================
func LoadConfig(path string) (*AppConfig, error) {
	// LAB 5: Read file with error handling
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// LAB 8: JSON Unmarshal — parse JSON bytes into struct
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err) // LAB 5: error wrapping
	}

	return &cfg, nil // LAB 7: return pointer
}

// ============================================================
// SaveConfig marshals AppConfig to JSON and writes to file.
// LAB 8: json.MarshalIndent for pretty-printed JSON output.
// ============================================================
func SaveConfig(path string, cfg *AppConfig) error {
	// LAB 8: JSON MarshalIndent — struct to pretty JSON bytes
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// LAB 5: Write file with error handling
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ============================================================
// ExportEventsJSON exports a slice of LogEvents as a JSON file.
// LAB 8: JSON marshalling of complex structs with nested types.
// ============================================================
func ExportEventsJSON(events []models.LogEvent, path string) error {
	// LAB 8: MarshalIndent for human-readable JSON export
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON export: %w", err)
	}

	fmt.Printf("Exported %d events to %s\n", len(events), path)
	return nil
}
