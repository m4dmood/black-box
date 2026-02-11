package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m4dmood/black-box/internal/parser"
)

// DB è il wrapper per il pool di connessioni
type DB struct {
	Pool *pgxpool.Pool
}

// Connect inizializza la connessione al database
func Connect(ctx context.Context) (*DB, error) {
	// In un progetto reale, questa stringa verrebbe da una variabile d'ambiente
	connStr := "postgres://admin:password@localhost:5432/blackbox_store"

	// Creiamo un pool di connessioni invece di una singola connessione
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("impossibile creare il pool: %w", err)
	}

	// Verifichiamo se il database è effettivamente raggiungibile
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping fallito: %w", err)
	}

	fmt.Println("Connessione a Postgres stabilita con successo (Pool)")
	return &DB{Pool: pool}, nil
}

func (db *DB) InsertWithFallback(ctx context.Context, entry parser.RegistryEntry) error {

	query := `ÌNSERT INTO registry (ts, device_id, event_type, val, event_message, level)
				VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Pool.Exec(ctx, query,
		entry.Timestamp,
		entry.DeviceID,
		"data", // Per ora usiamo un valore fisso, ma in futuro potrebbe essere dinamico
		entry.Value,
		entry.EventMessage,
		entry.Level,
	)

	if err != nil {
		db.logErrorToTable(ctx, entry.ToRawString(), err.Error())
	}

	return nil
}

func (db *DB) logErrorToTable(ctx context.Context, rawData string, reason string) {
	query := `INSERT INTO registry_errors (raw_data, error_message) VALUES ($1, $2)`
	db.Pool.Exec(ctx, query, rawData, reason)
}

func (db *DB) Close() {
	db.Pool.Close()
}
