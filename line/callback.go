package line

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/xxxVitoxxx/imgo/replicate"
)

func (b *lineBot) Router(r *gin.Engine) {
	r.POST("/line/callback", b.Callback)
}

// token effective time
var effectiveTime time.Duration = 50 * time.Second

func (b *lineBot) Callback(c *gin.Context) {
	// reply token effective
	replyTokenTime := time.After(effectiveTime)
	events, err := b.bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(http.StatusBadRequest)
		} else {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch eventType := event.Message.(type) {
			case *linebot.ImageMessage:
				imgURL, err := b.imgByteToURL(eventType.ID)
				if err != nil {
					log.Println("image byte to URL err: ", err)
					if err := b.replyText(event.ReplyToken, "process fail"); err != nil {
						log.Println(err)
					}

					return
				}

				finishCh := make(chan struct{})
				errCh := make(chan error)
				go b.restorationPhoto(imgURL, event.ReplyToken, finishCh, errCh)

				go func() {
					for {
						select {
						case err := <-errCh:
							log.Println("restoration err: ", err)
							if err := b.replyText(event.ReplyToken, "process fail"); err != nil {
								log.Println(err)
							}
							return
						case <-finishCh:
							return
						case <-replyTokenTime:
							log.Printf("user: %s processing time exceeded\n", event.Source.UserID)
							if err := b.replyText(event.ReplyToken, "processing time exceeded"); err != nil {
								log.Println(err)
							}
							return
						}
					}
				}()
			}
		}
	}
}

// imgByteToURL convert image byte to image URL
func (b *lineBot) imgByteToURL(eventTypeID string) (string, error) {
	contentCall := b.bot.GetMessageContent(eventTypeID)
	msgResp, err := contentCall.Do()
	if err != nil {
		return "", fmt.Errorf("message content call do error: %w", err)
	}

	imgByte, err := io.ReadAll(msgResp.Content)
	if err != nil {
		return "", fmt.Errorf("read image byte error: %w", err)
	}

	imgURL, err := b.UploadImage(imgByte)
	if err != nil {
		return "", fmt.Errorf("upload image err: %w", err)
	}
	return imgURL, nil
}

// Robust face restoration algorithm for old photos
func (b *lineBot) restorationPhoto(imgURL, replyToken string, finishCh chan<- struct{}, errCh chan<- error) {
	input := replicate.Input{
		"image": imgURL,
	}

	prediction, err := b.CreatePrediction(input)
	if err != nil {
		log.Println("create prediction err: ", err)
		if err := b.replyText(replyToken, "process fail"); err != nil {
			log.Println(err)
		}
		errCh <- err
	}
LOOP:
	// waiting replicate prediction callback
	// if replicate ID not your ID will throw data back to channel
	// and waiting next PredStatusCh
	predStatusCh := <-replicate.PredStatusCh
	if predStatusCh.ID != prediction.ID {
		replicate.PredStatusCh <- predStatusCh
		goto LOOP
	}

	switch predStatusCh.Status {
	case replicate.PredictionStart:
		// TODO start message
	case replicate.PredictionProcess:
		// TODO process message
	case replicate.PredictionFailure:
		// TODO  failure message
	case replicate.PredictionSuccess:
		// TODO success message
		pd, err := b.GetPrediction(prediction.ID)
		if err != nil {
			log.Println("get prediction err: ", err)
			errCh <- err
		}

		if err := b.replyImage(replyToken, pd.Output); err != nil {
			log.Println(err)
			errCh <- err
		}
	case replicate.PredictionCancel:
		// TODO cancel message
	}
	finishCh <- struct{}{}
}
