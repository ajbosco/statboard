package config

// Config contains information for running Statboard
type Config struct {
	Db struct {
		FilePath string `mapstructure:"file_path"`
	} `mapstructure:"db"`
	Charts struct {
		DirPath string `mapstructure:"dir_path"`
	} `mapstructure:"charts"`
	Fitbit  fitbitConfig                       `mapstructure:"fitbit"`
	Github  githubConfig                       `mapstructure:"github"`
	Metrics map[string]map[string]metricConfig `mapstructure:"metrics"`
}

type fitbitConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CacheFile    string `mapstructure:"cache_file"`
}

type githubConfig struct {
	Username    string `mapstructure:"user_name"`
	AccessToken string `mapstructure:"access_token"`
}

type metricConfig struct {
	ChartName string `mapstructure:"chart_name"`
	Color     string `mapstructure:"color"`
	DaysBack  int    `mapstructure:"days_back"`
}
