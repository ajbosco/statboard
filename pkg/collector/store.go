package collector

import (
	"database/sql"
	"time"
)

// Store represents the store that writes and reads metrics
type Store interface {
	GetMetric(source string, since time.Time) ([]metric, error)
	WriteMetric(m metric) error
	Close() error
}

type dbStore struct {
	db *sql.DB
}

// NewDBStore creates a new Collector db instance
func NewDBStore(db *sql.DB) (Store, error) {
	return &dbStore{db: db}, nil
}

// WriteMetric inserts or updates metric in database
func (d *dbStore) WriteMetric(m metric) error {
	return nil
}

// GetMetric returns the metrics since the given date
func (d *dbStore) GetMetric(source string, since time.Time) ([]metric, error) {
	return nil, nil
}

// Close closes the database connection
func (d *dbStore) Close() error {
	return nil
}
