package collector

import (
	"testing"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestGithubAggregateEvents(t *testing.T) {
	testEventType := "CreateEvent"
	testCreatedAt := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	tt := []struct {
		name     string
		events   []*github.Event
		metrics  []statboard.Metric
		expected []statboard.Metric
	}{
		{
			name: "single event",
			events: []*github.Event{
				{
					Type:      &testEventType,
					CreatedAt: &testCreatedAt,
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testCreatedAt.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  time.Date(testCreatedAt.Year(), testCreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testCreatedAt.AddDate(0, 0, -1),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  time.Date(testCreatedAt.Year(), testCreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 1.0,
				},
			},
		},
		{
			name: "multiple events",
			events: []*github.Event{
				{
					Type:      &testEventType,
					CreatedAt: &testCreatedAt,
				},
				{
					Type:      &testEventType,
					CreatedAt: &testCreatedAt,
				},
				{
					Type:      &testEventType,
					CreatedAt: &testCreatedAt,
				},
			},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  time.Date(testCreatedAt.Year(), testCreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  time.Date(testCreatedAt.Year(), testCreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour),
					Value: 3.0,
				},
			},
		},
		{
			name:   "no events",
			events: []*github.Event{},
			metrics: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testCreatedAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testCreatedAt.Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			actual := aggregateEvents(ts.events, ts.metrics)
			assert.Equal(t, ts.expected, actual)
		})
	}
}

func TestGithubIsContribEvent(t *testing.T) {
	tt := []struct {
		name      string
		eventType string
		expected  bool
	}{
		{
			name:      "valid contribution event",
			eventType: "CommitCommentEvent",
			expected:  true,
		},
		{
			name:      "invalid contribution event",
			eventType: "FakeEvent",
			expected:  false,
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			actual := isContribEvent(ts.eventType)
			assert.Equal(t, ts.expected, actual)
		})
	}
}

func TestGithubCollect_InvalidMetric(t *testing.T) {
	c := GithubCollector{}

	_, err := c.Collect("fake_metric", 1)
	assert.Error(t, err)
}
