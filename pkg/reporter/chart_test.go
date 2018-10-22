package reporter

import (
	"os"
	"testing"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/stretchr/testify/assert"
)

func TestRenderChart(t *testing.T) {
	defer os.Remove("test-chart.png")

	metrics := []statboard.Metric{
		statboard.Metric{
			Name:  "testMetric",
			Date:  time.Now(),
			Value: 3.0,
		},
		statboard.Metric{
			Name:  "testMetric",
			Date:  time.Now().AddDate(0, 0, -1),
			Value: 4.0,
		},
	}

	err := RenderChart("test", "2A7087", ".", metrics)
	assert.NoError(t, err)
}

func TestRenderChart_OneRecord(t *testing.T) {
	defer os.Remove("test-chart.png")

	metrics := []statboard.Metric{
		statboard.Metric{
			Name:  "testMetric",
			Date:  time.Now(),
			Value: 3.0,
		},
	}

	err := RenderChart("test", "2A7087", ".", metrics)
	assert.Error(t, err)
}

func TestRenderChart_RepeatYValue(t *testing.T) {
	defer os.Remove("test-chart.png")

	metrics := []statboard.Metric{
		statboard.Metric{
			Name:  "testMetric",
			Date:  time.Now(),
			Value: 3.0,
		},
		statboard.Metric{
			Name:  "testMetric",
			Date:  time.Now().AddDate(0, 1, -1),
			Value: 3.0,
		},
	}

	err := RenderChart("test", "2A7087", ".", metrics)
	assert.Error(t, err)
}

func TestRenderChart_RepeatXValue(t *testing.T) {
	defer os.Remove("test-chart.png")

	metrics := []statboard.Metric{
		statboard.Metric{
			Name:  "testMetric",
			Date:  time.Now(),
			Value: 3.0,
		},
		statboard.Metric{
			Name:  "testMetric",
			Date:  time.Now(),
			Value: 4.0,
		},
	}

	err := RenderChart("test", "2A7087", ".", metrics)
	assert.Error(t, err)
}
