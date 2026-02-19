package source

import (
  "bufio"
  "os"
)

// ReadFileLines reads all lines from a file and returns them as a slice.
// LAB 3: Returns a slice
// LAB 2: Uses loops and error handling
func ReadFileLines(path string) ([]string, error) {
  // LAB 3: Slice to store all lines (dynamic size)
  var lines []string

  // LAB 1: File handle variable
  file, err := os.Open(path)

  // LAB 2: If statement for error handling
  if err != nil {
    return nil, err
  }
  defer file.Close()

  // LAB 1: Scanner variable
  scanner := bufio.NewScanner(file)

  // LAB 2: Loop with condition (for-condition pattern)
  for scanner.Scan() {
    line := scanner.Text()     // LAB 1: string variable
    lines = append(lines, line) // LAB 3: append to slice
  }

  // LAB 2: Check scan error
  if scanner.Err() != nil {
    return nil, scanner.Err()
  }

  return lines, nil
}