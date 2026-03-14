# Log Aggregator

A modular, concurrent log processing tool written in Go. This application features a parallel pipeline for reading multiple log files, filtering entries, and visualizing results through a CLI or a real-time Web Dashboard.

## Features

- **Concurrent Pipeline**: Reads multiple log files in parallel using Goroutines and the Fan-In pattern.
- **Time Ordering**: Buffers and sorts log entries by timestamp for a unified chronological view.
- **Flexible Filtering**: Filter logs by level (INFO, WARN, ERROR) or keywords.
- **Web Dashboard**: Interactive real-time dashboard with secure Bcrypt authentication.
- **CLI Interface**: Robust command-line tool with statistics and export capabilities (JSON/Text).
- **Modular Design**: Clean separation of concerns across modules (Source, Parser, Filter, Ordering, Output).

## Prerequisites

- **Go**: 1.18 or higher.
- **Dependencies**: Uses `golang.org/x/crypto/bcrypt` for secure authentication.

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd LogAggregator
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o logmerge ./cmd/logmerge
   ```

## Usage

### CLI Mode

Process log files and display them in the terminal:

```bash
./logmerge --file logs/app.log --file logs/server.log --level WARN --stats
```

**Common Flags:**
- `--file`: Path to a log file (can be used multiple times).
- `--level`: Minimum log level to display (`INFO`, `WARN`, `ERROR`).
- `--match`: Keyword filter to show only logs containing specific text.
- `--stats`: Show aggregated statistics after processing.
- `--output`: Save filtered logs to a text file.
- `--json`: Export filtered logs as JSON.
- `--config`: Load configuration from a JSON file.

### Web Mode

Launch the interactive web dashboard:

```bash
./logmerge --file logs/app.log --file logs/server.log --web --port 8080 --password "secret"
```

Access the dashboard at `http://localhost:8080`.

## Testing

Run the comprehensive test suite (30+ tests):

```bash
go test ./... -v
```

## Project Structure

- `cmd/logmerge`: Application entry point and CLI logic.
- `internal/auth`: Secure Bcrypt authentication and session management.
- `internal/config`: Configuration loading and JSON exports.
- `internal/filter`: Log filtering logic and interfaces.
- `internal/models`: Data structures for log events and storage.
- `internal/ordering`: Time-based sorting and buffering.
- `internal/output`: Terminal coloring and file writing.
- `internal/parser`: Log parsing logic.
- `internal/source`: Concurrent multi-file reading logic.
- `web/`: Web server implementation and static assets.
