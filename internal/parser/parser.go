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
  if len(parts) < 3 {
    return models.NewLogEvent(
      time.Now(),
      models.UNKNOWN,
      line,
      source,
      line,
    )
  }

  
  timestamp, err := time.Parse(time.RFC3339, parts[0])
  if err != nil {
    //current time
    timestamp = time.Now()
  }

  level := models.ParseLogLevel(parts[1])

  //slicing
  messageParts := parts[2:]
  message := strings.Join(messageParts, " ")

  return models.NewLogEvent(
    timestamp,
    level,
    message,
    source,
    line,
  )
}
