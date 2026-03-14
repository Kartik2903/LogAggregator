package parser

// LAB 8: Unit tests for the Parser module
// Demonstrates: table-driven tests, testing error cases, test organization

import (
	"testing"

	"logmerge/internal/models"
)

// LAB 8: Table-driven test for ParseLogLine
// Tests valid parsing, invalid formats, and bad timestamps.
func TestParseLogLine(t *testing.T) {
	p := SimpleParser{}

	// LAB 8: Table-driven test — each entry is a test case struct
	tests := []struct {
		name      string          // test case name
		input     string          // raw log line
		source    string          // source file name
		wantLevel models.LogLevel // expected level
		wantMsg   string          // expected message substring
		wantErr   error           // expected error (nil if valid)
	}{
		{
			name:      "Valid INFO line",
			input:     "2025-06-15 10:00:01 INFO - Application started successfully",
			source:    "app.log",
			wantLevel: models.INFO,
			wantMsg:   "Application started successfully",
			wantErr:   nil,
		},
		{
			name:      "Valid ERROR line",
			input:     "2025-06-15 10:00:07 ERROR - Database connection timeout after 30s",
			source:    "app.log",
			wantLevel: models.ERROR,
			wantMsg:   "Database connection timeout after 30s",
			wantErr:   nil,
		},
		{
			name:      "Valid WARN line",
			input:     "2025-06-15 10:00:03 WARN - High memory usage detected: 82%",
			source:    "server.log",
			wantLevel: models.WARN,
			wantMsg:   "High memory usage detected: 82%",
			wantErr:   nil,
		},
		{
			name:    "Invalid format — too few fields",
			input:   "hello world",
			source:  "bad.log",
			wantErr: models.ErrInvalidFormat,
		},
		{
			name:    "Invalid format — empty string",
			input:   "",
			source:  "bad.log",
			wantErr: models.ErrInvalidFormat,
		},
		{
			name:    "Invalid timestamp",
			input:   "not-a-date not-a-time INFO - some message here",
			source:  "bad.log",
			wantErr: models.ErrInvalidTimestamp,
		},
	}

	// LAB 8: Run each test case as a subtest using t.Run
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := p.ParseLogLine(tt.input, tt.source)

			// LAB 8: Check error expectations
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if err != tt.wantErr {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			// LAB 8: No error expected — validate the parsed event
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if event.Level != tt.wantLevel {
				t.Errorf("level: got %v, want %v", event.Level, tt.wantLevel)
			}

			if event.Message != tt.wantMsg {
				t.Errorf("message: got %q, want %q", event.Message, tt.wantMsg)
			}

			if event.Source != tt.source {
				t.Errorf("source: got %q, want %q", event.Source, tt.source)
			}
		})
	}
}

// LAB 8: Benchmark test for parser performance
func BenchmarkParseLogLine(b *testing.B) {
	p := SimpleParser{}
	line := "2025-06-15 10:00:01 INFO - Application started successfully"

	// LAB 8: b.N is the benchmark iteration count
	for i := 0; i < b.N; i++ {
		_, _ = p.ParseLogLine(line, "bench.log")
	}
}
