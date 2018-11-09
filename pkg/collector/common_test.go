package collector

import (
	"testing"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/stretchr/testify/assert"
)

func TestCommonGenerateEmptyMetrics(t *testing.T) {
	testMonthDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour)

	tt := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected []statboard.Metric
	}{
		{
			name:  "monthly - range of dates",
			start: testMonthDate.AddDate(0, -2, 0),
			end:   testMonthDate,
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testMonthDate.AddDate(0, -2, 0),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  testMonthDate.AddDate(0, -1, 0),
					Value: 0.0,
				},
				{
					Name:  "testMetric",
					Date:  testMonthDate,
					Value: 0.0,
				},
			},
		},
		{
			name:  "monthly - equal date",
			start: testMonthDate,
			end:   testMonthDate,
			expected: []statboard.Metric{
				{
					Name:  "testMetric",
					Date:  testMonthDate,
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
