package collector

import (
	"testing"
	"time"

	"github.com/ajbosco/reads/goodreads"
	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/stretchr/testify/assert"
)

func TestGoodreadsAggregateBooks(t *testing.T) {
	testReadAt := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)

	tt := []struct {
		name     string
		books    []goodreads.Book
		metrics  []statboard.Metric
		expected []statboard.Metric
	}{
		{
			name: "single book",
			books: []goodreads.Book{
				{
					ReadAt: "Mon Jan 02 15:04:05 -0700 2006",
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 1.0,
				},
			},
		},
		{
			name: "multiple events",
			books: []goodreads.Book{
				{
					ReadAt: "Mon Jan 02 15:04:05 -0700 2006",
				},
				{
					ReadAt: "Mon Jan 03 15:04:05 -0700 2006",
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 2.0,
				},
			},
		},
		{
			name:  "no events",
			books: []goodreads.Book{},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			actual, _ := aggregateBooks(ts.books, ts.metrics)
			assert.Equal(t, ts.expected, actual)
		})
	}
}

func TestGoodreadsAggregatePages(t *testing.T) {
	testReadAt := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)

	tt := []struct {
		name     string
		books    []goodreads.Book
		metrics  []statboard.Metric
		expected []statboard.Metric
	}{
		{
			name: "single book",
			books: []goodreads.Book{
				{
					ReadAt: "Mon Jan 02 15:04:05 -0700 2006",
					Pages:  "100",
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 100.0,
				},
			},
		},
		{
			name: "multiple events",
			books: []goodreads.Book{
				{
					ReadAt: "Mon Jan 02 15:04:05 -0700 2006",
					Pages:  "200",
				},
				{
					ReadAt: "Mon Jan 03 15:04:05 -0700 2006",
					Pages:  "200",
				},
				{
					ReadAt: "Mon Jan 03 15:04:05 -0700 2016",
					Pages:  "200",
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 400.0,
				},
			},
		},
		{
			name:  "no events",
			books: []goodreads.Book{},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testReadAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			actual, _ := aggregatePages(ts.books, ts.metrics)
			assert.Equal(t, ts.expected, actual)
		})
	}
}

func TestGoodreadsCollect_InvalidMetric(t *testing.T) {
	c := GoodreadsCollector{}

	_, err := c.Collect("fake_metric", 1, "month")
	assert.Error(t, err)
}
