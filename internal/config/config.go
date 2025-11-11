package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string `mapstructure:"PORT"`
	DbUrl     string `mapstructure:"DB_URL"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
}

func Load() *Config {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Lỗi khi load .env")
	}

	var c Config
	viper.Unmarshal(&c)
	return &c
}
