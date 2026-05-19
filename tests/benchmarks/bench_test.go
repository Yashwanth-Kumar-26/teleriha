package benchmarks

import (
	"bytes"
	"encoding/json"
	"net/url"
	"testing"
)

type apiResponse struct {
	Ok     bool   `json:"ok"`
	Result string `json:"result"`
}

var benchDecodeResponseBody []byte

func init() {
	resp := map[string]interface{}{
		"ok": true,
		"result": map[string]interface{}{
			"message_id": int64(42),
			"sender": map[string]interface{}{
				"id":    int64(12345),
				"is_bot": true,
			},
			"chat": map[string]interface{}{
				"id":   int64(12345),
				"type": "private",
			},
			"date": int64(1_700_000_000),
			"text": "hello",
		},
	}
	var err error
	benchDecodeResponseBody, err = json.Marshal(resp)
	if err != nil {
		panic(err)
	}
}

const benchDecodeResponseLabel = "benchDecodeResponseBody"

func BenchmarkDecodeResponse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = decodeResponse(benchDecodeResponseLabel, benchDecodeResponseBody)
	}
}

func decodeResponse(k string, body []byte) *apiResponse {
	_ = k
	var resp apiResponse
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&resp); err != nil {
		return nil
	}
	return &resp
}

var benchMarshalPayload = map[string]interface{}{
	"chat_id": int64(12345),
	"text":    "hello",
}

func BenchmarkMarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(benchMarshalPayload)
	}
}

var benchBaseURL = "https://api.telegram.org"

func BenchmarkBuildURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u, err := url.Parse(benchBaseURL)
		if err != nil {
			b.Fatal(err)
		}
		_ = u.JoinPath("bot", "TOKEN", "sendMessage")
	}
}
