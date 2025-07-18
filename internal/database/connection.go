package database

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

func Migrate(db *sql.DB) error {
    migrationPath := "internal/database/migrations/001_init.sql"
    content, err := ioutil.ReadFile(migrationPath)
    if err != nil {
        // If file doesn't exist, create tables directly
        return createTables(db)
    }

    _, err = db.Exec(string(content))
    return err
}

func createTables(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS payments (
        id SERIAL PRIMARY KEY,
        correlation_id UUID UNIQUE NOT NULL,
        amount DECIMAL(10,2) NOT NULL,
        processed_by VARCHAR(10) NOT NULL,
        requested_at TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_payments_correlation_id ON payments(correlation_id);
    CREATE INDEX IF NOT EXISTS idx_payments_processed_by ON payments(processed_by);
    CREATE INDEX IF NOT EXISTS idx_payments_requested_at ON payments(requested_at);
    `
    
    _, err := db.Exec(query)
    return err
}