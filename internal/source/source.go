package source

// LAB 5 (Functions & Error Handling) + LAB 6 (Interface) +
// LAB 9 (Concurrency: WaitGroup) + LAB 10 (Goroutines & Channels):
// Source module — reads multiple log files concurrently and emits raw log lines.

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// ============================================================
// LAB 6: Source Interface — any implementation that can read logs.
// ============================================================
type Source interface {
	Read(path string) (<-chan string, <-chan error) // Single file read
}

// ============================================================
// FileSource reads log lines from files on disk.
// LAB 6: Implements the Source interface.
// ============================================================
type FileSource struct{}

// Read opens a single file and emits each line through a channel.
// LAB 5:  Error handling for file open and scanner errors.
// LAB 10: Goroutine sends lines into channel, closes on completion.
func (fs FileSource) Read(path string) (<-chan string, <-chan error) {
	lines := make(chan string)
	errs := make(chan error, 1)

	// LAB 10: Goroutine for asynchronous file reading
	go func() {
		defer close(lines)
		defer close(errs)

		// LAB 5: Open file with error handling
		file, err := os.Open(path)
		if err != nil {
			errs <- err
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		// LAB 2: Loop control flow — read file line by line
		for scanner.Scan() {
			lines <- scanner.Text() // LAB 10: Send line into channel
		}

		// LAB 5: Scanner error handling
		if err := scanner.Err(); err != nil {
			errs <- err
		}
	}()

	return lines, errs
}

// ============================================================
// ReadMultipleConcurrent reads multiple files concurrently.
// Each file gets its own goroutine; all lines are merged into
// a single output channel using the Fan-In pattern.
//
// LAB 9:  sync.WaitGroup to wait for all file-reading goroutines.
// LAB 10: One goroutine per file, channels for communication.
// ============================================================
func ReadMultipleConcurrent(paths []string) (<-chan string, <-chan error) {
	// LAB 10: Create unified output channels
	merged := make(chan string)
	errs := make(chan error, len(paths)) // LAB 10: Buffered error channel

	// LAB 9: WaitGroup to track all reader goroutines
	var wg sync.WaitGroup

	// LAB 10 + LAB 2: Launch one goroutine per file
	for _, path := range paths {
		wg.Add(1) // LAB 9: Increment WaitGroup counter

		// LAB 10: Goroutine for concurrent file reading
		go func(filePath string) {
			defer wg.Done() // LAB 9: Decrement when goroutine completes

			// LAB 5: Open file with error handling
			file, err := os.Open(filePath)
			if err != nil {
				errs <- fmt.Errorf("[%s] %w", filePath, err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			// LAB 2: Line-by-line reading
			for scanner.Scan() {
				// LAB 10: Prefix each line with source ID for downstream parsing
				merged <- filePath + "|" + scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				errs <- fmt.Errorf("[%s] scanner: %w", filePath, err)
			}
		}(path) // LAB 10: Pass path to avoid closure capture issue
	}

	// LAB 9: Goroutine that waits for all readers to finish, then closes channels
	go func() {
		wg.Wait()
		close(merged)
		close(errs)
	}()

	return merged, errs
}
