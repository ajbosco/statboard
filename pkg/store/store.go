package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	bolt "github.com/etcd-io/bbolt"
	"github.com/pkg/errors"
)

const (
	bucketName = "Metrics"
)

// Store represents the store that writes and reads metrics
type Store interface {
	GetMetric(metricName string, since time.Time) ([]statboard.Metric, error)
	WriteMetric(m statboard.Metric) error
	Close() error
}

type boltStore struct {
	db *bolt.DB
}

// NewBoltStore creates a new Collector BoltDB instance
func NewBoltStore(filePath string) (Store, error) {
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to open BoltDB with filepath %q", filePath))
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return errors.Wrap(err, "failed to create bucket")
		}
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &boltStore{db: db}, nil
}

// WriteMetric inserts or updates metric in database
func (b *boltStore) WriteMetric(m statboard.Metric) error {
	key := fmt.Sprintf("%s-%s", m.Date, m.Name)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return errors.Wrap(err, "encoding metric failed")
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Put([]byte(key), buf.Bytes())
	})
}

// GetMetric returns the metrics since the given date
func (b *boltStore) GetMetric(metricName string, since time.Time) ([]statboard.Metric, error) {
	var metrics []statboard.Metric
	key := fmt.Sprintf("%s-%s", since, metricName)

	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		search := []byte(key)

		for k, v := c.Seek(search); k != nil; k, v = c.Next() {
			met := statboard.Metric{}
			dec := gob.NewDecoder(bytes.NewBuffer(v))
			err := dec.Decode(&met)
			if err != nil {
				return errors.Wrap(err, "decoding metric failed")
			}
			metrics = append(metrics, met)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to query metric:%q", metricName))
	}

	return metrics, nil
}

// Close closes the database connection
func (b *boltStore) Close() error {
	return b.db.Close()
}
