package source

import (
  "bufio"
  "os"
)

// Interface (Lab 6)
type Source interface {
  Read(path string) (<-chan string, <-chan error)
}

// File implementation
type FileSource struct{}

// LAB 5: Function using error handling
func (fs FileSource) Read(path string) (<-chan string, <-chan error) {

  lines := make(chan string)
  errs := make(chan error, 1)

  // LAB 5: goroutine function
  go func() {

    defer close(lines)
    defer close(errs)

    // LAB 5: open file + handle error
    file, err := os.Open(path)
    if err != nil {
      errs <- err
      return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // LAB 2: loop control flow
    for scanner.Scan() {
      lines <- scanner.Text()
    }

    // LAB 5: scanner error handling
    if err := scanner.Err(); err != nil {
      errs <- err
    }
  }()

  return lines, errs
}