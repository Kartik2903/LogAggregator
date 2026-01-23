// package main

// import (
// 	"fmt"
// 	"time"

// 	"main/internal/models"
// 	"main/internal/parser"
// )

// func main() {
// //loglevel.go
// 	level := models.ParseLogLevel("ERROR")

// 	if level.Enabled(models.WARN) {
// 		fmt.Println("LogLevel check: ERROR is >= WARN")
// 	}

// //event.go
// 	event1 := models.NewLogEvent(
// 		time.Now(),
// 		level,
// 		"Manual log event created",
// 		"MAIN",
// 		"",
// 	)

// 	fmt.Println("\n--- LogEvent created using constructor ---")
// 	fmt.Println("Timestamp:", event1.Timestamp)
// 	fmt.Println("Level:", event1.Level)
// 	fmt.Println("Message:", event1.Message)
// 	fmt.Println("Source:", event1.Source)

// //parser.go
// 	rawLine := "2026-01-08T10:30:15 INFO User clicked reserve table"
// 	source := "FRONTEND"

// 	event2 := parser.ParseLogLine(rawLine, source)

// 	fmt.Println("\n--- LogEvent created using parser ---")
// 	fmt.Println("Timestamp:", event2.Timestamp)
// 	fmt.Println("Level:", event2.Level)
// 	fmt.Println("Message:", event2.Message)
// 	fmt.Println("Source:", event2.Source)
// }
