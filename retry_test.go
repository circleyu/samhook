package samhook

import (
	"net/http"
	"testing"
	"time"
)

func TestSendWithRetry_Success(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	msg := createTestMessage()
	opts := DefaultRetryOptions
	opts.MaxRetries = 1

	err := SendWithRetry(server.URL, msg, opts)
	if err != nil {
		t.Fatalf("SendWithRetry() error = %v", err)
	}
}

func TestSendWithRetry_RetryOnNetworkError(t *testing.T) {
	attempts := 0
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			// 模擬暫時性網路錯誤（通過返回 500）
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
	})

	msg := createTestMessage()
	opts := DefaultRetryOptions
	opts.MaxRetries = 3
	opts.Interval = 10 * time.Millisecond

	err := SendWithRetry(server.URL, msg, opts)
	if err != nil {
		t.Fatalf("SendWithRetry() error = %v", err)
	}
	if attempts < 2 {
		t.Error("expected retry to occur")
	}
}

func TestSendWithRetry_NoRetryOnClientError(t *testing.T) {
	attempts := 0
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	})

	msg := createTestMessage()
	opts := DefaultRetryOptions
	opts.MaxRetries = 3
	opts.Interval = 10 * time.Millisecond

	err := SendWithRetry(server.URL, msg, opts)
	if err == nil {
		t.Fatal("expected error")
	}

	// 4xx 錯誤不應該重試
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestExponentialBackoff_NextInterval(t *testing.T) {
	backoff := &ExponentialBackoff{
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		Jitter:          false,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
		{10, 30 * time.Second}, // 應該被限制在 MaxInterval
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.attempt)), func(t *testing.T) {
			interval := backoff.NextInterval(tt.attempt)
			if interval != tt.expected {
				t.Errorf("attempt %d: expected %v, got %v", tt.attempt, tt.expected, interval)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      *WebhookError
		expected bool
	}{
		{
			name:     "網路錯誤可重試",
			err:      NewNetworkError("https://example.com", http.ErrServerClosed),
			expected: true,
		},
		{
			name:     "5xx 錯誤可重試",
			err:      NewAPIError("https://example.com", 500, "server error"),
			expected: true,
		},
		{
			name:     "429 錯誤可重試",
			err:      NewAPIError("https://example.com", 429, "rate limited"),
			expected: true,
		},
		{
			name:     "4xx 錯誤不可重試",
			err:      NewAPIError("https://example.com", 400, "bad request"),
			expected: false,
		},
		{
			name:     "401 錯誤不可重試",
			err:      NewAPIError("https://example.com", 401, "unauthorized"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
