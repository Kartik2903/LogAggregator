package main

import (
  "flag"
  "fmt"

  "logmerge/internal/models"
  "logmerge/internal/parser"
  "logmerge/internal/source"
)

func main() {

  // LAB 1: variables + slice
  var files []string
  var showStats bool

  levelStr := flag.String("level", "INFO", "Minimum log level")

  flag.BoolVar(&showStats, "stats", false, "Show statistics")

  flag.Func("file", "Log file path", func(v string) error {
    files = append(files, v) // LAB 3: append slice
    return nil
  })

  flag.Parse()

  // LAB 2: control flow
  if len(files) == 0 {
    fmt.Println("No files provided")
    return
  }

  // LAB 1: custom type
  minLevel := models.ParseLogLevel(*levelStr)

  // LAB 2: switch statement
  switch minLevel {
  case models.INFO, models.WARN, models.ERROR:
    fmt.Println("Filter level set to:", minLevel)
  default:
    minLevel = models.INFO
  }

  // LAB 5: using functions through interface
  var p parser.Parser = parser.SimpleParser{}
  var src source.Source = source.FileSource{}

  storage := models.NewLogStorage()
  buffer := models.NewLogBuffer()

  // LAB 2 + LAB 3: loop over slice
  for fileIdx, file := range files {

    fmt.Printf("\n[%d] Reading file: %s\n", fileIdx+1, file)

    lines, errs := src.Read(file)

    lineNum := 0

    // LAB 2: loop over channel
    for line := range lines {

      lineNum++

      // LAB 5: FUNCTION CALL WITH ERROR HANDLING
      event, err := p.ParseLogLine(line, file)
      if err != nil {
        fmt.Println("Skipped:", err) // error handled safely
        continue
      }

      storage.AddEvent(event)
      buffer.AddEvent(event)

      if event.Level.Enabled(minLevel) {
        fmt.Printf("  Line %3d: [%s] [%-5v] [%s] %s\n",
          lineNum,
          event.Timestamp.Format("15:04:05"),
          event.Level,
          event.Source,
          event.Message,
        )
      }
    }

    // LAB 5: handle source errors
    if err := <-errs; err != nil {
      fmt.Println("Source error:", err)
    }
  }

  if showStats {
    printStatistics(storage, buffer)
  }
}