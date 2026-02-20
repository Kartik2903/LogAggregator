package parser

import (
  "strings"
  "time"

  "logmerge/internal/models"
)

// Interface (Lab 6 concept — OK to keep)
type Parser interface {
  ParseLogLine(line string, source string) (models.LogEvent, error)
}

// Implementation struct
type SimpleParser struct{}

// LAB 5: Function returning (value, error)
func (p SimpleParser) ParseLogLine(line string, source string) (models.LogEvent, error) {

  // LAB 3: Split string into slice
  parts := strings.Fields(line)

  // LAB 2 + LAB 5: validation + error return
  if len(parts) < 4 {
    return models.LogEvent{}, models.ErrInvalidFormat
  }

  // Combine date + time
  timestampStr := parts[0] + " " + parts[1]

  // LAB 5: function with error handling
  timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
  if err != nil {
    // fallback format
    timestamp, err = time.Parse(time.RFC3339, parts[0])
    if err != nil {
      return models.LogEvent{}, models.ErrInvalidTimestamp
    }
  }

  // LAB 1: custom type conversion
  level := models.ParseLogLevel(parts[2])

  // LAB 3: slice from index
  messageParts := parts[4:]
  message := strings.Join(messageParts, " ")

  // LAB 5: constructor function usage
  event := models.NewLogEvent(timestamp, level, message, source, line)

  return event, nil
}