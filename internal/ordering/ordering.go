package ordering

// LAB 3 (Array & Slice) + LAB 7 (Pointers):
// Time Ordering module — buffers incoming LogEvents and sorts by timestamp.
// Implements a 1-second flush window so logs from multiple sources appear chronologically.

import (
	"sort"
	"sync"
	"time"

	"logmerge/internal/models"
)

// ============================================================
// TimeOrderBuffer stores incoming events and flushes them sorted.
// LAB 3: Slice operations (append, sort, clear)
// LAB 7: Pointer receiver methods for efficient mutation
// ============================================================
type TimeOrderBuffer struct {
	events []models.LogEvent // LAB 3: dynamic slice as buffer
	mu     sync.Mutex        // LAB 9: mutex for thread-safe access
	window time.Duration     // flush window duration
}

// NewTimeOrderBuffer creates a buffer with the specified flush window.
// LAB 7: Returns pointer to avoid copying the struct.
func NewTimeOrderBuffer(window time.Duration) *TimeOrderBuffer {
	return &TimeOrderBuffer{
		events: make([]models.LogEvent, 0, 64), // LAB 3: make with capacity
		window: window,
	}
}

// AddEvent appends a log event to the buffer.
// LAB 3: append to slice
// LAB 7: pointer receiver — modifies the buffer in-place
func (tb *TimeOrderBuffer) AddEvent(event models.LogEvent) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	// LAB 3: append to dynamic slice
	tb.events = append(tb.events, event)
}

// Flush sorts all buffered events by timestamp and returns them.
// The buffer is then cleared for the next window.
// LAB 3: Slice sorting + slicing to zero length
// LAB 7: Pointer receiver for efficient in-place sort
func (tb *TimeOrderBuffer) Flush() []models.LogEvent {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if len(tb.events) == 0 {
		return nil
	}

	// LAB 3: sort.Slice — sort the buffered events by timestamp
	sort.Slice(tb.events, func(i, j int) bool {
		return tb.events[i].Timestamp.Before(tb.events[j].Timestamp)
	})

	// LAB 3: Copy sorted events into a new slice to return
	sorted := make([]models.LogEvent, len(tb.events))
	copy(sorted, tb.events) // LAB 3: copy builtin

	// LAB 3: Reset slice to zero length but keep capacity
	tb.events = tb.events[:0]

	return sorted
}

// Count returns the number of events currently in the buffer.
// LAB 3: len() on slice
func (tb *TimeOrderBuffer) Count() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return len(tb.events)
}

// Window returns the flush window duration.
func (tb *TimeOrderBuffer) Window() time.Duration {
	return tb.window
}

// SortEvents is a standalone utility that sorts a slice of LogEvents by timestamp.
// LAB 3: Slice sorting demonstration
// LAB 7: Modifies the slice in-place (slices are reference types)
func SortEvents(events []models.LogEvent) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})
}
