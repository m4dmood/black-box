package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m4dmood/black-box/internal/parser"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context) (*DB, error) {
	connStr := "postgres://admin:password@localhost:5432/blackbox_store"

	// Connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Printf("[DATABASE] Error while creating pool -> %v", err)
		return nil, err
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		fmt.Printf("[DATABASE] Ping failed -> %v", err)
		return nil, err
	}

	fmt.Println("[DATABASE] Database connected successfully! (Pool)")
	return &DB{Pool: pool}, nil
}

func (db *DB) InsertWithFallback(ctx context.Context, entry parser.RegistryEntry) error {

	query := `INSERT INTO registry (ts, device_id, event_type, val, event_message, level)
				VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Pool.Exec(ctx, query,
		entry.Timestamp,
		entry.DeviceID,
		"data",
		entry.Value,
		entry.EventMessage,
		entry.Level,
	)

	if err != nil {
		fmt.Printf("[DATABASE] Error while inserting data, moving record to error log -> %v\n", err)
		db.logErrorToTable(ctx, entry.ToRawString(), err.Error())
		return err
	}

	fmt.Printf("[DATABASE] Logged event -> %v\n", entry)
	return nil
}

func (db *DB) InsertBatch(ctx context.Context, entries []parser.RegistryEntry) (int64, error) {
	// Definiamo quali colonne vogliamo mappare e da dove vengono
	columns := []string{"ts", "device_id", "event_type", "val", "event_message", "level"}

	// Trasformiamo la nostra slice di struct in una matrice di dati che pgx capisce
	rows := [][]any{}
	for _, e := range entries {
		rows = append(rows, []any{e.Timestamp, e.DeviceID, "data", e.Value, e.EventMessage, e.Level})
	}

	// CopyFrom è un'operazione atomica: o entrano tutte o nessuna (molto veloce)
	count, err := db.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"registry"}, // Nome tabella
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return 0, fmt.Errorf("[DATABASE] Batch insert failed: %w", err)
	}

	return count, nil
}

func (db *DB) logErrorToTable(ctx context.Context, rawData string, reason string) {
	fmt.Printf("[DATABASE] Moving record to error log -> %s\n", rawData)
	fmt.Printf("[DATABASE] Invalid record due to %s\n", reason)

	query := `INSERT INTO registry_errors (raw_data, error_message) VALUES ($1, $2)`
	db.Pool.Exec(ctx, query, rawData, reason)
}

func (db *DB) Close() {
	db.Pool.Close()
}
