package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Replicate Replicate `toml:"replicate"`
	Imgur     Imgur     `toml:"imgur"`
	Line      Line      `toml:"line"`
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
