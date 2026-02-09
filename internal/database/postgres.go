package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context) (*DB, error) {
	connStr := "postgres://admin:password@localhost:5432/blackbox_store"

	// Connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("[DATABASE] Error while creating pool -> %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("[DATABASE] Ping failed -> %w", err)
	}

	fmt.Println("[DATABASE] Database connected successfully! (Pool)")
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
