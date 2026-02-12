package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)

func (r *RegistryEntry) ToRawString() string {
	return "{ DeviceID: " + r.DeviceID +
		", Timestamp: " + r.Timestamp +
		", Value: " + strconv.FormatFloat(r.Value, 'f', -1, 64) +
		", Level: " + r.Level +
		", EventMessage: " + r.EventMessage + " }"
}

func IsValidUUID(u string) bool {
	return uuidRegex.MatchString(u)
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

		if !IsValidUUID(parts[0]) {
			fmt.Printf("[PARSER] Salto riga: UUID malformato -> %s\n", parts[0])
			continue // Salta alla prossima riga senza mandarla al DB
		}

		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("[PARSER] Non-numeric value in row: %s", value)
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
