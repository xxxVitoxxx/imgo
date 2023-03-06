# flow

先用 create a prediction 

request

```curl
curl -s -X POST \
-d '{"version": "7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56", "input": {"text": "vito"}}' \
-H "Authorization: Token 787c600f06e8f674725b45dee573067a27c04d4d" \
-H 'Content-Type: application/json' \
https://api.replicate.com/v1/predictions
```

response

```json
{
    "completed_at": null,
    "created_at": "2023-02-19T21:11:21.781455Z",
    "error": null,
    "id": "h2ochsbznveatar4f4wp2mlwmy",
    "input": {},
    "logs": "",
    "metrics": {},
    "output": null,
    "started_at": null,
    "status": "starting",
    "urls": {
        "get": "https://api.replicate.com/v1/predictions/h2ochsbznveatar4f4wp2mlwmy",
        "cancel": "https://api.replicate.com/v1/predictions/h2ochsbznveatar4f4wp2mlwmy/cancel"
    },
    "version": "7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56",
    "webhook_completed": null
}
```

之後可以透過 prediction id 取得 prediction 資訊

`GET https://api.replicate.com/v1/predictions/{prediction_id}`

request

```curl
curl -s GET \
-H "Authorization: Token 787c600f06e8f674725b45dee573067a27c04d4d" \
https://api.replicate.com/v1/predictions/h2ochsbznveatar4f4wp2mlwmy
```

response

```json
{
  "completed_at": "2023-02-19T21:11:21.892710Z",
  "created_at": "2023-02-19T21:11:21.781455Z",
  "error": "Prediction input failed validation: {\"detail\":[{\"loc\":[\"body\",\"input\",\"image\"],\"msg\":\"field required\",\"type\":\"value_error.missing\"}]}",
  "id": "h2ochsbznveatar4f4wp2mlwmy",
  "input": {},
  "logs": "",
  "metrics": {
    "predict_time": 0.000004
  },
  "output": null,
  "started_at": "2023-02-19T21:11:21.892706Z",
  "status": "failed",
  "urls": {
    "get": "https://api.replicate.com/v1/predictions/h2ochsbznveatar4f4wp2mlwmy",
    "cancel": "https://api.replicate.com/v1/predictions/h2ochsbznveatar4f4wp2mlwmy/cancel"
  },
  "version": "7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56",
  "webhook_completed": null
}
```

line 可以設定 replicate model 參數，會將參數跟 line ID 記錄一起，當 line user send image 會去取得參數並連同照片的 URL 發送到 replicate 

issue 同時送多張有時候 imgURLCh 會是空的，但實際上所有照片都有成功轉換

```
replicate callback:  {eg7ozdqtnbfevl2y5unneqcbly 7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56 map[image:https://i.imgur.com/KUH0Fta.jpg] https://replicate.delivery/pbxt/heCSeHfVziLbkJSIhZXV5ViMIsd1Qaj8vrGEgPZ99hL8MUFhA/output.png detect 0 faces <nil> succeeded 2023-02-27 17:59:53.141056 +0000 UTC 2023-02-27 17:59:53.217091 +0000 UTC 2023-02-27 17:59:59.559035 +0000 UTC https://7042-220-135-89-54.jp.ngrok.io/webhook {https://api.replicate.com/v1/predictions/eg7ozdqtnbfevl2y5unneqcbly https://api.replicate.com/v1/predictions/eg7ozdqtnbfevl2y5unneqcbly/cancel} map[predict_time:6.341944]}
[GIN] 2023/02/28 - 01:59:59 | 200 |     195.292µs |  35.184.121.158 | POST     "/webhook"
imgURLCh:  https://replicate.delivery/pbxt/heCSeHfVziLbkJSIhZXV5ViMIsd1Qaj8vrGEgPZ99hL8MUFhA/output.png
replicate callback:  {gqllr4kg7beudedci5pm6s6fy4 7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56 map[image:https://i.imgur.com/ntJRoJ2.jpg] https://replicate.delivery/pbxt/RBVAvWio0N6RLp8YaCbYjWDSfeQReGf9GgLa9zo8tN6LaoKCB/output.png detect 1 faces <nil> succeeded 2023-02-27 17:59:54.368251 +0000 UTC 2023-02-27 17:59:54.441952 +0000 UTC 2023-02-27 18:00:03.518778 +0000 UTC https://7042-220-135-89-54.jp.ngrok.io/webhook {https://api.replicate.com/v1/predictions/gqllr4kg7beudedci5pm6s6fy4 https://api.replicate.com/v1/predictions/gqllr4kg7beudedci5pm6s6fy4/cancel} map[predict_time:9.076826]}
[GIN] 2023/02/28 - 02:00:04 | 200 |     208.084µs |  35.226.223.228 | POST     "/webhook"
imgURLCh:  
2023/02/28 02:00:04 reply image: linebot: APIError 400 The request body has 2 error(s)
[messages[0].originalContentUrl] May not be empty
[messages[0].previewImageUrl] May not be empty
replicate callback:  {twlyaabvxje4rc6kyludasdyca 7de2ea26c616d5bf2245ad0d5e24f0ff9a6204578a5c876db53142edd9d2cd56 map[image:https://i.imgur.com/GwsJsdR.jpg] https://replicate.delivery/pbxt/Khjl1K0zXlKVA59B5Ea6YJQdTDflgW0wxzYeNxJCJH7kGqiQA/output.png detect 3 faces <nil> succeeded 2023-02-27 17:59:53.809177 +0000 UTC 2023-02-27 17:59:53.889223 +0000 UTC 2023-02-27 18:00:05.053406 +0000 UTC https://7042-220-135-89-54.jp.ngrok.io/webhook {https://api.replicate.com/v1/predictions/twlyaabvxje4rc6kyludasdyca https://api.replicate.com/v1/predictions/twlyaabvxje4rc6kyludasdyca/cancel} map[predict_time:11.164183]}
[GIN] 2023/02/28 - 02:00:05 | 200 |     230.291µs |  35.184.121.158 | POST     "/webhook"
imgURLCh:  https://replicate.delivery/pbxt/RBVAvWio0N6RLp8YaCbYjWDSfeQReGf9GgLa9zo8tN6LaoKCB/output.png
```

失敗原因：
假設傳三張照片（三個 goroutine），同時 create prediction ，當 replica†e callback 回傳 status ，此時三個 goroutine 都有可能接收到該 status ，
若剛好接收到 status 的 goroutine 還在 processing ，這樣 get prediction 就會取得 processing 的結果，取得的 image URL 會是空白，當 line reply image 
接收到空白的 image URL 就會出錯

解法：
三個都丟 goroutine ，另外開一個 goroutine 接收 status 及 prediction ID ，當接收到 success 就拿 prediction ID 去 get prediction 

replicate rate limit not do it

