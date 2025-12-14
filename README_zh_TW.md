# samhook

[English](README.md) | 繁體中文

一個輕量級的 Go 語言函式庫，用於發送 Slack 和 Mattermost webhook 訊息。

## 功能特性

- ✅ **輕量級**: 極少外部依賴，使用高性能 JSON 庫
- ✅ **類型安全**: 完整的 Go 類型定義
- ✅ **鏈式調用**: 支援方法鏈式調用以提升開發體驗
- ✅ **靈活輸入**: 支援結構體和 Reader 兩種輸入方式
- ✅ **錯誤處理**: 詳細的錯誤類型和分類
- ✅ **可配置**: 支援自訂 HTTP 客戶端和超時
- ✅ **重試機制**: 可選的重試功能，支援指數退避

## 安裝

```bash
go get github.com/circleyu/samhook
```

## 快速開始

### 基本使用

```go
package main

import (
    "log"
    "github.com/circleyu/samhook"
)

func main() {
    webhookURL := "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
    
    msg := samhook.Message{
        Text:     "Hello from samhook!",
        Username: "samhook-bot",
    }
    
    err := samhook.Send(webhookURL, msg)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 使用附件

```go
msg := samhook.Message{
    Text: "系統通知",
}

attachment := samhook.Attachment{
    Color: samhook.Good,
    Title: "操作成功",
    Text:  "所有任務已完成",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### 錯誤處理

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        if webhookErr.IsNetworkError() {
            // 處理網路錯誤，可以重試
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            if statusCode == 429 {
                // 處理速率限制
            }
        }
    }
}
```

### 使用自訂客戶端

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

// 使用自訂超時
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)

// 使用 Context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := samhook.SendWithContext(ctx, webhookURL, msg)
```

### 使用重試機制

```go
opts := samhook.DefaultRetryOptions
opts.MaxRetries = 5

err := samhook.SendWithRetry(webhookURL, msg, opts)
```

## 文檔

- [API 文檔](docs/api_zh_TW.md) - 完整的 API 參考
- [架構文檔](docs/architecture_zh_TW.md) - 專案架構與設計
- [使用範例](docs/examples_zh_TW.md) - 詳細的使用範例

## 專案結構

```
samhook/
├── go.mod          # Go 模組定義
├── message.go      # 訊息資料結構定義
├── samhook.go       # 核心發送功能
├── error.go        # 錯誤類型定義
├── client.go       # HTTP 客戶端配置
├── retry.go        # 重試機制
├── message_test.go # 資料結構測試
├── samhook_test.go # 核心功能測試
└── README.md       # 專案說明
```

## 授權

本專案採用 MIT 授權 - 詳見 [LICENSE](LICENSE) 文件。
