package models

import "strings" //importing the strings package

type LogLevel int //severity of log message

const(           //log levels
  UNKNOWN LogLevel = iota
  INFO
  WARN
  ERROR
)

func ParseLogLevel(s string) LogLevel {

  // convert input string to uppercase
  switch strings.ToUpper(s) {

  case "INFO":
    return INFO

  case "WARN", "WARNING":
    return WARN

  case "ERROR":
    return ERROR

  default:
    return UNKNOWN //if log not registered, return UNKNOWN
  }
}

func (l LogLevel) Enabled(min LogLevel) bool {
  return l >= min
}