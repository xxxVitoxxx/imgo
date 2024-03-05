package replicate

import (
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

func (c *Client) CallBack(ctx *gin.Context) {
	var resp getReplicate
	if err := ctx.BindJSON(&resp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	PredStatusCh <- PredictionStatus{resp.ID, resp.Status}
	ctx.Status(http.StatusOK)
}
