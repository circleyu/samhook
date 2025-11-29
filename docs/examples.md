# 使用範例

本文檔提供 samhook 函式庫的各種使用範例。

## 基本使用

### 發送簡單文字訊息

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
    
    if err := samhook.Send(webhookURL, msg); err != nil {
        log.Fatal(err)
    }
}
```

### 發送帶有 Emoji 圖示的訊息

```go
msg := samhook.Message{
    Text:      "Notification",
    Username:  "bot",
    IconEmoji: ":robot_face:",
}
samhook.Send(webhookURL, msg)
```

### 發送帶有自訂圖示的訊息

```go
msg := samhook.Message{
    Text:     "Notification",
    Username: "bot",
    IconURL:  "https://example.com/icon.png",
}
samhook.Send(webhookURL, msg)
```

## 使用附件

### 基本附件

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

### 帶有欄位的附件

```go
msg := samhook.Message{
    Text: "部署通知",
}

attachment := samhook.Attachment{
    Color: samhook.Good,
    Title: "部署成功",
    Fields: []samhook.Field{
        {
            Title: "環境",
            Value: "生產環境",
            Short: true,
        },
        {
            Title: "版本",
            Value: "v1.2.3",
            Short: true,
        },
        {
            Title: "部署時間",
            Value: "2024-01-15 10:30:00",
            Short: false,
        },
    },
    Footer: "部署系統",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### 警告訊息

```go
msg := samhook.Message{
    Text: "警告通知",
}

attachment := samhook.Attachment{
    Color:   samhook.Warning,
    Title:   "資源使用率過高",
    Pretext: "請注意以下資源使用情況：",
    Fields: []samhook.Field{
        {
            Title: "CPU 使用率",
            Value: "85%",
            Short: true,
        },
        {
            Title: "記憶體使用率",
            Value: "78%",
            Short: true,
        },
    },
    Footer: "監控系統",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### 錯誤訊息

```go
msg := samhook.Message{
    Text: "錯誤通知",
}

attachment := samhook.Attachment{
    Color: samhook.Danger,
    Title: "系統錯誤",
    Text:  "資料庫連接失敗",
    Fields: []samhook.Field{
        {
            Title: "錯誤代碼",
            Value: "DB_CONN_001",
            Short: true,
        },
        {
            Title: "發生時間",
            Value: "2024-01-15 10:30:00",
            Short: true,
        },
    },
    Footer: "錯誤追蹤系統",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

## 鏈式調用

### 添加多個附件

```go
msg := samhook.Message{Text: "多個通知"}
    .AddAttachment(samhook.Attachment{
        Color: samhook.Good,
        Title: "任務 1 完成",
    })
    .AddAttachment(samhook.Attachment{
        Color: samhook.Warning,
        Title: "任務 2 需要關注",
    })

samhook.Send(webhookURL, msg)
```

### 批量添加附件

```go
attachments := []samhook.Attachment{
    {
        Color: samhook.Good,
        Title: "服務 A 正常",
        Text:  "所有檢查通過",
    },
    {
        Color: samhook.Good,
        Title: "服務 B 正常",
        Text:  "所有檢查通過",
    },
}

msg := samhook.Message{Text: "健康檢查報告"}
    .AddAttachments(attachments)

samhook.Send(webhookURL, msg)
```

## 進階使用

### 帶有作者資訊的附件

```go
attachment := samhook.Attachment{
    Color:      samhook.Good,
    AuthorName: "部署機器人",
    AuthorLink: "https://example.com/bot",
    AuthorIcon: "https://example.com/bot-icon.png",
    Title:      "部署完成",
    TitleLink:  "https://example.com/deployment/123",
    Text:       "應用程式已成功部署到生產環境",
}

msg := samhook.Message{}
msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### 帶有圖片的附件

```go
attachment := samhook.Attachment{
    Color:    samhook.Good,
    Title:    "圖表報告",
    Text:     "本週數據分析",
    ImageURL: "https://example.com/chart.png",
    ThumbURL: "https://example.com/thumb.png",
}

msg := samhook.Message{}
msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### 使用 SendReader 發送

當您已經有 JSON 格式的訊息時：

```go
package main

import (
    "bytes"
    "github.com/circleyu/samhook"
)

func main() {
    webhookURL := "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
    
    // 假設您已經有 JSON 格式的訊息
    jsonMsg := `{
        "text": "Hello",
        "username": "bot",
        "attachments": [{
            "color": "#00FF00",
            "title": "Success"
        }]
    }`
    
    reader := bytes.NewReader([]byte(jsonMsg))
    samhook.SendReader(webhookURL, reader)
}
```

## 實際應用場景

### 部署通知

```go
func notifyDeployment(webhookURL, env, version string, success bool) error {
    msg := samhook.Message{
        Username: "部署系統",
        IconEmoji: ":rocket:",
    }
    
    attachment := samhook.Attachment{
        Title: "部署通知",
        Fields: []samhook.Field{
            {Title: "環境", Value: env, Short: true},
            {Title: "版本", Value: version, Short: true},
        },
        Footer: "CI/CD 系統",
    }
    
    if success {
        attachment.Color = samhook.Good
        attachment.Text = "部署成功"
    } else {
        attachment.Color = samhook.Danger
        attachment.Text = "部署失敗"
    }
    
    msg.AddAttachment(attachment)
    return samhook.Send(webhookURL, msg)
}
```

## 錯誤處理

### 基本錯誤處理

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    log.Printf("發送失敗: %v", err)
}
```

### 使用 WebhookError 進行詳細錯誤處理

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        // 檢查錯誤類型
        if webhookErr.IsNetworkError() {
            log.Printf("網路錯誤: %v", webhookErr)
            // 可以重試
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            log.Printf("API 錯誤 (狀態碼 %d): %v", statusCode, webhookErr)
            
            // 處理特定狀態碼
            switch statusCode {
            case 401:
                log.Println("認證失敗，請檢查 webhook URL")
            case 429:
                log.Println("速率限制，請稍後重試")
            case 500, 502, 503, 504:
                log.Println("伺服器錯誤，可以重試")
            }
        } else if webhookErr.IsSerializationError() {
            log.Printf("序列化錯誤: %v", webhookErr)
        }
        
        // 獲取詳細錯誤訊息
        log.Println(webhookErr.DetailedMessage())
    } else {
        log.Printf("未知錯誤: %v", err)
    }
}
```

## 進階配置

### 使用自訂超時

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "重要通知"}

