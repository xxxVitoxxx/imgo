package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Replicate Replicate `toml:"replicate"`
	Imgur     Imgur     `toml:"imgur"`
	Line      Line      `toml:"line"`
	RabbitMQ  RabbitMQ  `toml:"rabbitmq"`
	Redis     Redis     `toml:"redis"`
	Bucket    Bucket    `toml:"bucket"`
}

type Replicate struct {
	URL      string `toml:"url"`
	Token    string `toml:"token"`
	Version  string `toml:"version"`
	Callback string `toml:"callback"`
}

type Imgur struct {
	URL string `toml:"url"`
	ID  string `toml:"id"`
}

type Line struct {
	Secret string `toml:"secret"`
	Token  string `toml:"token"`
}

type RabbitMQ struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Address  string `toml:"address"`
}

type Redis struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Address  string `toml:"address"`
	DB       string `toml:"db"`
}

type Bucket struct {
	ID     string `toml:"id"`
	Secret string `toml:"secret"`
	Name   string `toml:"name"`
	Region string `toml:"region"`
}

func loadConfig() (cfg config) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("read config: ", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("unmarshal config: ", err)
	}
	return
}

// LoadConfig load config from toml file
func LoadConfig() config {
	return loadConfig()
}
