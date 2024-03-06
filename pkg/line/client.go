package line

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/xxxVitoxxx/imgo/pkg/imgur"
	"github.com/xxxVitoxxx/imgo/pkg/replicate"
)

type lineBot struct {
	bot *linebot.Client
	imgur.Imgur
	replicate.Replicate
}

// NewLineBot return a new lineBot instance
func NewLineBot(secret, token string, img imgur.Imgur, rep replicate.Replicate) (*lineBot, error) {
	bot, err := linebot.New(secret, token)
	if err != nil {
		return &lineBot{}, err
	}
	return &lineBot{bot, img, rep}, err
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