// 設置 30 秒超時
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)
```

### 使用自訂 HTTP 客戶端

```go
import (
    "net/http"
    "time"
    "github.com/circleyu/samhook"
)

customClient := &http.Client{
    Timeout: 20 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 10,
    },
}

err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithClient(customClient),
)
```

### 使用 Context 進行超時和取消

```go
import (
    "context"
    "time"
    "github.com/circleyu/samhook"
)

// 使用 Context 超時
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := samhook.SendWithContext(ctx, webhookURL, msg)

// 使用 Context 取消
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // 5 秒後取消請求
}()

err := samhook.SendWithContext(ctx, webhookURL, msg)
```

## 重試機制

### 基本重試

```go
import "github.com/circleyu/samhook"

msg := samhook.Message{Text: "重要通知"}

// 使用預設重試選項（最多重試 3 次，間隔 1 秒）
opts := samhook.DefaultRetryOptions
err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### 自訂重試配置

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "重要通知"}

opts := samhook.RetryOptions{
    MaxRetries: 5,
    Interval:   2 * time.Second,
    Backoff: &samhook.ExponentialBackoff{
        InitialInterval: 1 * time.Second,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,
        Jitter:          true, // 添加隨機抖動避免雷群效應
    },
}

err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### 重試場景範例

