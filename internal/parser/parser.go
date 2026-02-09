package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RegistryEntry struct {
	DeviceID     string
	Timestamp    string
	Value        float64
	Level        string
	EventMessage string
}

func ProcessFile(path string) error {
	file, err := os.Open(path)

	if err != nil {
		fmt.Errorf("[PARSER] Failed to open file -> %s\n", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	fmt.Printf("-------- File content %s --------\n", path)

	for scanner.Scan() {
		// Cast row to string
		line := scanner.Text()

		// Skip empty rows
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Parsing single fields
		parts := strings.Split(line, ";")

		if len(parts) < 4 {
			fmt.Printf("[PARSER] Invalid line format: %s\n", line)
			continue
		}

		id := parts[0]
		value := parts[1]
		timestamp := parts[2]
		status := parts[3]

		fmt.Printf("[PARSER] Device code: %s | Value: %s | Timestamp: %s | Status: %s\n", id, value, timestamp, status)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("[PARSER] Error while reading data: %w", err)
	}

	return nil
}
