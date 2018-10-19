package collector

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sajal/fitbitclient"
	"github.com/spf13/viper"
)

var _ Collector = &FitbitCollector{}

// FitbitCollector is used to collect metrics from Fitbit API and implements Collector interface
type FitbitCollector struct {
	client *http.Client
}

// NewFitbitCollector parses config file and creates a new FitbitCollector
func NewFitbitCollector(configPath string) (*FitbitCollector, error) {
	viper.SetConfigName(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not read config from filepath: %v", configPath))
	}

	clientID := viper.GetString("fitbit.client_id")
	if clientID == "" {
		return nil, errors.New("'fitbit.client_id' must be present")
	}
	clientSecret := viper.GetString("fitbit.client_secret")
	if clientSecret == "" {
		return nil, errors.New("'fitbit.client_secret' must be present")
	}
	cacheFile := viper.GetString("fitbit.cache_file")
	if cacheFile == "" {
		return nil, errors.New("'fitbit.cache_file' must be present")
	}

	cfg := &fitbitclient.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"activity"},
		CredFile:     cacheFile,
	}
	client, err := fitbitclient.NewFitBitClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not create fitbit client")
	}

	return &FitbitCollector{client: client}, nil
}

func (c *FitbitCollector) Collect(metricName string, daysBack int) (*metric, error) {

	return nil, nil
}
