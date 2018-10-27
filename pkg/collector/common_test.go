package collector

import (
	"testing"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/stretchr/testify/assert"
)

func TestCommonGenerateEmptyMetrics(t *testing.T) {

	tt := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected []statboard.Metric
	}{
		{
			name:  "range of dates",
			start: time.Now().AddDate(0, 0, -1),
			end:   time.Now(),
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  time.Now().Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
		},
		{
			name:  "equal date",
			start: time.Now().AddDate(0, 0, -1),
			end:   time.Now().AddDate(0, 0, -1),
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour),
					Value: 0.0,
				},
			},
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			actual := generateEmptyMetrics("testMetric", ts.start, ts.end)
			assert.Equal(t, ts.expected, actual)
		})
	}
}
