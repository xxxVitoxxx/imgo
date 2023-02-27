package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xxxVitoxxx/imgo/imgur"
	"github.com/xxxVitoxxx/imgo/line"
	"github.com/xxxVitoxxx/imgo/replicate"
)

type config struct {
	Replicate Replicate `toml:"replicate"`
	Imgur     Imgur     `toml:"imgur"`
	Line      Line      `toml:"line"`
}

type Replicate struct {
	URL     string `toml:"url"`
	Token   string `toml:"token"`
	Version string `toml:"version"`
}

type Imgur struct {
	URL string `toml:"url"`
	ID  string `toml:"id"`
}

type Line struct {
	Secret string `toml:"secret"`
	Token  string `toml:"token"`
}

var cfg config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("read config: ", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("unmarshal config: ", err)
	}
}

func main() {
	imgurer := imgur.NewImgur(
		cfg.Imgur.URL,
		cfg.Imgur.ID,
	)

	replicater := replicate.NewClient(
		cfg.Replicate.URL,
		cfg.Replicate.Token,
		cfg.Replicate.Version,
	)

	bot, err := line.NewLineBot(
		cfg.Line.Secret,
		cfg.Line.Token,
		imgurer,
		replicater,
	)
	if err != nil {
		log.Fatal("new bot: ", err)
	}

	r := gin.Default()
	replicater.Router(r)
	bot.Router(r)

	if err := r.Run(":8090"); err != nil {
		log.Fatal("run error: ", err)
	}
}
