package samhook

import (
	"io"
	"log"
	"time"
)

// Logger 定義日誌記錄介面
type Logger interface {
	LogRequest(url string, method string, duration time.Duration, err error)
}

// defaultLogger 預設的日誌記錄器實現
type defaultLogger struct {
	logger *log.Logger
}

// LogRequest 記錄請求資訊
func (l *defaultLogger) LogRequest(url string, method string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	l.logger.Printf("[samhook] %s %s - %s - duration: %v", method, url, status, duration)
	if err != nil {
		l.logger.Printf("[samhook] error: %v", err)
	}
}

// noOpLogger 空操作日誌記錄器（不記錄任何內容）
type noOpLogger struct{}

func (l *noOpLogger) LogRequest(url string, method string, duration time.Duration, err error) {
	// 不執行任何操作
}

// 包級別的日誌記錄器（可選功能，用於簡單的日誌記錄）
var defaultPackageLogger Logger = &noOpLogger{}

// SetLogger 設置包級別的日誌記錄器
func SetLogger(logger Logger) {
	if logger == nil {
		defaultPackageLogger = &noOpLogger{}
	} else {
		defaultPackageLogger = logger
	}
}

// SetLoggerWriter 使用 io.Writer 設置包級別的日誌記錄器
func SetLoggerWriter(w io.Writer) {
	if w == nil {
		SetLogger(nil)
		return
	}
	SetLogger(&defaultLogger{
		logger: log.New(w, "", log.LstdFlags),
	})
}
