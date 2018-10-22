package storage

import (
	"fmt"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Store represents the store that writes and reads metrics
type Store interface {
	GetMetric(name string, since time.Time) ([]statboard.Metric, error)
	WriteMetric(m statboard.Metric) error
	Close() error
}

type stormStore struct {
	db *storm.DB
}

// NewStormStore creates a new metric db instance
func NewStormStore(filePath string) (Store, error) {
	db, err := storm.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to open db with filepath %q", filePath))
	}
	return &stormStore{db: db}, nil
}

// WriteMetric inserts or updates metric in database
func (s *stormStore) WriteMetric(m statboard.Metric) error {
	m.ID = fmt.Sprintf("%s-%s", m.Date, m.Name)
	return s.db.Save(&m)
}

// GetMetric returns the metrics since the given date
func (s *stormStore) GetMetric(name string, since time.Time) ([]statboard.Metric, error) {
	var metrics []statboard.Metric
	err := s.db.Select(q.And(q.Eq("Name", name), q.Gt("Date", since))).Find(&metrics)
	if err != nil {
		if err == storm.ErrNotFound {
			return metrics, nil
		}
		return nil, err
	}
	logrus.Info(fmt.Sprintf("found %d %q records", len(metrics), name))
	return metrics, nil
}

// Close closes the database connection
func (s *stormStore) Close() error {
	return s.db.Close()
}
