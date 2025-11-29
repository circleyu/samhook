package samhook

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// 錯誤類型常數
const (
	ErrorTypeNetwork       = "network"
	ErrorTypeSerialization = "serialization"
	ErrorTypeAPI           = "api"
	ErrorTypeUnknown       = "unknown"
)

// 錯誤代碼常數
const (
	ErrorCodeNetworkTimeout    = "NETWORK_TIMEOUT"
	ErrorCodeNetworkConnection = "NETWORK_CONNECTION"
	ErrorCodeNetworkDNS        = "NETWORK_DNS"
	ErrorCodeSerializationJSON = "SERIALIZATION_JSON"
	ErrorCodeAPIUnauthorized   = "API_UNAUTHORIZED"
	ErrorCodeAPIForbidden      = "API_FORBIDDEN"
	ErrorCodeAPINotFound       = "API_NOT_FOUND"
	ErrorCodeAPIRateLimit      = "API_RATE_LIMIT"
	ErrorCodeAPIServerError    = "API_SERVER_ERROR"
)

// WebhookError 表示 webhook 操作中的錯誤
type WebhookError struct {
	// Type 錯誤類型
	Type string

	// StatusCode HTTP 狀態碼（如果是 API 錯誤）
	StatusCode int

	// Message 錯誤訊息
	Message string

	// ResponseBody API 回應體（如果是 API 錯誤）
	ResponseBody string

	// Err 原始錯誤
	Err error

	// URL webhook URL（用於上下文）
	URL string
}

// Error 實現 error 介面，提供詳細的錯誤訊息
func (e *WebhookError) Error() string {
	var parts []string

	// 錯誤類型
	parts = append(parts, fmt.Sprintf("[%s]", e.Type))

	// HTTP 狀態碼（如果適用）
	if e.StatusCode > 0 {
		parts = append(parts, fmt.Sprintf("HTTP %d", e.StatusCode))
	}

	// 主要錯誤訊息
	parts = append(parts, e.Message)

	// URL 上下文（如果有）
	if e.URL != "" {
		parts = append(parts, fmt.Sprintf("(URL: %s)", e.URL))
	}

	// 原始錯誤（如果有）
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Err))
	}

	return strings.Join(parts, " ")
}

// Unwrap 返回原始錯誤（支援 errors.Unwrap）
func (e *WebhookError) Unwrap() error {
	return e.Err
}

// IsNetworkError 判斷是否為網路錯誤
func (e *WebhookError) IsNetworkError() bool {
	return e.Type == ErrorTypeNetwork
}

// IsSerializationError 判斷是否為序列化錯誤
func (e *WebhookError) IsSerializationError() bool {
	return e.Type == ErrorTypeSerialization
}

// IsAPIError 判斷是否為 API 錯誤
func (e *WebhookError) IsAPIError() bool {
	return e.Type == ErrorTypeAPI
}

// GetStatusCode 返回 HTTP 狀態碼（如果是 API 錯誤）
func (e *WebhookError) GetStatusCode() int {
	return e.StatusCode
}

// GetResponseBody 返回 API 回應體（如果是 API 錯誤）
func (e *WebhookError) GetResponseBody() string {
	return e.ResponseBody
}

// GetErrorCode 返回錯誤代碼
func (e *WebhookError) GetErrorCode() string {
	// 根據類型和狀態碼返回具體錯誤代碼
	if e.IsNetworkError() {
		// 可以根據原始錯誤進一步分類
		return ErrorCodeNetworkConnection
	}
	if e.IsAPIError() {
		switch e.StatusCode {
		case 401:
			return ErrorCodeAPIUnauthorized
		case 403:
			return ErrorCodeAPIForbidden
		case 404:
			return ErrorCodeAPINotFound
		case 429:
			return ErrorCodeAPIRateLimit
		case 500, 502, 503, 504:
			return ErrorCodeAPIServerError
		}
	}
	if e.IsSerializationError() {
		return ErrorCodeSerializationJSON
	}
	return ErrorTypeUnknown
}

// DetailedMessage 返回詳細的錯誤訊息（多行格式）
func (e *WebhookError) DetailedMessage() string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("Webhook Error [%s]\n", e.Type))

	if e.StatusCode > 0 {
		buf.WriteString(fmt.Sprintf("  Status Code: %d\n", e.StatusCode))
	}

	buf.WriteString(fmt.Sprintf("  Message: %s\n", e.Message))

	if e.URL != "" {
		buf.WriteString(fmt.Sprintf("  URL: %s\n", e.URL))
	}

	if e.ResponseBody != "" {
		buf.WriteString(fmt.Sprintf("  Response: %s\n", e.ResponseBody))
	}

	if e.Err != nil {
		buf.WriteString(fmt.Sprintf("  Cause: %v\n", e.Err))
	}

	return buf.String()
}

// NewNetworkError 創建網路錯誤
func NewNetworkError(url string, err error) *WebhookError {
	return &WebhookError{
		Type:    ErrorTypeNetwork,
		Message: fmt.Sprintf("network error: %v", err),
		Err:     err,
		URL:     url,
	}
}

// NewSerializationError 創建序列化錯誤
func NewSerializationError(err error) *WebhookError {
	return &WebhookError{
		Type:    ErrorTypeSerialization,
		Message: fmt.Sprintf("serialization error: %v", err),
		Err:     err,
	}
}

// NewAPIError 創建 API 錯誤
func NewAPIError(url string, statusCode int, responseBody string) *WebhookError {
	message := fmt.Sprintf("API returned status %d", statusCode)
	if responseBody != "" {
		message = fmt.Sprintf("%s: %s", message, responseBody)
	}

	return &WebhookError{
		Type:         ErrorTypeAPI,
		StatusCode:   statusCode,
		Message:      message,
		ResponseBody: responseBody,
		URL:          url,
	}
}

// classifyError 分類標準錯誤為 WebhookError
func classifyError(webhookURL string, err error) *WebhookError {
	if err == nil {
		return nil
	}

	// 檢查是否已經是 WebhookError
	if webhookErr, ok := err.(*WebhookError); ok {
		return webhookErr
	}

	// 嘗試分類錯誤類型
	if urlErr, ok := err.(*url.Error); ok {
		if _, ok := urlErr.Err.(*net.OpError); ok {
			return NewNetworkError(webhookURL, err)
		}
	}

	// 預設為未知錯誤
	return &WebhookError{
		Type:    ErrorTypeUnknown,
		Message: err.Error(),
		Err:     err,
		URL:     webhookURL,
	}
}
