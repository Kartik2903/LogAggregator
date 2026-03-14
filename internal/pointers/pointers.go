package pointers

// LAB 7: Pointers, method sets — Efficient passing of log events
// This module demonstrates call-by-value vs. call-by-pointer for LogEvent processing.

import (
	"logmerge/internal/models"
	"time"
)

// LAB 7: EnrichEvent modifies a LogEvent IN-PLACE through a pointer receiver.
// This avoids copying the entire struct and allows the caller to see the changes.
// Demonstrates: pointer parameter (*models.LogEvent) vs value parameter.
func EnrichEvent(event *models.LogEvent, tag string) {
	// LAB 7: Dereferencing pointer to modify the original struct
	event.Message = "[" + tag + "] " + event.Message
}

// LAB 7: CompareByValue receives two LogEvents BY VALUE (copies are made).
// Changes inside this function do NOT affect the originals.
// Demonstrates: call-by-value semantics — each parameter is a full copy.
func CompareByValue(a models.LogEvent, b models.LogEvent) bool {
	// LAB 7: These are copies; modifying them here has no side effect.
	return a.Timestamp.Before(b.Timestamp)
}

// LAB 7: SwapByPointer swaps two LogEvent values using pointers.
// Without pointers, a swap function would only swap local copies.
// Demonstrates: pointer-based swap — the classic pointer exercise.
func SwapByPointer(a *models.LogEvent, b *models.LogEvent) {
	// LAB 7: Dereferencing pointers to swap the underlying values
	*a, *b = *b, *a
}

// LAB 7: BatchEnrich processes a slice of *LogEvent pointers efficiently.
// Because we use a slice of pointers, each element is 8 bytes (pointer size)
// rather than the full LogEvent struct size.
// Demonstrates: []*T pattern for efficient batch processing.
func BatchEnrich(events []*models.LogEvent, tag string) {
	// LAB 7: Range over pointer slice — each 'event' is already a pointer
	for _, event := range events {
		EnrichEvent(event, tag) // modifies original through pointer
	}
}

// LAB 7: FindLatest returns a POINTER to the most recent LogEvent in the slice.
// Returning a pointer lets the caller read or modify the found event directly
// without an extra copy.
// Demonstrates: returning *T from a function.
func FindLatest(events []models.LogEvent) *models.LogEvent {
	if len(events) == 0 {
		// LAB 7: Returning nil pointer when no result is found
		return nil
	}

	// LAB 7: Use pointer to track the latest event by index
	latestIdx := 0
	for i := 1; i < len(events); i++ {
		if events[i].Timestamp.After(events[latestIdx].Timestamp) {
			latestIdx = i
		}
	}

	// LAB 7: Return address of the slice element — caller gets a pointer
	return &events[latestIdx]
}

// LAB 7: UpdateTimestamp demonstrates modifying a struct field through a pointer.
// Demonstrates: pointer receiver pattern used in method sets.
func UpdateTimestamp(event *models.LogEvent, newTime time.Time) {
	// LAB 7: Direct field access through pointer (Go auto-dereferences)
	event.Timestamp = newTime
}

// LAB 7: CloneEvent creates a deep copy of a LogEvent.
// Demonstrates: value semantics — assigning a struct creates a copy.
func CloneEvent(original *models.LogEvent) models.LogEvent {
	// LAB 7: Dereferencing pointer creates an independent copy
	clone := *original
	return clone
}
