package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

func OpenOracle(dsn string) (*sql.DB, error) {
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("oracle ping: %w", err)
	}
	return db, nil
}
