package samhook

import (
	"math"
	"math/rand"
	"time"
)

// RetryOptions 重試選項
type RetryOptions struct {
	MaxRetries int
	Interval   time.Duration
	Backoff    *ExponentialBackoff
}

// ExponentialBackoff 指數退避
type ExponentialBackoff struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          bool
}

// DefaultRetryOptions 預設重試選項
var DefaultRetryOptions = RetryOptions{
	MaxRetries: 3,
	Interval:   1 * time.Second,
	Backoff: &ExponentialBackoff{
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		Jitter:          true,
	},
}

// NextInterval 計算下一次重試間隔
func (eb *ExponentialBackoff) NextInterval(attempt int) time.Duration {
	interval := float64(eb.InitialInterval) * math.Pow(eb.Multiplier, float64(attempt))
	if interval > float64(eb.MaxInterval) {
		interval = float64(eb.MaxInterval)
	}

	if eb.Jitter {
		// 添加隨機抖動（±10%）
		jitter := interval * 0.1 * (rand.Float64()*2 - 1)
		interval += jitter
	}

	return time.Duration(interval)
}

// SendWithRetry 帶重試的發送，支援自訂客戶端配置
func SendWithRetry(url string, msg Message, opts RetryOptions, clientOpts ...ClientOption) error {
	var lastErr error
	interval := opts.Interval

	for i := 0; i <= opts.MaxRetries; i++ {
		err := SendWithOptions(url, msg, clientOpts...)
		if err == nil {
			return nil
		}

		// 檢查是否可重試
		if webhookErr, ok := err.(*WebhookError); ok {
			if !isRetryable(webhookErr) {
				return err
			}
		} else {
			// 非 WebhookError 預設不重試
			return err
		}

		lastErr = err
		if i < opts.MaxRetries {
			// 計算重試間隔
			if opts.Backoff != nil {
				interval = opts.Backoff.NextInterval(i)
			}
			time.Sleep(interval)
		}
	}
	return lastErr
}

// isRetryable 判斷錯誤是否可重試
func isRetryable(err *WebhookError) bool {
	// 網路錯誤可以重試
	if err.IsNetworkError() {
		return true
	}
	// 5xx 錯誤可以重試
	if err.IsAPIError() && err.StatusCode >= 500 {
		return true
	}
	// 429 速率限制可以重試
	if err.IsAPIError() && err.StatusCode == 429 {
		return true
	}
	// 4xx 錯誤不重試
	return false
}
