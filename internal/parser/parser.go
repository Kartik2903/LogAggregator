package parser

import (
  "strings"
  "time"

  "logmerge/internal/models"
)

//converts a raw log line string into a LogEvent struct
func ParseLogLine(line string, source string) models.LogEvent {

  // split the log line into parts using spaces
  parts := strings.Fields(line)

  // itoo short
  if len(parts) < 4 {
    return models.NewLogEvent(
      time.Now(),
      models.UNKNOWN,
      line,
      source,
      line,
    )
  }

  // Parse date and time (YYYY-MM-DD HH:MM:SS)
  timestampStr := parts[0] + " " + parts[1]
  timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
  if err != nil {
    // try RFC3339 as fallback
    timestamp, err = time.Parse(time.RFC3339, parts[0])
    if err != nil {
      timestamp = time.Now()
    }
  }

  level := models.ParseLogLevel(parts[2])

  //slicing - message starts from index 3 or 4 depending on format
  // For our test logs: [date] [time] [level] [source] [message...]
  // Wait, our test logs are: 2026-01-23 10:00:00 INFO frontend User logged in
  // parts[0]=2026-01-23, parts[1]=10:00:00, parts[2]=INFO, parts[3]=frontend
  messageParts := parts[4:]
  message := strings.Join(messageParts, " ")

  return models.NewLogEvent(
    timestamp,
    level,
    message,
    source,
    line,
  )
}
