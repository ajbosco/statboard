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
	Metrics map[string]map[string]metricConfig `mapstructure:"metrics"`
}

type fitbitConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CacheFile    string `mapstructure:"cache_file"`
}

type metricConfig struct {
	Color    string `mapstructure:"color"`
	DaysBack int    `mapstructure:"days_back"`
}

// func main() {

// 	viper.SetConfigName("test")
// 	viper.AddConfigPath(".")
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		log.Panic("error:", err)
// 	}

// 	var config Config

// 	err = viper.Unmarshal(&config)
// 	if err != nil {
// 		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
// 	}
// 	// fmt.Println(config)

// 	fmt.Println(config.Db.FilePath)
// 	fmt.Println(config.Charts.DirPath)
// 	fmt.Println(config.Fitbit.ClientID)
// 	fmt.Println(config.Fitbit.CacheFile)
// 	// fmt.Println(config.Metrics["fitbit"]["steps"].color)
// 	// fmt.Println(config.Metrics["fitbit"]["other_met"].daysBack)

// 	for k, m := range config.Metrics {
// 		fmt.Println(k)
// 		for k, v := range m {
// 			fmt.Println(k)
// 			fmt.Println(v.Color)
// 			fmt.Println(v.DaysBack)
// 		}
// 	}
// }
