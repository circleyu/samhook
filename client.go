package samhook

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// ClientOption 客戶端選項
type ClientOption func(*http.Client)

// WithClient 使用自訂 HTTP 客戶端
func WithClient(client *http.Client) ClientOption {
	return func(c *http.Client) {
		if client != nil {
			*c = *client
		}
	}
}

// WithTimeout 設置超時
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}

// DefaultTimeout 預設超時時間
const DefaultTimeout = 10 * time.Second

// SendWithOptions 使用選項發送訊息
func SendWithOptions(url string, msg Message, opts ...ClientOption) error {
	client := &http.Client{
		Timeout: DefaultTimeout,
	}
	// 應用選項
	for _, opt := range opts {
		opt(client)
	}

	payloadBytes, err := json.Marshal(msg)
	if err != nil {
		return NewSerializationError(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return NewNetworkError(url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return NewNetworkError(url, err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		responseBody := string(bodyBytes)
		return NewAPIError(url, resp.StatusCode, responseBody)
	}

	return nil
}

// SendWithContext 使用 Context 發送訊息
func SendWithContext(ctx context.Context, url string, msg Message, opts ...ClientOption) error {
	client := &http.Client{
		Timeout: DefaultTimeout,
	}
	// 應用選項
	for _, opt := range opts {
		opt(client)
	}

	payloadBytes, err := json.Marshal(msg)
	if err != nil {
		return NewSerializationError(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return NewNetworkError(url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return NewNetworkError(url, err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		responseBody := string(bodyBytes)
		return NewAPIError(url, resp.StatusCode, responseBody)
	}

	return nil
}