```go
// 處理暫時性網路故障
func sendWithRetry(webhookURL string, msg samhook.Message) error {
    opts := samhook.DefaultRetryOptions
    opts.MaxRetries = 3
    
    err := samhook.SendWithRetry(webhookURL, msg, opts)
    if err != nil {
        // 檢查是否為可重試的錯誤
        if webhookErr, ok := err.(*samhook.WebhookError); ok {
            if webhookErr.IsNetworkError() {
                log.Printf("網路錯誤，已重試 %d 次: %v", opts.MaxRetries, err)
            } else if webhookErr.IsAPIError() && webhookErr.GetStatusCode() >= 500 {
                log.Printf("伺服器錯誤，已重試 %d 次: %v", opts.MaxRetries, err)
            }
        }
    }
    return err
}
```

### 監控告警

```go
func sendAlert(webhookURL, alertType, message string, severity int) error {
    msg := samhook.Message{
        Username: "監控系統",
        IconEmoji: ":warning:",
    }
    
    var color string
    switch severity {
    case 1:
        color = samhook.Danger
    case 2:
        color = samhook.Warning
    default:
        color = samhook.Good
    }
    
    attachment := samhook.Attachment{
        Color: color,
        Title: alertType,
        Text:  message,
        Footer: "監控系統",
    }
    
    msg.AddAttachment(attachment)
    return samhook.Send(webhookURL, msg)
}
```

## 錯誤處理

### 基本錯誤處理

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    log.Printf("發送失敗: %v", err)
}
```

### 使用 WebhookError 進行詳細錯誤處理

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        // 檢查錯誤類型
        if webhookErr.IsNetworkError() {
            log.Printf("網路錯誤: %v", webhookErr)
            // 可以重試
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            log.Printf("API 錯誤 (狀態碼 %d): %v", statusCode, webhookErr)
            
            // 處理特定狀態碼
            switch statusCode {
            case 401:
                log.Println("認證失敗，請檢查 webhook URL")
            case 429:
                log.Println("速率限制，請稍後重試")
            case 500, 502, 503, 504:
                log.Println("伺服器錯誤，可以重試")
            }
        } else if webhookErr.IsSerializationError() {
            log.Printf("序列化錯誤: %v", webhookErr)
        }
        
        // 獲取詳細錯誤訊息
        log.Println(webhookErr.DetailedMessage())
    } else {
        log.Printf("未知錯誤: %v", err)
    }
}
```

## 進階配置

### 使用自訂超時

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "重要通知"}

// 設置 30 秒超時
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)
```

### 使用自訂 HTTP 客戶端

```go
import (
    "net/http"
    "time"
    "github.com/circleyu/samhook"
)

customClient := &http.Client{
    Timeout: 20 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 10,
    },
}

err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithClient(customClient),
)
```

### 使用 Context 進行超時和取消

```go
import (
    "context"
    "time"
    "github.com/circleyu/samhook"
)

// 使用 Context 超時
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := samhook.SendWithContext(ctx, webhookURL, msg)

// 使用 Context 取消
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // 5 秒後取消請求
}()

err := samhook.SendWithContext(ctx, webhookURL, msg)
```

## 重試機制

### 基本重試

```go
import "github.com/circleyu/samhook"

msg := samhook.Message{Text: "重要通知"}

// 使用預設重試選項（最多重試 3 次，間隔 1 秒）
opts := samhook.DefaultRetryOptions
err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### 自訂重試配置

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "重要通知"}

opts := samhook.RetryOptions{
    MaxRetries: 5,
    Interval:   2 * time.Second,
    Backoff: &samhook.ExponentialBackoff{
        InitialInterval: 1 * time.Second,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,
        Jitter:          true, // 添加隨機抖動避免雷群效應
    },
}

err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### 重試場景範例

```go
// 處理暫時性網路故障
func sendWithRetry(webhookURL string, msg samhook.Message) error {
    opts := samhook.DefaultRetryOptions
    opts.MaxRetries = 3
    
    err := samhook.SendWithRetry(webhookURL, msg, opts)
    if err != nil {
        // 檢查是否為可重試的錯誤
        if webhookErr, ok := err.(*samhook.WebhookError); ok {
            if webhookErr.IsNetworkError() {
                log.Printf("網路錯誤，已重試 %d 次: %v", opts.MaxRetries, err)
            } else if webhookErr.IsAPIError() && webhookErr.GetStatusCode() >= 500 {
                log.Printf("伺服器錯誤，已重試 %d 次: %v", opts.MaxRetries, err)
            }
        }
    }
    return err
}
```

