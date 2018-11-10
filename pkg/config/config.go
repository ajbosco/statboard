package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config contains information for running Statboard
type Config struct {
	Fitbit    fitbitConfig                       `mapstructure:"fitbit" yaml:"fitbit"`
	Github    githubConfig                       `mapstructure:"github" yaml:"github"`
	Goodreads goodreadsConfig                    `mapstructure:"goodreads" yaml:"goodreads"`
	Metrics   map[string]map[string]MetricConfig `mapstructure:"metrics" yaml:"metrics"`
}

type fitbitConfig struct {
	ClientID     string `mapstructure:"client_id" yaml:"client_id"`
	ClientSecret string `mapstructure:"client_secret" yaml:"client_secret"`
	CacheFile    string `mapstructure:"cache_file" yaml:"cache_file"`
}

type githubConfig struct {
	Username    string `mapstructure:"user_name" yaml:"user_name"`
	AccessToken string `mapstructure:"access_token" yaml:"access_token"`
}

type goodreadsConfig struct {
	DeveloperKey    string `mapstructure:"developer_key" yaml:"developer_key"`
	DeveloperSecret string `mapstructure:"developer_secret" yaml:"developer_secret"`
	AccessToken     string `mapstructure:"access_token" yaml:"access_token"`
	AccessSecret    string `mapstructure:"access_secret" yaml:"access_secret"`
}

// MetricConfig contains information for collecting and visualizing a metric
type MetricConfig struct {
	ChartName         string `mapstructure:"chart_name" yaml:"chart_name"`
	ChartColor        string `mapstructure:"chart_color" yaml:"chart_color"`
	ChartMonthsBack   int    `mapstructure:"chart_months_back" yaml:"chart_months_back"`
	CollectMonthsBack int    `mapstructure:"collect_months_back" yaml:"collect_months_back"`
}

// Write writes a Config object to the config file
func (c *Config) Write() error {
	// get config file path from env variable
	path := os.Getenv("STATBOARD_CONFIGFILEPATH")

	configBytes, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}
	if err = ioutil.WriteFile(path, configBytes, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not write config to filepath: %v", path))
	}
	return nil
}
