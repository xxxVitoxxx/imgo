package replicate_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xxxVitoxxx/imgo/replicate"
)

func TestCreatePrediction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"completed_at": null,
			"created_at": "2023-02-25T16:39:31.097550Z",
			"error": null,
			"id": "76d7fp6tifd6pprabbw6gtjnzy",
			"input": {
				"image": "https://cdn.hk01.com/di/media/images/dw/20200805/367761876806930432.jpeg/LypEctZ0sCisiMo4oUbMTtpO8bn7a6xHt7nyrre58q4?v=w1920"
			},
			"logs": "",
			"metrics": {},
			"output": null,
			"started_at": null,
			"status": "starting",
			"urls": {
				"get": "https://api.replicate.com/v1/predictions/76d7fp6tifd6pprabbw6gtjnzy",
				"cancel": "https://api.replicate.com/v1/predictions/76d7fp6tifd6pprabbw6gtjnzy/cancel"
			},
			"version": "7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56",
			"webhook_completed": null
		}`))
	}))
	defer ts.Close()

	c := replicate.NewClient(ts.URL, "", "", "")
	input := replicate.Input{
		"image": "https://cdn.hk01.com/di/media/images/dw/20200805/367761876806930432.jpeg/LypEctZ0sCisiMo4oUbMTtpO8bn7a6xHt7nyrre58q4?v=w1920",
	}
	_, err := c.CreatePrediction(input)
	assert.NoError(t, err)
}

func TestGetPrediction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"completed_at": "2023-02-25T16:40:01.312425Z",
			"created_at": "2023-02-25T16:39:31.097550Z",
			"error": null,
			"id": "76d7fp6tifd6pprabbw6gtjnzy",
			"input": {
				"image": "https://cdn.hk01.com/di/media/images/dw/20200805/367761876806930432.jpeg/LypEctZ0sCisiMo4oUbMTtpO8bn7a6xHt7nyrre58q4?v=w1920"
			},
			"logs": "detect 4 faces",
			"metrics": {
				"predict_time": 15.779969
			},
			"output": "https://replicate.delivery/pbxt/blTMMYRfJz1TIKfRs18OmTYP0f1QlnGloEk9KTT4hSZBf6HCB/output.png",
			"started_at": "2023-02-25T16:39:45.532456Z",
			"status": "succeeded",
			"urls": {
				"get": "https://api.replicate.com/v1/predictions/76d7fp6tifd6pprabbw6gtjnzy",
				"cancel": "https://api.replicate.com/v1/predictions/76d7fp6tifd6pprabbw6gtjnzy/cancel"
			},
			"version": "7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56",
			"webhook_completed": null
		}`))
	}))
	defer ts.Close()

	c := replicate.NewClient(ts.URL, "", "", "")
	_, err := c.GetPrediction("76d7fp6tifd6pprabbw6gtjnzy")
	assert.NoError(t, err)
}
