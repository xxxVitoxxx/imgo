package imgur_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xxxVitoxxx/imgo/pkg/imgur"
)

func TestUploadImage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header()["X-Post-Rate-Limit-Remaining"] = []string{"1000"}
			w.Write([]byte(`{
				"data": {
				  "id": "7yClKye",
				  "title": null,
				  "description": null,
				  "datetime": 1677436325,
				  "type": "image/jpeg",
				  "animated": false,
				  "width": 500,
				  "height": 491,
				  "size": 376922,
				  "views": 0,
				  "bandwidth": 0,
				  "vote": null,
				  "favorite": false,
				  "nsfw": null,
				  "section": null,
				  "account_url": null,
				  "account_id": 0,
				  "is_ad": false,
				  "in_most_viral": false,
				  "has_sound": false,
				  "tags": [],
				  "ad_type": 0,
				  "ad_url": "",
				  "edited": "0",
				  "in_gallery": false,
				  "deletehash": "wUMtREfMcTCLMCe",
				  "name": "",
				  "link": "https://i.imgur.com/7yClKye.jpg"
				},
				"success": true,
				"status": 200
			  }`))
		}))
		defer ts.Close()

		imgur := imgur.NewImgur(ts.URL, "")
		link, err := imgur.UploadImage([]byte(`test`))
		assert.NoError(t, err)
		assert.Equal(t, "https://i.imgur.com/7yClKye.jpg", link)
	})

	t.Run("fail", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header()["X-Post-Rate-Limit-Remaining"] = []string{"1000"}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"data": {
					"error": {
						"code": 1003,
						"message": "File type invalid (1)",
						"type": "ImgurException",
						"exception": []
					},
					"request": "\/3\/image",
					"method": "POST"
				},
				"success": false,
				"status": 400
			}`))
		}))
		defer ts.Close()

		imgur := imgur.NewImgur(ts.URL, "")
		_, err := imgur.UploadImage(nil)
		assert.EqualError(t, err, "status code: 400")
	})
}
