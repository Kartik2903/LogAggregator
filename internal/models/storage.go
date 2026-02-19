package models

import (
  "sync"
  "time"
)

// LAB 4: Map and structs demonstration
// This module provides map-based storage for log events

// LogStorage manages log events using maps
// Demonstrates: maps, structs, struct composition
type LogStorage struct {
  // LAB 4: Map - source name -> list of events from that source
  EventsBySource map[string][]LogEvent

  // LAB 4: Map - log level -> count of events at that level
  LevelCounts map[LogLevel]int

  // LAB 4: Map - timestamp (as string) -> event
  EventsByTime map[string]LogEvent

  // LAB 4: Map - source -> metadata
  SourceMetadata map[string]SourceInfo

  mu sync.RWMutex // For thread-safe operations (will be used in concurrency lab)
}

// LAB 4: Struct definition for source metadata
type SourceInfo struct {
  Name          string    // Source identifier
  FilePath      string    // Path to log file
  LastSeen      time.Time // Last time we saw a log from this source
  TotalEvents   int       // Total events from this source
  ErrorCount    int       // Number of errors from this source
  LastEventTime time.Time // Timestamp of most recent event
}

// LAB 4: Struct with embedded struct
type EnrichedLogEvent struct {
  LogEvent           // LAB 4: Embedded struct (composition)
  SourceMeta SourceInfo // Additional metadata
  Tags       []string   // Custom tags
}

// NewLogStorage creates a new log storage
// Demonstrates: map initialization with make()
func NewLogStorage() *LogStorage {
  return &LogStorage{
    // LAB 4: Initialize maps using make()
    EventsBySource: make(map[string][]LogEvent),
    LevelCounts:    make(map[LogLevel]int),
    EventsByTime:   make(map[string]LogEvent),
    SourceMetadata: make(map[string]SourceInfo),
  }
}

// AddEvent stores a log event
// Demonstrates: map insertion, map access, map update
func (ls *LogStorage) AddEvent(event LogEvent) {
  ls.mu.Lock()
  defer ls.mu.Unlock()

  // LAB 4: Map insertion - add event to source-based map
  ls.EventsBySource[event.Source] = append(ls.EventsBySource[event.Source], event)

  // LAB 4: Map update - increment level count
  ls.LevelCounts[event.Level]++

  // LAB 4: Map insertion with composite key
  timeKey := event.Timestamp.Format(time.RFC3339Nano)
  ls.EventsByTime[timeKey] = event

  // LAB 4: Update source metadata
  ls.updateSourceMetadata(event)
}

// updateSourceMetadata updates metadata for a source
// Demonstrates: map access, map update with struct
func (ls *LogStorage) updateSourceMetadata(event LogEvent) {
  // LAB 4: Map access - check if key exists
  meta, exists := ls.SourceMetadata[event.Source]

  if !exists {
    // LAB 4: Create new struct and add to map
    meta = SourceInfo{
      Name:          event.Source,
      FilePath:      event.Source, // In real scenario, this would be actual path
      LastSeen:      event.Timestamp,
      TotalEvents:   0,
      ErrorCount:    0,
      LastEventTime: event.Timestamp,
    }
  }

  // LAB 4: Struct field access and modification
  meta.TotalEvents++
  meta.LastSeen = time.Now()
  meta.LastEventTime = event.Timestamp

  if event.Level == ERROR {
    meta.ErrorCount++
  }

  // LAB 4: Map update with modified struct
  ls.SourceMetadata[event.Source] = meta
}

// GetEventsBySource retrieves all events from a specific source
// Demonstrates: map lookup, returning slice from map
func (ls *LogStorage) GetEventsBySource(source string) []LogEvent {
  ls.mu.RLock()
  defer ls.mu.RUnlock()

  // LAB 4: Map lookup
  events, exists := ls.EventsBySource[source]
  if !exists {
    return []LogEvent{} // Return empty slice if not found
  }

  return events
}

// GetLevelCount returns the count of events at a specific level
// Demonstrates: map access
func (ls *LogStorage) GetLevelCount(level LogLevel) int {
  ls.mu.RLock()
  defer ls.mu.RUnlock()

  // LAB 4: Map access (returns zero value if key doesn't exist)
  return ls.LevelCounts[level]
}

// GetAllSources returns a list of all known sources
// Demonstrates: iterating over map keys
func (ls *LogStorage) GetAllSources() []string {
  ls.mu.RLock()
  defer ls.mu.RUnlock()

  // LAB 4: Iterate over map to get keys
  sources := make([]string, 0, len(ls.EventsBySource))
  for source := range ls.EventsBySource {
    sources = append(sources, source)
  }

  return sources
}

// GetSourceMetadata retrieves metadata for a source
// Demonstrates: map lookup with existence check
func (ls *LogStorage) GetSourceMetadata(source string) (SourceInfo, bool) {
  ls.mu.RLock()
  defer ls.mu.RUnlock()

  // LAB 4: Map lookup with "comma ok" idiom
  meta, exists := ls.SourceMetadata[source]
  return meta, exists
}

// GetAllMetadata returns all source metadata
// Demonstrates: returning map
func (ls *LogStorage) GetAllMetadata() map[string]SourceInfo {
  ls.mu.RLock()
  defer ls.mu.RUnlock()

  // LAB 4: Copy map to prevent external modification
  metaCopy := make(map[string]SourceInfo, len(ls.SourceMetadata))
  for k, v := range ls.SourceMetadata {
    metaCopy[k] = v
  }

  return metaCopy
}

// RemoveSource removes all data for a source
// Demonstrates: map deletion
func (ls *LogStorage) RemoveSource(source string) {
  ls.mu.Lock()
  defer ls.mu.Unlock()

  // LAB 4: Delete from map using delete() builtin
  delete(ls.EventsBySource, source)
  delete(ls.SourceMetadata, source)
}

// GetLevelStatistics returns a summary of events by level
// Demonstrates: map iteration, creating new map from existing map
func (ls *LogStorage) GetLevelStatistics() map[string]int {
  ls.mu.RLock()
  defer ls.mu.RUnlock()

  // LAB 4: Create new map and populate from existing map
  stats := make(map[string]int)

  for level, count := range ls.LevelCounts {
    // Convert LogLevel to string for the map
    levelStr := ""
    switch level {
    case INFO:
      levelStr = "INFO"
    case WARN:
      levelStr = "WARN"
    case ERROR:
      levelStr = "ERROR"
    default:
      levelStr = "UNKNOWN"
    }
    stats[levelStr] = count
  }

  return stats
}

// Clear removes all stored data
// Demonstrates: map re-initialization
func (ls *LogStorage) Clear() {
  ls.mu.Lock()
  defer ls.mu.Unlock()

  // LAB 4: Re-initialize maps to clear them
  ls.EventsBySource = make(map[string][]LogEvent)
  ls.LevelCounts = make(map[LogLevel]int)
  ls.EventsByTime = make(map[string]LogEvent)
  ls.SourceMetadata = make(map[string]SourceInfo)
}

// TotalEvents returns total number of events across all sources
func (ls *LogStorage) TotalEvents() int {
  ls.mu.RLock()
  defer ls.mu.RUnlock()

  total := 0
  // LAB 4: Iterate over map values
  for _, events := range ls.EventsBySource {
    total += len(events)
  }
  return total
}