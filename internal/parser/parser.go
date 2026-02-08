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
		fmt.Errorf("Errore apertura file : %s\n", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	fmt.Printf("-------- Contenuto file %s --------\n", path)

	for scanner.Scan() {
		// Cast riga file -> string
		line := scanner.Text()

		// Skip riga se vuota
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Parsing campi nella riga
		parts := strings.Split(line, ";")

		if len(parts) < 4 {
			fmt.Printf("Formato riga non valido: %s\n", line)
			continue
		}

		id := parts[0]
		value := parts[1]
		timestamp := parts[2]
		status := parts[3]

		fmt.Printf("[DATA] Codice sensore: %s | Valore rilevato: %s | Timestamp: %s | Status: %s\n", id, value, timestamp, status)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Errore durante la lettura: %w", err)
	}

	return nil
}
