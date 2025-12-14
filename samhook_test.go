package samhook

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bytedance/sonic"
)

// mockWebhookServer 創建一個 mock webhook 伺服器
func mockWebhookServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})
	return server
}

// createTestMessage 創建測試用的 Message
func createTestMessage() Message {
	return Message{
		Text:     "Test message",
		Username: "test-bot",
	}
}

// createTestAttachment 創建測試用的 Attachment
func createTestAttachment() Attachment {
	return Attachment{
		Color: Good,
		Title: "Test Attachment",
		Text:  "Test attachment text",
	}
}

func TestSend_Success(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	msg := createTestMessage()
	err := Send(server.URL, msg)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestSend_RequestFormat(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		// 驗證 JSON 格式
		var data map[string]interface{}
		if err := sonic.Unmarshal(body, &data); err != nil {
			t.Errorf("invalid JSON: %v", err)
		}

		// 驗證欄位
		if data["text"] != "Test" {
			t.Errorf("text field mismatch")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	msg := Message{Text: "Test"}
	if err := Send(server.URL, msg); err != nil {
		t.Fatal(err)
	}
}

func TestSend_ErrorScenarios(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		expectedError bool
		expectedType  string
	}{
		{
			name: "HTTP 400 Bad Request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid_payload"))
			},
			expectedError: true,
			expectedType:  ErrorTypeAPI,
		},
		{
			name: "HTTP 401 Unauthorized",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid_token"))
			},
			expectedError: true,
			expectedType:  ErrorTypeAPI,
		},
		{
			name: "HTTP 429 Rate Limit",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("rate_limited"))
			},
			expectedError: true,
			expectedType:  ErrorTypeAPI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockWebhookServer(t, tt.handler)

			msg := Message{Text: "Test"}
			err := Send(server.URL, msg)

			if tt.expectedError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// 驗證錯誤類型
			if err != nil {
				if webhookErr, ok := err.(*WebhookError); ok {
					if webhookErr.Type != tt.expectedType {
						t.Errorf("expected error type %s, got %s", tt.expectedType, webhookErr.Type)
					}
				}
			}
		})
	}
}

func TestSendReader_Success(t *testing.T) {
	server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	jsonData := `{"text":"Test","username":"bot"}`
	reader := bytes.NewReader([]byte(jsonData))

	err := SendReader(server.URL, reader)
	if err != nil {
		t.Fatalf("SendReader() error = %v", err)
	}
}

func TestSendReader_VariousReaders(t *testing.T) {
	tests := []struct {
		name   string
		reader io.Reader
	}{
		{"bytes.Reader", bytes.NewReader([]byte(`{"text":"test"}`))},
		{"strings.Reader", strings.NewReader(`{"text":"test"}`)},
		{"bytes.Buffer", bytes.NewBuffer([]byte(`{"text":"test"}`))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			if err := SendReader(server.URL, tt.reader); err != nil {
				t.Errorf("SendReader() error = %v", err)
			}
		})
	}
}

func TestAddAttachment(t *testing.T) {
	msg := Message{Text: "Test"}
	attachment := createTestAttachment()

	msg.AddAttachment(attachment)

	if len(msg.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(msg.Attachments))
	}
	if msg.Attachments[0].Title != attachment.Title {
		t.Errorf("attachment title mismatch")
	}
}

func TestAddAttachments(t *testing.T) {
	msg := Message{Text: "Test"}
	attachments := []Attachment{
		createTestAttachment(),
		createTestAttachment(),
	}

	msg.AddAttachments(attachments)

	if len(msg.Attachments) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(msg.Attachments))
	}
}

func TestSend_NetworkError(t *testing.T) {
	// 使用無效 URL 模擬網路錯誤
	invalidURL := "http://invalid-domain-that-does-not-exist.local"
	msg := Message{Text: "Test"}

	err := Send(invalidURL, msg)
	if err == nil {
		t.Fatal("expected network error")
	}

	// 驗證錯誤類型
	if webhookErr, ok := err.(*WebhookError); ok {
		if !webhookErr.IsNetworkError() {
			t.Error("expected network error type")
		}
	}
}

func TestSend_HTTPErrorStatusCodes(t *testing.T) {
	statusCodes := []int{400, 401, 403, 404, 429, 500, 502, 503}

	for _, statusCode := range statusCodes {
		t.Run(string(rune(statusCode)), func(t *testing.T) {
			server := mockWebhookServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(statusCode)
				w.Write([]byte("error"))
			})

			msg := Message{Text: "Test"}
			err := Send(server.URL, msg)

			if err == nil {
				t.Error("expected error for non-200 status")
			}

			// 驗證是 API 錯誤
			if webhookErr, ok := err.(*WebhookError); ok {
				if !webhookErr.IsAPIError() {
					t.Error("expected API error type")
				}
				if webhookErr.GetStatusCode() != statusCode {
					t.Errorf("expected status %d, got %d", statusCode, webhookErr.GetStatusCode())
				}
			}
		})
	}
}

func TestWebhookError_Methods(t *testing.T) {
	err := NewAPIError("https://example.com", 401, "invalid_token")

	if !err.IsAPIError() {
		t.Error("expected API error")
	}
	if err.IsNetworkError() {
		t.Error("should not be network error")
	}
	if err.GetStatusCode() != 401 {
		t.Errorf("expected status 401, got %d", err.GetStatusCode())
	}
	if err.GetResponseBody() != "invalid_token" {
		t.Errorf("expected response body 'invalid_token', got %s", err.GetResponseBody())
	}
	if err.GetErrorCode() != ErrorCodeAPIUnauthorized {
		t.Errorf("expected error code %s, got %s", ErrorCodeAPIUnauthorized, err.GetErrorCode())
	}
}
