package line

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/xxxVitoxxx/imgo/pkg/bucket"
	"github.com/xxxVitoxxx/imgo/pkg/imgur"
	"github.com/xxxVitoxxx/imgo/pkg/replicate"
	"github.com/xxxVitoxxx/imgo/pkg/storage"
)

type lineBot struct {
	bot *linebot.Client
	imgur.Imgur
	replicate.Replicate
	bucket *bucket.Bucket
	redis  *storage.Redis
}

// NewLineBot return a new lineBot instance
func NewLineBot(
	secret string,
	token string,
	img imgur.Imgur,
	rep replicate.Replicate,
	bucket *bucket.Bucket,
	redis *storage.Redis,
) (*lineBot, error) {
	bot, err := linebot.New(secret, token)
	if err != nil {
		return &lineBot{}, err
	}
	return &lineBot{bot, img, rep, bucket, redis}, err
}

func (b *lineBot) ReplyText(replyToken, text string) error {
	replyCall := b.bot.ReplyMessage(replyToken, linebot.NewTextMessage(text))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	replyCall = replyCall.WithContext(ctx)
	if _, err := replyCall.Do(); err != nil {
		return fmt.Errorf("reply text: %w", err)
	}
	return nil
}

func (b *lineBot) ReplyImage(replyToken, imageURL string) error {
	if imageURL == "" {
		return errors.New("image url is empty")
	}

	replyCall := b.bot.ReplyMessage(replyToken, linebot.NewImageMessage(imageURL, imageURL))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	replyCall = replyCall.WithContext(ctx)
	if _, err := replyCall.Do(); err != nil {
		return fmt.Errorf("reply image: %w", err)
	}
	return nil
}
