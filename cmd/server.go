package cmd

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
	"github.com/spf13/cobra"
	"github.com/xxxVitoxxx/imgo/internal/config"
	"github.com/xxxVitoxxx/imgo/pkg/bucket"
	"github.com/xxxVitoxxx/imgo/pkg/line"
	"github.com/xxxVitoxxx/imgo/pkg/rabbitmq"
	"github.com/xxxVitoxxx/imgo/pkg/replicate"
	"github.com/xxxVitoxxx/imgo/pkg/storage"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start imgo server to improve the clarity of the photo through the AI model",
	Long: `start the imgo server to receive photos provided by user from line and 
improve the clarity of the photo through the AI model.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()

		mq, err := rabbitmq.NewMessageQueue(
			cfg.RabbitMQ.User,
			cfg.RabbitMQ.Password,
			cfg.RabbitMQ.Address,
		)
		if err != nil {
			log.Fatal(err)
		}

		replicateQueue, err := mq.DeclareReplicateQueue()
		if err != nil {
			log.Fatal(err)
		}

		replicator := replicate.NewClient(
			cfg.Replicate.URL,
			cfg.Replicate.Token,
			cfg.Replicate.Version,
			cfg.Replicate.Callback,
			mq,
			replicateQueue,
		)

		redis, err := storage.NewRedis(
			cfg.Redis.User,
			cfg.Redis.Password,
			cfg.Redis.Address,
			cfg.Redis.DB,
		)
		if err != nil {
			log.Fatal(err)
		}

		bucket, err := bucket.NewBucket(
			cfg.Bucket.ID,
			cfg.Bucket.Secret,
			cfg.Bucket.Name,
			cfg.Bucket.Region,
		)
		if err != nil {
			log.Fatal(err)
		}

		bot, err := line.NewLineBot(
			cfg.Line.Secret,
			cfg.Line.Token,
			replicator,
			bucket,
			redis,
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
	},
}
