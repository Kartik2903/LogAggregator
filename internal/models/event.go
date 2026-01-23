package models

import "time"

type LogEvent struct {
	Timestamp time.Time `json:"timestamp"` //to sort later
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"` //actual log text
	Source    string    `json:"source"`
	RawLine   string    `json:"raw_line,omitempty"` //preserves data if parsing fails
}
// NewLogEvent is a constructor that creates a fully initialized LogEvent.
func NewLogEvent(
	ts time.Time,
	level LogLevel,
	msg string,
	source string,
	raw string,
) LogEvent {
	return LogEvent{
		Timestamp: ts,
		Level:     level,
		Message:   msg,
		Source:    source,
		RawLine:   raw,
	}
}
