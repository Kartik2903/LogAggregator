package main

// LAB 7 (CLI & Configuration):
// Command-line interface for the Log Aggregator.
// Integrates all modules: Source → Fan-In → Parser → Filter → TimeOrdering → Output.
// Demonstrates: flag package, modular integration, concurrent pipeline.

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"logmerge/internal/config"
	"logmerge/internal/filter"
	"logmerge/internal/models"
	"logmerge/internal/ordering"
	"logmerge/internal/output"
	"logmerge/internal/parser"
	"logmerge/internal/source"
	"logmerge/web"
)

func main() {
	// ============================================================
	// LAB 1: Variable declarations + LAB 7 (CLI flag parsing)
	// ============================================================
	var files []string
	var showStats bool
	var outputFile string
	var jsonExport string
	var configFile string
	var webMode bool
	var webPort int
	var webPassword string

	// LAB 7: CLI flags using the flag package
	levelStr := flag.String("level", "INFO", "Minimum log level (INFO, WARN, ERROR)")
	matchStr := flag.String("match", "", "Keyword filter — show only logs containing this text")

	flag.BoolVar(&showStats, "stats", false, "Show statistics after processing")
	flag.StringVar(&outputFile, "output", "", "Save filtered logs to a file (e.g., output.log)")
	flag.StringVar(&jsonExport, "json", "", "Export filtered logs as JSON (e.g., output.json)")
	flag.StringVar(&configFile, "config", "", "Load configuration from a JSON file")
	flag.BoolVar(&webMode, "web", false, "Launch web dashboard instead of CLI output")
	flag.IntVar(&webPort, "port", 8080, "Port for the web dashboard (default: 8080)")
	flag.StringVar(&webPassword, "password", "", "Password to protect the web dashboard (uses bcrypt hashing)")

	// LAB 7: flag.Func for repeatable --file flags
	flag.Func("file", "Log file path (can be specified multiple times)", func(v string) error {
		files = append(files, v) // LAB 3: append to slice
		return nil
	})

	flag.Parse()

	// ============================================================
	// LAB 8: Load config from JSON if --config flag is provided
	// ============================================================
	if configFile != "" {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		// LAB 2: Conditionals — merge config values with CLI flags
		if len(files) == 0 {
			files = cfg.Files
		}
		if *levelStr == "INFO" && cfg.Level != "" {
			*levelStr = cfg.Level
		}
		if *matchStr == "" && cfg.Match != "" {
			*matchStr = cfg.Match
		}
		if outputFile == "" {
			outputFile = cfg.OutputFile
		}
		showStats = showStats || cfg.ShowStats
	}

	// ============================================================
	// LAB 2: Input validation
	// ============================================================
	if len(files) == 0 {
		fmt.Println("No files provided. Usage:")
		fmt.Println("  logmerge --file app.log --file server.log [--level WARN] [--match keyword]")
		fmt.Println("  logmerge --file app.log --file server.log --web")
		flag.PrintDefaults()
		return
	}

	// ============================================================
	// Web mode: launch dashboard and exit
	// ============================================================
	if webMode {
		server := web.NewServer(files, webPassword)
		if err := server.Start(webPort); err != nil {
			fmt.Println("Web server error:", err)
			os.Exit(1)
		}
		return
	}

	// ============================================================
	// LAB 1: Custom type conversion
	// ============================================================
	minLevel := models.ParseLogLevel(*levelStr)

	// LAB 2: Switch statement for validation
	switch minLevel {
	case models.INFO, models.WARN, models.ERROR:
		// valid
	default:
		minLevel = models.INFO
	}

	// ============================================================
	// Build the filter chain
	// LAB 6: Interface-based filter composition
	// ============================================================
	var filters []filter.Filter

	// LAB 6: Level filter always active
	filters = append(filters, filter.LevelFilter{MinLevel: minLevel})

	// LAB 6: Keyword filter if --match specified
	if *matchStr != "" {
		filters = append(filters, filter.KeywordFilter{Keyword: *matchStr})
	}

	// LAB 6: Chain all filters into a composite filter
	activeFilter := filter.ChainFilters(filters...)

	// ============================================================
	// Print header
	// ============================================================
	output.PrintHeader()
	fmt.Printf("  Files:  %s\n", strings.Join(files, ", "))
	fmt.Printf("  Level:  %s\n", minLevel)
	if *matchStr != "" {
		fmt.Printf("  Match:  %q\n", *matchStr)
	}
	output.PrintSeparator()
	fmt.Println()

	// ============================================================
	// LAB 10: Concurrent pipeline — read all files concurrently
	// LAB 9:  Fan-In merges all channels via WaitGroup
	// ============================================================
	mergedLines, fileErrs := source.ReadMultipleConcurrent(files)

	// ============================================================
	// LAB 5: Parser interface + LAB 6: Interface usage
	// ============================================================
	var p parser.Parser = parser.SimpleParser{}

	// LAB 7: Time ordering buffer with 1-second window
	buffer := ordering.NewTimeOrderBuffer(1 * time.Second)

	// LAB 4: Storage for statistics
	storage := models.NewLogStorage()

	// Collect all filtered events for file/JSON export
	var allFiltered []models.LogEvent

	// ============================================================
	// LAB 10: Range over merged channel — process each line as it arrives
	// ============================================================
	for line := range mergedLines {
		// LAB 10: Split source prefix from line (added by ReadMultipleConcurrent)
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		sourceName := parts[0]
		rawLine := parts[1]

		// LAB 5: Parse raw line into structured LogEvent
		event, err := p.ParseLogLine(rawLine, sourceName)
		if err != nil {
			continue // skip unparseable lines
		}

		// LAB 4: Store event for statistics
		storage.AddEvent(event)

		// LAB 7: Add to time-ordered buffer
		buffer.AddEvent(event)

		// LAB 6: Apply filter chain
		if activeFilter.Match(event) {
			// LAB 5: Color-coded terminal output
			output.PrintEventWithSource(event)
			allFiltered = append(allFiltered, event) // LAB 3: append
		}
	}

	// ============================================================
	// Flush any remaining buffered events and sort them
	// LAB 3: Slice operations in time ordering
	// ============================================================
	remaining := buffer.Flush()
	_ = remaining // events already printed in real-time above

	// ============================================================
	// LAB 5: Handle file-reading errors
	// ============================================================
	for err := range fileErrs {
		if err != nil {
			fmt.Printf("\n%sSource error: %v%s\n", output.ColorRed, err, output.ColorReset)
		}
	}

	// ============================================================
	// Output summary
	// ============================================================
	fmt.Println()
	output.PrintSeparator()
	fmt.Printf("  Processed: %d total events | Displayed: %d filtered events\n",
		storage.TotalEvents(), len(allFiltered))
	output.PrintSeparator()

	// ============================================================
	// LAB 5: Save to file if --output specified
	// ============================================================
	if outputFile != "" {
		// Sort before saving
		ordering.SortEvents(allFiltered)
		if err := output.WriteToFile(allFiltered, outputFile); err != nil {
			fmt.Println("Error writing output:", err)
		} else {
			fmt.Printf("  Saved to: %s\n", outputFile)
		}
	}

	// ============================================================
	// LAB 8: Export as JSON if --json specified
	// ============================================================
	if jsonExport != "" {
		ordering.SortEvents(allFiltered)
		if err := config.ExportEventsJSON(allFiltered, jsonExport); err != nil {
			fmt.Println("Error exporting JSON:", err)
		}
	}

	// ============================================================
	// Statistics — LAB 4: Map iteration
	// ============================================================
	if showStats {
		printStatistics(storage)
	}
}

// printStatistics prints aggregated log statistics.
// LAB 4: Map iteration and struct access
func printStatistics(storage *models.LogStorage) {
	fmt.Println()
	output.PrintSeparator()
	fmt.Printf("%s%s  STATISTICS  %s\n", output.ColorBold, output.ColorCyan, output.ColorReset)
	output.PrintSeparator()

	// LAB 4: Iterate over level statistics map
	stats := storage.GetLevelStatistics()
	for level, count := range stats {
		fmt.Printf("  %-8s %d events\n", level, count)
	}

	fmt.Println()

	// LAB 4: Iterate over source metadata
	sources := storage.GetAllSources()
	fmt.Printf("  Sources: %d\n", len(sources))
	for idx, src := range sources {
		events := storage.GetEventsBySource(src)
		meta, exists := storage.GetSourceMetadata(src)

		fmt.Printf("  %d. %s\n", idx+1, src)
		fmt.Printf("     Events: %d\n", len(events))
		if exists {
			fmt.Printf("     Errors: %d\n", meta.ErrorCount)
		}
	}
	output.PrintSeparator()
}
