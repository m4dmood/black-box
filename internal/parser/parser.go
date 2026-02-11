package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RegistryEntry struct {
	DeviceID     string
	Timestamp    string
	Value        float64
	Level        string
	EventMessage string
}

func (r *RegistryEntry) ToRawString() string {
	return "{ DeviceID: " + r.DeviceID +
		", Timestamp: " + r.Timestamp +
		", Value: " + strconv.FormatFloat(r.Value, 'f', -1, 64) +
		", Level: " + r.Level +
		", EventMessage: " + r.EventMessage + " }"
}

func ProcessFile(path string) ([]RegistryEntry, error) {
	var entries []RegistryEntry
	file, err := os.Open(path)

	if err != nil {
		return entries, fmt.Errorf("[PARSER] Failed to open file -> %s\n", err)
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

		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("[DATA] Non-numeric value in row: %s", value)
			continue
		}

		if err := scanner.Err(); err != nil {
			return entries, fmt.Errorf("[PARSER] Error while reading data: %w", err)
		}

		fmt.Printf("[PARSER] Device code: %s | Value: %s | Timestamp: %s | Status: %s\n", id, value, timestamp, status)
		entries = append(entries, RegistryEntry{
			id,
			timestamp,
			val,
			"INFO",
			status,
		})
	}

	return entries, nil
}
