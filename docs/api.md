# API 文檔

## 概述

本文檔提供 samhook 函式庫的完整 API 參考。

## 常數

### 顏色常數

預定義的顏色常數，用於設定附件的顏色：

```go
const Warning string = "#FFBB00"  // 警告顏色（黃色）
const Danger string = "#FF0000"   // 危險顏色（紅色）
const Good string = "#00FF00"      // 正常資訊顏色（綠色）
```

## 資料結構

### Message

訊息主體結構，代表要發送的 webhook 訊息。

```go
type Message struct {
    Parse       string       `json:"parse,omitempty"`
    Username    string       `json:"username,omitempty"`
    IconURL     string       `json:"icon_url,omitempty"`
    IconEmoji   string       `json:"icon_emoji,omitempty"`
    Channel     string       `json:"channel,omitempty"`
    Text        string       `json:"text,omitempty"`
    Attachments []Attachment `json:"attachments,omitempty"`
}
```

#### 欄位說明

- `Parse` - 解析模式（可選）
- `Username` - 發送者使用者名稱（可選）
- `IconURL` - 圖示 URL（可選）
- `IconEmoji` - 圖示 Emoji（可選）
- `Channel` - 目標頻道（可選）
- `Text` - 訊息文字內容（可選）
- `Attachments` - 附件列表（可選）

### Attachment

附件結構，用於豐富訊息內容。

```go
type Attachment struct {
    Fallback   string  `json:"fallback,omitempty"`
    Color      string  `json:"color,omitempty"`
    Pretext    string  `json:"pretext,omitempty"`
    AuthorName string  `json:"author_name,omitempty"`
    AuthorLink string  `json:"author_link,omitempty"`
    AuthorIcon string  `json:"author_icon,omitempty"`
    Title      string  `json:"title,omitempty"`
    TitleLink  string  `json:"title_link,omitempty"`
    Text       string  `json:"text,omitempty"`
    ImageURL   string  `json:"image_url,omitempty"`
    Fields     []Field `json:"fields,omitempty"`
    Footer     string  `json:"footer,omitempty"`
    FooterIcon string  `json:"footer_icon,omitempty"`
    ThumbURL   string  `json:"thumb_url,omitempty"`
}
```

#### 欄位說明

- `Fallback` - 回退文字（用於不支援附件的客戶端）
- `Color` - 附件左側顏色條（可使用預定義常數）
- `Pretext` - 附件前的文字
- `AuthorName` - 作者名稱
- `AuthorLink` - 作者連結
- `AuthorIcon` - 作者圖示 URL
- `Title` - 附件標題
- `TitleLink` - 標題連結
- `Text` - 附件文字內容
- `ImageURL` - 圖片 URL
- `Fields` - 欄位列表
- `Footer` - 頁腳文字
- `FooterIcon` - 頁腳圖示 URL
- `ThumbURL` - 縮圖 URL

### Field

欄位結構，用於在附件中顯示鍵值對資訊。

```go
type Field struct {
    Title string `json:"title,omitempty"`
    Value string `json:"value,omitempty"`
    Short bool   `json:"short,omitempty"`
}
```

#### 欄位說明

- `Title` - 欄位標題
- `Value` - 欄位值
- `Short` - 是否為短欄位（用於並排顯示）

## 函數

### Send

發送訊息到指定的 webhook URL。

```go
func Send(url string, msg Message) error
```

#### 參數

- `url` - webhook URL（字串）
- `msg` - 要發送的訊息結構

#### 返回值

- `error` - 如果發送失敗則返回錯誤，否則返回 nil

#### 行為

1. 將訊息結構序列化為 JSON
2. 建立 HTTP POST 請求
3. 設定 Content-Type 為 application/json
4. 發送請求並關閉回應主體

### SendReader

從 io.Reader 發送訊息到指定的 webhook URL。

```go
func SendReader(url string, r io.Reader) error
```

#### 參數

- `url` - webhook URL（字串）
- `r` - 包含 JSON 訊息的 Reader

#### 返回值

- `error` - 如果發送失敗則返回錯誤，否則返回 nil

#### 使用場景

當您已經有 JSON 格式的訊息資料（例如從檔案讀取或從其他來源獲取）時，可以使用此函數直接發送。

## 方法

### AddAttachment

為訊息添加單個附件。返回訊息指標以支援鏈式調用。

```go
func (m *Message) AddAttachment(attachment Attachment) *Message
```

#### 參數

- `attachment` - 要添加的附件

#### 返回值

- `*Message` - 訊息指標（用於鏈式調用）

#### 範例

```go
msg := samhook.Message{Text: "Alert!"}
msg.AddAttachment(samhook.Attachment{
    Color: samhook.Danger,
    Title: "Error",
    Text:  "Something went wrong",
})
```

### AddAttachments

為訊息添加多個附件。返回訊息指標以支援鏈式調用。

```go
func (m *Message) AddAttachments(attachments []Attachment) *Message
```

#### 參數

- `attachments` - 要添加的附件列表

#### 返回值

- `*Message` - 訊息指標（用於鏈式調用）

#### 範例

```go
msg := samhook.Message{Text: "Multiple alerts"}
attachments := []samhook.Attachment{
    {Color: samhook.Warning, Title: "Warning 1"},
    {Color: samhook.Danger, Title: "Error 1"},
}
msg.AddAttachments(attachments)
```

### SendWithOptions

使用選項模式發送訊息，支援自訂 HTTP 客戶端配置。

