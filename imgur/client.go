package imgur

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type imgur struct {
	url string
	id  string
}

// NewImgur return a new imgur instance
func NewImgur(url, id string) *imgur {
	return &imgur{url, id}
}

const defaultTime time.Duration = 20 * time.Second

type Imgur interface {
	UploadImage(imageByte []byte) (string, error)
}

// respImgur imgur APIs response
type respImgur struct {
	Data    data `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

type data struct {
	Err  err    `json:"error"`
	Link string `json:"link"`
}

type err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Method  string `json:"method"`
}

// UploadImage will return image URL of imgur
func (img *imgur) UploadImage(imageByte []byte) (string, error) {
	values := url.Values{}
	values.Set("image", string(imageByte))
	enc := values.Encode()
	req, err := http.NewRequest(
		http.MethodPost,
		img.url,
		bytes.NewBuffer([]byte(enc)),
	)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Add("Authorization", "Client-ID "+img.id)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	ctx, cancel := context.WithTimeout(context.Background(), defaultTime)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", fmt.Errorf("client do: %w", err)
	}

	respHeader := resp.Header
	// this limit is 1,250 POST requests per hour.
	if respHeader["X-Post-Rate-Limit-Remaining"][0] == "0" {
		return "", errors.New("imgur APIs post request hit a limit")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r respImgur
	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}

	if !r.Success || r.Status != http.StatusOK {
		return "", fmt.Errorf("error code: %d\nmessage: %s", r.Data.Err.Code, r.Data.Err.Message)
	}

	return r.Data.Link, nil
}
