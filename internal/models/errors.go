package models

import "errors"

var (
  ErrInvalidFormat    = errors.New("invalid log format")
  ErrInvalidTimestamp = errors.New("invalid timestamp")
)