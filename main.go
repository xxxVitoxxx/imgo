package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xxxVitoxxx/imgo/internal/config"
	"github.com/xxxVitoxxx/imgo/pkg/imgur"
	"github.com/xxxVitoxxx/imgo/pkg/line"
	"github.com/xxxVitoxxx/imgo/pkg/replicate"
)

func main() {
	cfg := config.LoadConfig()
	imgurer := imgur.NewImgur(
		cfg.Imgur.URL,
		cfg.Imgur.ID,
	)

	replicator := replicate.NewClient(
		cfg.Replicate.URL,
		cfg.Replicate.Token,
		cfg.Replicate.Version,
		cfg.Replicate.Callback,
	)

	bot, err := line.NewLineBot(
		cfg.Line.Secret,
		cfg.Line.Token,
		imgurer,
		replicator,
	)
	if err != nil {
		log.Fatal("new bot: ", err)
	}

	router := gin.Default()
	replicator.Router(router)
	bot.Router(router)

	srv := &http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("run error: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("shutdown err: ", err)
	}

	select {
	case <-ctx.Done():
		log.Println("server exiting")
	}
}