```go
func SendWithOptions(url string, msg Message, opts ...ClientOption) error
```

#### 參數

- `url` - webhook URL（字串）
- `msg` - 要發送的訊息結構
- `opts` - 客戶端選項（可變參數）

#### 返回值

- `error` - 如果發送失敗則返回錯誤，否則返回 nil

#### 範例

```go
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)
```

### SendWithContext

使用 Context 發送訊息，支援超時和取消。

```go
func SendWithContext(ctx context.Context, url string, msg Message, opts ...ClientOption) error
```

#### 參數

- `ctx` - Context 用於超時和取消控制
- `url` - webhook URL（字串）
- `msg` - 要發送的訊息結構
- `opts` - 客戶端選項（可變參數）

#### 返回值

- `error` - 如果發送失敗則返回錯誤，否則返回 nil

#### 範例

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := samhook.SendWithContext(ctx, webhookURL, msg)
```

### SendWithRetry

帶重試機制的發送函數。

```go
func SendWithRetry(url string, msg Message, opts RetryOptions) error
```

#### 參數

- `url` - webhook URL（字串）
- `msg` - 要發送的訊息結構
- `opts` - 重試選項

#### 返回值

- `error` - 如果所有重試都失敗則返回錯誤，否則返回 nil

#### 重試條件

- 網路錯誤會自動重試
- 5xx 伺服器錯誤會自動重試
- 429 速率限制錯誤會自動重試
- 4xx 客戶端錯誤不會重試

#### 範例

```go
opts := samhook.DefaultRetryOptions
opts.MaxRetries = 5
err := samhook.SendWithRetry(webhookURL, msg, opts)
```

## 選項類型

### ClientOption

HTTP 客戶端選項函數類型。

```go
type ClientOption func(*http.Client)
```

#### 可用選項

- `WithTimeout(timeout time.Duration)` - 設置超時時間
- `WithClient(client *http.Client)` - 使用自訂 HTTP 客戶端

### RetryOptions

重試選項結構。

```go
type RetryOptions struct {
    MaxRetries int
    Interval   time.Duration
    Backoff    *ExponentialBackoff
}
```

#### 欄位說明

- `MaxRetries` - 最大重試次數
- `Interval` - 固定重試間隔（如果未設置 Backoff）
- `Backoff` - 指數退避配置（可選）

### ExponentialBackoff

指數退避配置。

```go
type ExponentialBackoff struct {
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
    Jitter          bool
}
```

#### 欄位說明

- `InitialInterval` - 初始重試間隔
- `MaxInterval` - 最大重試間隔
- `Multiplier` - 退避倍數（通常為 2.0）
- `Jitter` - 是否添加隨機抖動

## 錯誤類型

### WebhookError

自訂錯誤類型，提供詳細的錯誤資訊。

```go
type WebhookError struct {
    Type         string
    StatusCode   int
    Message      string
    ResponseBody string
    Err          error
    URL          string
}
```

#### 方法

- `Error() string` - 實現 error 介面
- `Unwrap() error` - 返回原始錯誤
- `IsNetworkError() bool` - 判斷是否為網路錯誤
- `IsSerializationError() bool` - 判斷是否為序列化錯誤
- `IsAPIError() bool` - 判斷是否為 API 錯誤
- `GetStatusCode() int` - 返回 HTTP 狀態碼
- `GetResponseBody() string` - 返回 API 回應體
- `GetErrorCode() string` - 返回錯誤代碼
- `DetailedMessage() string` - 返回詳細的多行錯誤訊息

#### 錯誤類型常數

```go
const (
    ErrorTypeNetwork       = "network"
    ErrorTypeSerialization = "serialization"
    ErrorTypeAPI           = "api"
    ErrorTypeUnknown       = "unknown"
)
```

#### 錯誤代碼常數

```go
const (
    ErrorCodeNetworkTimeout     = "NETWORK_TIMEOUT"
    ErrorCodeNetworkConnection  = "NETWORK_CONNECTION"
    ErrorCodeNetworkDNS         = "NETWORK_DNS"
    ErrorCodeSerializationJSON  = "SERIALIZATION_JSON"
    ErrorCodeAPIUnauthorized    = "API_UNAUTHORIZED"
    ErrorCodeAPIForbidden       = "API_FORBIDDEN"
    ErrorCodeAPINotFound        = "API_NOT_FOUND"
    ErrorCodeAPIRateLimit       = "API_RATE_LIMIT"
    ErrorCodeAPIServerError    = "API_SERVER_ERROR"
)
```

#### 錯誤構造函數

- `NewNetworkError(url string, err error) *WebhookError` - 創建網路錯誤
- `NewSerializationError(err error) *WebhookError` - 創建序列化錯誤
- `NewAPIError(url string, statusCode int, responseBody string) *WebhookError` - 創建 API 錯誤

## 錯誤處理

所有函數在發生錯誤時都會返回 `error`。從 `Send()` 和 `SendReader()` 返回的錯誤可能是 `*WebhookError` 類型，提供更詳細的錯誤資訊。

### 常見錯誤情況

- **網路錯誤**: 連線失敗、超時等
- **序列化錯誤**: JSON 序列化失敗
- **API 錯誤**: HTTP 狀態碼非 200（如 401、429、500 等）

### 錯誤處理範例

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        if webhookErr.IsNetworkError() {
            // 處理網路錯誤
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            // 根據狀態碼處理
        }
    }
}
```

建議在生產環境中適當處理這些錯誤，並根據錯誤類型實現相應的重試或回退策略。

