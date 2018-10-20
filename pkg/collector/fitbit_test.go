package collector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFitbitCollect_Steps(t *testing.T) {
	server := mockServer("activities-steps.json")

	c := FitbitCollector{client: server.Client(), baseURI: server.URL}

	metrics, err := c.Collect("steps", 1)
	if err != nil {
		t.Fatal(err)
	}

	expectedDt, err := time.Parse("2006-01-02", "2011-04-30")
	assert.NoError(t, err)

	assert.Equal(t, 7, len(metrics))
	assert.Equal(t, float64(2344), metrics[1].Value)
	assert.Equal(t, "fitbit.steps", metrics[2].Name)
	assert.Equal(t, expectedDt, metrics[3].Date)
}

func TestFitbitCollect_InvalidMetric(t *testing.T) {
	c := FitbitCollector{}

	_, err := c.Collect("fake_metric", 1)
	assert.Error(t, err)
}
