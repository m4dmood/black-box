package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (db *DB) Close() {
	db.Pool.Close()
}
