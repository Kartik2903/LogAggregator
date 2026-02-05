package main

import (
  "flag"
  "fmt"

  "logmerge/internal/models"
  "logmerge/internal/parser"
  "logmerge/internal/source"
)

func main() {

  var files []string
  levelStr := flag.String("level", "INFO", "Minimum log level (INFO/WARN/ERROR)")
  flag.Func("file", "Log file path (can be repeated)", func(v string) error {
    files = append(files, v)
    return nil
  })
  flag.Parse()


  if len(files) == 0 {
    fmt.Println("No files provided")
    fmt.Println(`go run ./cmd/logmerge --file frontend.log --file api.log --level WARN`)
    return
  }

  minLevel := models.ParseLogLevel(*levelStr)

  //switch
  switch minLevel {
  case models.INFO, models.WARN, models.ERROR:
    fmt.Println("Filter level set to:", minLevel)
  default:
    fmt.Println("Invalid level provided, using INFO as default")
    minLevel = models.INFO
  }

//LOOP
  for _, file := range files {

    fmt.Println("\nReading file:", file)

    lines, err := source.ReadFileLines(file)
    if err != nil {
      fmt.Println("Failed to read file:", err)
      continue
    }

    // Looping through slice of lines
    for _u, line := range lines {

      // Parse raw log line -> LogEvent
      event := parser.ParseLogLine(line, file)

      // Filter using control flow + comparison
      if event.Level.Enabled(minLevel) {
        fmt.Printf("[%s] [%v] [%s] %s\n",
          event.Timestamp.Format("15:04:05"),
          event.Level,
          event.Source,
          event.Message,
        )
      }
    }
  }
}
