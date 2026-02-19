package main

import (
  "flag"
  "fmt"
  "logmerge/internal/models"
  "logmerge/internal/parser"
  "logmerge/internal/source"
)

func main() {
  // LAB 1: Variables and types
  var files []string                     // LAB 3: Slice declaration
  var showStats bool                     // LAB 1: Boolean variable

  levelStr := flag.String("level", "INFO", "Minimum log level (INFO/WARN/ERROR)")
  flag.BoolVar(&showStats, "stats", false, "Show statistics at the end")
  flag.Func("file", "Log file path (can be repeated)", func(v string) error {
    files = append(files, v) // LAB 3: Append to slice
    return nil
  })
  flag.Parse()

  // LAB 2: Control flow - if statement
  if len(files) == 0 { // LAB 3: len() on slice
    fmt.Println("No files provided")
    fmt.Println("Usage: go run ./cmd/logmerge --file frontend.log --file api.log --level WARN")
    return
  }

  // LAB 1: Type conversion and custom type
  minLevel := models.ParseLogLevel(*levelStr)

  // LAB 2: Switch statement
  switch minLevel {
  case models.INFO, models.WARN, models.ERROR:
    fmt.Println("Filter level set to:", minLevel)
  default:
    fmt.Println("Invalid level provided, using INFO as default")
    minLevel = models.INFO
  }

  // LAB 4: Create map-based storage
  storage := models.NewLogStorage()

  // LAB 3: Create buffer for recent events
  buffer := models.NewLogBuffer()

  // LAB 2: Loop - for range over slice
  for fileIdx, file := range files { // LAB 3: Range with index and value
    fmt.Printf("\n[%d] Reading file: %s\n", fileIdx+1, file)

    lines, err := source.ReadFileLines(file)

    // LAB 2: If-else for error handling
    if err != nil {
      fmt.Println("Failed to read file:", err)
      continue // LAB 2: Continue statement
    }

    fmt.Printf("Found %d lines in %s\n", len(lines), file) // LAB 3: len()

    // LAB 2: Nested loop - for range over slice
    for lineNum, line := range lines { // LAB 3: Range with both index and value
      // LAB 1: Function call returning struct
      event := parser.ParseLogLine(line, file)

      // LAB 4: Store in map-based storage
      storage.AddEvent(event)

      // LAB 3: Add to buffer
      buffer.AddEvent(event)

      // LAB 2: If statement with method call
      if event.Level.Enabled(minLevel) {
        // LAB 1: String formatting with multiple types
        fmt.Printf("  Line %3d: [%s] [%-5v] [%s] %s\n",
          lineNum+1,
          event.Timestamp.Format("15:04:05"),
          event.Level,
          event.Source,
          event.Message,
        )
      }
    }
  }

  // LAB 2: If statement to show stats
  if showStats {
    printStatistics(storage, buffer)
  }
}

// LAB 4: Function demonstrating maps and structs
func printStatistics(storage *models.LogStorage, buffer *models.LogBuffer) {
  fmt.Println("\n" + "=== LOG STATISTICS ===")

  // LAB 4: Get all sources from map
  sources := storage.GetAllSources()
  fmt.Printf("Total sources: %d\n", len(sources)) // LAB 3: len()

  // LAB 2: Loop over slice
  for idx, source := range sources { // LAB 3: Range over slice
    // LAB 4: Map lookup
    events := storage.GetEventsBySource(source)
    // LAB 4: Map lookup with existence check
    meta, exists := storage.GetSourceMetadata(source)

    fmt.Printf("%d. Source: %s\n", idx+1, source)
    fmt.Printf("   Events: %d\n", len(events)) // LAB 3: len()

    // LAB 2: If-else
    if exists {
      // LAB 4: Struct field access
      fmt.Printf("   Errors: %d\n", meta.ErrorCount)
      fmt.Printf("   Last seen: %s\n", meta.LastSeen.Format("15:04:05"))
    }
  }

  // LAB 4: Get level statistics (map)
  fmt.Println("\n=== LEVEL BREAKDOWN ===")
  levelStats := storage.GetLevelStatistics()

  // LAB 4: Iterate over map
  for level, count := range levelStats {
    fmt.Printf("%s: %d\n", level, count)
  }

  // LAB 3: Array and slice operations
  fmt.Println("\n=== BUFFER STATISTICS ===")
  fmt.Printf("Total events in buffer: %d\n", buffer.Count())      // LAB 3: Slice len
  fmt.Printf("Buffer capacity: %d\n", buffer.Capacity())          // LAB 3: Slice cap
  fmt.Printf("Unique sources: %d\n", len(buffer.GetSourcesList())) // LAB 3: len()

  // LAB 3: Demonstrate slicing operations
  recentThree := buffer.GetRecentN(3) // LAB 3: Slice operation
  if len(recentThree) > 0 {
    fmt.Println("\n=== LAST 3 EVENTS ===")
    // LAB 2 & 3: For range over returned slice
    for i, event := range recentThree {
      fmt.Printf("%d. [%s] %s: %s\n", 
        i+1,
        event.Timestamp.Format("15:04:05"),
        event.Source,
        event.Message)
    }
  }

  // LAB 3: Demonstrate range slicing
  if buffer.Count() > 5 {
    midRange := buffer.GetEventsInRange(2, 5) // LAB 3: [start:end] slicing
    fmt.Printf("\nEvents in range [2:5]: %d events\n", len(midRange))
  }
}