package models

// LAB 3: Array and slice demonstration
// This module handles buffering of log events for batch processing

const MaxBufferSize = 100 // Fixed size constant for array demonstration

// LogBuffer manages a collection of log events with both arrays and slices
type LogBuffer struct {
  // Array: Fixed-size recent events (demonstrates array concept)
  RecentEvents [10]LogEvent // LAB 3: Fixed-size array

  // Slice: Dynamic event storage (demonstrates slice concept)
  Events []LogEvent // LAB 3: Dynamic slice

  // Slice for tracking sources
  Sources []string // LAB 3: String slice
}

// NewLogBuffer creates a new log buffer
func NewLogBuffer() *LogBuffer {
  return &LogBuffer{
    Events:  make([]LogEvent, 0, MaxBufferSize), // LAB 3: make() with capacity
    Sources: make([]string, 0, 10),
  }
}

// AddEvent adds a log event to the buffer
// Demonstrates: append, slicing, array assignment
func (lb *LogBuffer) AddEvent(event LogEvent) {
  // LAB 3: Append to slice
  lb.Events = append(lb.Events, event)

  // LAB 3: Array usage - keep last 10 events in fixed array
  // Shift all elements left and add new one at end
  for i := 0; i < len(lb.RecentEvents)-1; i++ {
    lb.RecentEvents[i] = lb.RecentEvents[i+1]
  }
  lb.RecentEvents[len(lb.RecentEvents)-1] = event

  // Track unique sources
  lb.addSource(event.Source)
}

// addSource adds source if not already present
func (lb *LogBuffer) addSource(source string) {
  // LAB 3: Range over slice (demonstrates range with index and value)
  for _, existingSource := range lb.Sources {
    if existingSource == source {
      return
    }
  }
  lb.Sources = append(lb.Sources, source)
}

// GetEventsByLevel returns events matching a specific level
// Demonstrates: slicing, filtering with range
func (lb *LogBuffer) GetEventsByLevel(level LogLevel) []LogEvent {
  // LAB 3: Create new slice for filtered results
  filtered := make([]LogEvent, 0)

  // LAB 3: Range over slice with both index and value
  for idx, event := range lb.Events {
    if event.Level == level {
      filtered = append(filtered, event)
      // Example of using index
      _ = idx // Acknowledging we have access to index
    }
  }

  return filtered
}

// GetRecentN returns the last N events
// Demonstrates: slicing with range syntax
func (lb *LogBuffer) GetRecentN(n int) []LogEvent {
  // LAB 3: Slicing operation
  if n > len(lb.Events) {
    n = len(lb.Events)
  }

  // Return last N events using slice range [start:]
  startIdx := len(lb.Events) - n
  return lb.Events[startIdx:] // LAB 3: Slice from index to end
}

// GetEventsInRange returns events between two indices
// Demonstrates: slice range [start:end]
func (lb *LogBuffer) GetEventsInRange(start, end int) []LogEvent {
  // LAB 3: Bounds checking
  if start < 0 {
    start = 0
  }
  if end > len(lb.Events) {
    end = len(lb.Events)
  }
  if start >= end {
    return []LogEvent{}
  }

  // LAB 3: Slice with start and end [start:end]
  return lb.Events[start:end]
}

// Clear removes all events but keeps capacity
func (lb *LogBuffer) Clear() {
  // LAB 3: Slice to zero length but maintain capacity
  lb.Events = lb.Events[:0]
  lb.Sources = lb.Sources[:0]

  // Reset array to zero values
  lb.RecentEvents = [10]LogEvent{} // LAB 3: Array initialization
}

// Count returns number of events in buffer
func (lb *LogBuffer) Count() int {
  // LAB 3: len() function on slice
  return len(lb.Events)
}

// Capacity returns the current capacity of the buffer
func (lb *LogBuffer) Capacity() int {
  // LAB 3: cap() function on slice
  return cap(lb.Events)
}

// GetSourcesList returns all unique sources
// Demonstrates: returning a slice
func (lb *LogBuffer) GetSourcesList() []string {
  // LAB 3: Return slice copy to prevent external modification
  sourcesCopy := make([]string, len(lb.Sources))
  copy(sourcesCopy, lb.Sources) // LAB 3: copy() builtin
  return sourcesCopy
}