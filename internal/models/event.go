package models

import "time"

// LAB 1 & 4: Struct definition with multiple field types
type LogEvent struct {
	Timestamp time.Time `json:"timestamp"` // LAB 1: time.Time type, LAB 4: struct tags
	Level     LogLevel  `json:"level"`     // LAB 1: Custom type
	Message   string    `json:"message"`   // LAB 1: string type
	Source    string    `json:"source"`    // LAB 1: string type
	RawLine   string    `json:"raw_line,omitempty"` // LAB 4: omitempty tag
}

// NewLogEvent is a constructor that creates a fully initialized LogEvent.
// LAB 4: Constructor function for struct
func NewLogEvent(
	ts time.Time,
	level LogLevel,
	msg string,
	source string,
	raw string,
) LogEvent {
	// LAB 4: Struct literal initialization
	return LogEvent{
		Timestamp: ts,
		Level:     level,
		Message:   msg,
		Source:    source,
		RawLine:   raw,
	}
}