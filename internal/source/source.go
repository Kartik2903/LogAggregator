package source

import (
  "bufio"
  "os"
)

// ReadFileLines reads all lines from a file and returns them as a slice.
// This uses variables, slices, loops, and if-based error handling.
func ReadFileLines(path string) ([]string, error) {

  // Slice to store all lines (dynamic size)
  var lines []string

  // Open the file
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  // Read file line by line
  scanner := bufio.NewScanner(file)
  for scanner.Scan() { // loop control flow
    line := scanner.Text()     // variable
    lines = append(lines, line) // slice append
  }

  // Check scan error
  if scanner.Err() != nil {
    return nil, scanner.Err()
  }

  return lines, nil
}
