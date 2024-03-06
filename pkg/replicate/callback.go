package replicate

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Client) Router(r *gin.Engine) {
	r.POST("/replicate/callback", c.CallBack)
}

type Status string

// prediction status
var (
	PredictionStart   Status = "starting"
	PredictionProcess Status = "processing"
	PredictionFailure Status = "failed"
	PredictionSuccess Status = "succeeded"
	PredictionCancel  Status = "canceled"
)

var PredStatusCh = make(chan PredictionStatus)

type PredictionStatus struct {
	ID     string
	Status Status
}

type Prediction struct {
	ID     string
	URL    string
	Status Status
}

func (c *Client) CallBack(ctx *gin.Context) {
	var resp getReplicate
	if err := ctx.BindJSON(&resp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	prediction := Prediction{resp.ID, resp.Output, resp.Status}
	b, err := json.Marshal(prediction)
	if err != nil {
		log.Println("failed to marshal prediction: %w", err)
	}

	if err := c.mq.Publish(context.Background(), c.queue, b); err != nil {
		log.Println("failed to publish prediction: %w", err)
	}

	ctx.Status(http.StatusOK)
}
