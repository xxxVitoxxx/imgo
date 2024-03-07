package line

import (
	"context"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/xxxVitoxxx/imgo/pkg/replicate"
)

func (b *lineBot) Router(r *gin.Engine) {
	r.POST("/line/callback", b.Callback)
}

func (b *lineBot) Callback(c *gin.Context) {
	events, err := b.bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(http.StatusBadRequest)
		} else {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	defer c.Status(http.StatusOK)

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch eventType := event.Message.(type) {
			case *linebot.ImageMessage:
				contentCall := b.bot.GetMessageContent(eventType.ID)
				msgResp, err := contentCall.Do()
				if err != nil {
					log.Println(err)
					return
				}

				filename := randFilename(16)
				if err := b.bucket.PutImage(filename, msgResp.Content); err != nil {
					log.Println(err)
					return
				}

				imgURL, err := b.bucket.GetImageURL(filename)
				if err != nil {
					log.Println(err)
					return
				}

				prediction, err := b.CreatePrediction(replicate.Input{"imgae": imgURL})
				if err != nil {
					log.Println(err)
					return
				}

				if err := b.redis.SetData(context.Background(), prediction.ID, event.ReplyToken); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func randFilename(length int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	filename := make([]rune, length)
	for i := range filename {
		filename[i] = letter[rand.Intn(len(letter))]
	}
	return string(filename)
}
