package samhook

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
)

// AddAttachment 添加一個attachment
func (m *Message) AddAttachment(attachment Attachment) *Message {
	m.Attachments = append(m.Attachments, attachment)
	return m
}

// AddAttachments 添加多個attachment
func (m *Message) AddAttachments(attachments []Attachment) *Message {
	m.Attachments = append(m.Attachments, attachments...)
	return m
}

// sendRequest 內部函數，統一處理 HTTP 請求
func sendRequest(client *http.Client, req *http.Request) error {
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	// 記錄請求日誌（如果設置了日誌記錄器）
	if defaultPackageLogger != nil {
		defaultPackageLogger.LogRequest(req.URL.String(), req.Method, duration, err)
	}

	if err != nil {
		return NewNetworkError(req.URL.String(), err)
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		// 讀取錯誤回應體
		bodyBytes, _ := io.ReadAll(resp.Body)
		responseBody := string(bodyBytes)
		apiErr := NewAPIError(req.URL.String(), resp.StatusCode, responseBody)

		// 記錄 API 錯誤日誌
		if defaultPackageLogger != nil {
			defaultPackageLogger.LogRequest(req.URL.String(), req.Method, duration, apiErr)
		}

		return apiErr
	}

	return nil
}

// Send 發送message
func Send(url string, msg Message) error {
	payloadBytes, err := sonic.Marshal(msg)
	if err != nil {
		return NewSerializationError(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return NewNetworkError(url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	return sendRequest(http.DefaultClient, req)
}

// SendReader 發送message
func SendReader(url string, r io.Reader) error {
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return NewNetworkError(url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	return sendRequest(http.DefaultClient, req)
}
