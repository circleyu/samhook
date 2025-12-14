package samhook

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidateWebhookURL 驗證 webhook URL 格式和協議
func ValidateWebhookURL(webhookURL string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// 檢查協議
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https, got: %s", parsedURL.Scheme)
	}

	// 檢查主機名
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}
