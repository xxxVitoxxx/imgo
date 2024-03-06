package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/xxxVitoxxx/imgo/pkg/rabbitmq"
)

type Client struct {
	url         string
	token       string
	version     string
	callbackURL string
	mq          *rabbitmq.MessageQueue
	queue       amqp091.Queue
}

// NewClient return a new client instance
func NewClient(url, token, version, callbackURL string, mq *rabbitmq.MessageQueue, queue amqp091.Queue) *Client {
	return &Client{url, token, version, callbackURL, mq, queue}
}

// Input depends on what model you are running
type Input map[string]interface{}

const defaultTime time.Duration = 20 * time.Second

type Replicate interface {
	CreatePrediction(input Input) (createReplicate, error)
	GetPrediction(predictionID string) (getReplicate, error)
}

func (c *Client) CreatePrediction(input Input) (createReplicate, error) {
	reqBody := Request{
		Version:      c.version,
		Input:        input,
		Webhook:      c.callbackURL,
		EventsFilter: []string{"completed"},
	}

	jsonByte, err := json.Marshal(reqBody)
	if err != nil {
		return createReplicate{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return createReplicate{}, err
	}

	req.Header.Add("Authorization", "Token "+c.token)
	req.Header.Add("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), defaultTime)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return createReplicate{}, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return createReplicate{}, err
	}

	var r createReplicate
	if err := json.Unmarshal(b, &r); err != nil {
		return createReplicate{}, err
	}
	return r, nil
}

func (c *Client) GetPrediction(predictionID string) (getReplicate, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s", c.url, predictionID),
		nil,
	)
	if err != nil {
		return getReplicate{}, err
	}

	req.Header.Add("Authorization", "Token "+c.token)
	req.Header.Add("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), defaultTime)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return getReplicate{}, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return getReplicate{}, err
	}

	var r getReplicate
	if err := json.Unmarshal(b, &r); err != nil {
		return getReplicate{}, err
	}
	return r, nil
}
