package parser

import (
  "strings"
  "time"
  "logmerge/internal/models"
)

// ParseLogLine converts a raw log line string into a LogEvent struct
// LAB 4: Function that returns a struct
func ParseLogLine(line string, source string) models.LogEvent {
  // LAB 3: split the log line into parts using spaces (creates slice)
  parts := strings.Fields(line)

  // LAB 2: If statement - check if too short
  if len(parts) < 4 { // LAB 3: len() on slice
    // LAB 4: Return struct using constructor
    return models.NewLogEvent(
      time.Now(),
      models.UNKNOWN,
      line,
      source,
      line,
    )
  }

  // LAB 3: Slice indexing - Parse date and time (YYYY-MM-DD HH:MM:SS)
  timestampStr := parts[0] + " " + parts[1]
  timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)

  // LAB 2: If-else for error handling
  if err != nil {
    // try RFC3339 as fallback
    timestamp, err = time.Parse(time.RFC3339, parts[0])
    if err != nil {
      timestamp = time.Now()
    }
  }

  // LAB 1: Custom type from string
  level := models.ParseLogLevel(parts[2])

  // LAB 3: Slicing - message starts from index 4 onwards
  // For our test logs: [date] [time] [level] [source] [message...]
  messageParts := parts[4:] // LAB 3: Slice from index to end [start:]
  message := strings.Join(messageParts, " ")

  // LAB 4: Return struct
  return models.NewLogEvent(
    timestamp,
    level,
    message,
    source,
    line,
  )
}