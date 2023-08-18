package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
)

type HTTPClient struct {
	path string
}

func NewHTTPClient(path string) *HTTPClient {
	return &HTTPClient{path}
}

func (h *HTTPClient) Request(method string, body any, headers map[string]string) (respBody []byte, err error) {
	jsonBody, err := comjson.Marshal(body)
	if err != nil {
		return
	}
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(method, h.path, bodyReader)
	if err != nil {
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	return
}
