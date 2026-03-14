package fanin

// LAB 9 (Concurrency: WaitGroup, Mutex) + LAB 10 (Goroutines & Channels, Select):
// Fan-In module — merges multiple input channels into a single unified stream.
// This is the concurrency core of the Log Aggregator pipeline.

import (
	"sync"
)

// ============================================================
// MergeChannels takes multiple read-only string channels and
// merges them into a single output channel using the Fan-In pattern.
//
// LAB 9:  sync.WaitGroup to track when all source goroutines finish
// LAB 10: One goroutine per input channel; select is used implicitly
//
//	via the range-over-channel pattern.
//
// ============================================================
func MergeChannels(channels []<-chan string) <-chan string {
	// LAB 10: Create the unified output channel
	merged := make(chan string)

	// LAB 9: WaitGroup tracks completion of all reader goroutines
	var wg sync.WaitGroup

	// LAB 10: Launch one goroutine per input channel
	for _, ch := range channels {
		wg.Add(1) // LAB 9: Increment wait counter

		// LAB 10: Goroutine reads from its assigned channel
		// and forwards every value to the merged channel.
		go func(c <-chan string) {
			defer wg.Done() // LAB 9: Decrement counter when done

			// LAB 10: Range over channel — blocks until channel closes
			for line := range c {
				merged <- line // Forward to merged output
			}
		}(ch) // LAB 10: Pass channel as argument to avoid closure capture issues
	}

	// LAB 9 + LAB 10: Separate goroutine waits for all readers,
	// then closes the merged channel to signal downstream consumers.
	go func() {
		wg.Wait()     // LAB 9: Block until all goroutines call Done()
		close(merged) // LAB 10: Close channel — downstream range loops will exit
	}()

	return merged
}

// ============================================================
// MergeWithSelect merges channels using an explicit select statement.
// This is an alternative Fan-In implementation that demonstrates
// the select keyword directly.
//
// LAB 10: select statement for multiplexing channels
// ============================================================
func MergeWithSelect(ch1 <-chan string, ch2 <-chan string) <-chan string {
	merged := make(chan string)

	go func() {
		defer close(merged)

		// LAB 10: Track which channels are still open
		ch1Open, ch2Open := true, true

		// LAB 2: Loop until both channels are closed
		for ch1Open || ch2Open {
			// LAB 10: select statement — waits on multiple channel operations
			select {
			case line, ok := <-ch1:
				if !ok {
					// LAB 10: Channel closed — stop selecting from it
					ch1Open = false
					continue
				}
				merged <- line

			case line, ok := <-ch2:
				if !ok {
					ch2Open = false
					continue
				}
				merged <- line
			}
		}
	}()

	return merged
}
