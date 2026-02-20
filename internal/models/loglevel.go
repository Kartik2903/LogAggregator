package models

import (
  "strings"
)

type LogLevel int

const (
  UNKNOWN LogLevel = iota
  INFO
  WARN
  ERROR
)

// String implements fmt.Stringer (safe, already in your labs)
func (l LogLevel) String() string {
  switch l {
  case INFO:
    return "INFO"
  case WARN:
    return "WARN"
  case ERROR:
    return "ERROR"
  default:
    return "UNKNOWN"
  }
}

// Parse string -> LogLevel
func ParseLogLevel(s string) LogLevel {
  switch strings.ToUpper(strings.TrimSpace(s)) {
  case "INFO":
    return INFO
  case "WARN", "WARNING":
    return WARN
  case "ERROR", "ERR":
    return ERROR
  default:
    return UNKNOWN
  }
}

// Comparison helpers
func (l LogLevel) Enabled(min LogLevel) bool {
  return l >= min
}

func (l LogLevel) IsError() bool {
  return l == ERROR
}

func (l LogLevel) IsWarningOrHigher() bool {
  return l >= WARN
}

// Conversions
func (l LogLevel) ToInt() int {
  return int(l)
}

func FromInt(i int) LogLevel {
  if i < 0 || i > 3 {
    return UNKNOWN
  }
  return LogLevel(i)
}

// Compare two levels
func (l LogLevel) Compare(other LogLevel) int {
  if l < other {
    return -1
  } else if l > other {
    return 1
  }
  return 0
}