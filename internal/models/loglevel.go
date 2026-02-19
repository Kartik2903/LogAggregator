package models

import "strings"

// LAB 1: Custom type definition based on int
type LogLevel int

// LAB 1: Constants with iota (enumeration)
const (
  UNKNOWN LogLevel = iota // 0
  INFO                    // 1
  WARN                    // 2
  ERROR                   // 3
)

// LAB 1 & 4: Method on custom type
// String implements the Stringer interface
func (l LogLevel) String() string {
  // LAB 2: Switch statement on custom type
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

// ParseLogLevel converts a string to LogLevel
// Demonstrates: string manipulation, switch, custom type
func ParseLogLevel(s string) LogLevel {
  // LAB 1: Type conversion - string to uppercase
  upperStr := strings.ToUpper(s)

  // LAB 2: Switch with string cases
  switch upperStr {
  case "INFO":
    return INFO
  case "WARN", "WARNING": // Multiple values in one case
    return WARN
  case "ERROR", "ERR":
    return ERROR
  default:
    return UNKNOWN
  }
}

// Enabled checks if this log level meets the minimum threshold
// LAB 1: Method with custom type receiver
// LAB 2: Boolean comparison
func (l LogLevel) Enabled(min LogLevel) bool {
  return l >= min // Numeric comparison on custom int type
}

// IsError checks if this is an error level
func (l LogLevel) IsError() bool {
  return l == ERROR
}

// IsWarningOrHigher checks if level is WARN or ERROR
func (l LogLevel) IsWarningOrHigher() bool {
  return l >= WARN
}

// ToInt converts LogLevel to int
// LAB 1: Type conversion from custom type to int
func (l LogLevel) ToInt() int {
  return int(l)
}

// FromInt creates a LogLevel from an integer
// LAB 1: Type conversion from int to custom type
func FromInt(i int) LogLevel {
  // LAB 2: If-else chain for validation
  if i < 0 || i > 3 {
    return UNKNOWN
  }
  return LogLevel(i)
}

// Compare returns -1 if l < other, 0 if equal, 1 if l > other
func (l LogLevel) Compare(other LogLevel) int {
  if l < other {
    return -1
  } else if l > other {
    return 1
  }
  return 0
}