package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xxxVitoxxx/imgo/internal/config"
	"github.com/xxxVitoxxx/imgo/pkg/bucket"
	"github.com/xxxVitoxxx/imgo/pkg/line"
	"github.com/xxxVitoxxx/imgo/pkg/rabbitmq"
	"github.com/xxxVitoxxx/imgo/pkg/replicate"
	"github.com/xxxVitoxxx/imgo/pkg/storage"
)

func init() {
	rootCmd.AddCommand(consumerCmd)
}

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "start a consumer to handle messages from queue of replicate",
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

		messages, err := mq.Receive(replicateQueue)
		if err != nil {
			log.Fatal(err)
		}

		replicator := replicate.NewClient(
			cfg.Replicate.URL,
			cfg.Replicate.Token,
			cfg.Replicate.Version,
			cfg.Replicate.Callback,
			nil,
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

		go func() {
			for message := range messages {
				var prediction replicate.Prediction
				if err := json.Unmarshal(message.Body, &prediction); err != nil {
					log.Println("failed to unmarshal message: ", err)
				}
				fmt.Println(prediction)

				// use replicate id to get line token from redis
				token, err := redis.GetData(context.Background(), prediction.ID)
				if err != nil {
					log.Println("failed to get token from replicate id: ", err)
				}

				switch prediction.Status {
				case replicate.PredictionSuccess:
					if err := bot.ReplyImage(token, prediction.URL); err != nil {
						log.Println("failed to reply image: ", err)
					}
				case replicate.PredictionFailure:
					if err := bot.ReplyText(token, "an error occurred"); err != nil {
						log.Println("failed to reply text: ", err)
					}
				}
			}
		}()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	},
}