### 任務完成通知

```go
func notifyTaskCompletion(webhookURL, taskName, duration string) error {
    msg := samhook.Message{
        Text:     "任務完成通知",
        Username: "任務系統",
    }
    
    attachment := samhook.Attachment{
        Color: samhook.Good,
        Title: "任務已完成",
        Fields: []samhook.Field{
            {Title: "任務名稱", Value: taskName, Short: false},
            {Title: "執行時間", Value: duration, Short: true},
        },
    }
    
    msg.AddAttachment(attachment)
    return samhook.Send(webhookURL, msg)
}
```

## 錯誤處理

### 基本錯誤處理

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    log.Printf("發送失敗: %v", err)
}
```

### 使用 WebhookError 進行詳細錯誤處理

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        // 檢查錯誤類型
        if webhookErr.IsNetworkError() {
            log.Printf("網路錯誤: %v", webhookErr)
            // 可以重試
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            log.Printf("API 錯誤 (狀態碼 %d): %v", statusCode, webhookErr)
            
            // 處理特定狀態碼
            switch statusCode {
            case 401:
                log.Println("認證失敗，請檢查 webhook URL")
            case 429:
                log.Println("速率限制，請稍後重試")
            case 500, 502, 503, 504:
                log.Println("伺服器錯誤，可以重試")
            }
        } else if webhookErr.IsSerializationError() {
            log.Printf("序列化錯誤: %v", webhookErr)
        }
        
        // 獲取詳細錯誤訊息
        log.Println(webhookErr.DetailedMessage())
    } else {
        log.Printf("未知錯誤: %v", err)
    }
}
```

## 進階配置

### 使用自訂超時

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "重要通知"}

// 設置 30 秒超時
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)
```

### 使用自訂 HTTP 客戶端

```go
import (
    "net/http"
    "time"
    "github.com/circleyu/samhook"
)

customClient := &http.Client{
    Timeout: 20 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 10,
    },
}

err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithClient(customClient),
)
```

### 使用 Context 進行超時和取消

```go
import (
    "context"
    "time"
    "github.com/circleyu/samhook"
)

// 使用 Context 超時
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := samhook.SendWithContext(ctx, webhookURL, msg)

// 使用 Context 取消
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // 5 秒後取消請求
}()

err := samhook.SendWithContext(ctx, webhookURL, msg)
```

## 重試機制

### 基本重試

```go
import "github.com/circleyu/samhook"

msg := samhook.Message{Text: "重要通知"}

// 使用預設重試選項（最多重試 3 次，間隔 1 秒）
opts := samhook.DefaultRetryOptions
err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### 自訂重試配置

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "重要通知"}

opts := samhook.RetryOptions{
    MaxRetries: 5,
    Interval:   2 * time.Second,
    Backoff: &samhook.ExponentialBackoff{
        InitialInterval: 1 * time.Second,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,
        Jitter:          true, // 添加隨機抖動避免雷群效應
    },
}

err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### 重試場景範例

```go
// 處理暫時性網路故障
func sendWithRetry(webhookURL string, msg samhook.Message) error {
    opts := samhook.DefaultRetryOptions
    opts.MaxRetries = 3
    
    err := samhook.SendWithRetry(webhookURL, msg, opts)
    if err != nil {
        // 檢查是否為可重試的錯誤
        if webhookErr, ok := err.(*samhook.WebhookError); ok {
            if webhookErr.IsNetworkError() {
                log.Printf("網路錯誤，已重試 %d 次: %v", opts.MaxRetries, err)
            } else if webhookErr.IsAPIError() && webhookErr.GetStatusCode() >= 500 {
                log.Printf("伺服器錯誤，已重試 %d 次: %v", opts.MaxRetries, err)
            }
        }
    }
    return err
}
```

