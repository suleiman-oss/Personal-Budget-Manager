package config

import "github.com/spf13/viper"

type Config struct {
	DBDriver   string `mapstructure:"DBDriver"`
	DBHost     string `mapstructure:"DBHost"`
	DBPort     int64  `mapstructure:"DBPort"`
	DBUser     string `mapstructure:"DBUser"`
	DBPassword string `mapstructure:"DBPassword"`
	DBName     string `mapstructure:"DBName"`
	Port       string `mapstructure:"Port"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
