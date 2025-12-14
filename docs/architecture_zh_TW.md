# 架構文檔

[English](architecture.md) | 繁體中文

## 專案架構概述

samhook 是一個輕量級的 Go 函式庫，設計目標是提供簡潔、易用的 API 來發送 Slack 和 Mattermost webhook 訊息。

## 設計原則

### 1. 輕量級設計

- **極少外部依賴**: 僅使用高性能 JSON 庫
- **最小化抽象**: 直接映射 Slack/Mattermost webhook API
- **簡單易用**: 提供直觀的 API 介面

### 2. 類型安全

- **強類型**: 所有資料結構都有明確的類型定義
- **JSON 標籤**: 使用標準 JSON 標籤進行序列化
- **可選欄位**: 使用 `omitempty` 標籤，僅序列化有值的欄位

### 3. 靈活性

- **多種輸入方式**: 支援結構體和 Reader 兩種輸入方式
- **鏈式調用**: 支援方法鏈式調用以提升開發體驗
- **可擴展**: 易於添加新功能

## 模組結構

### 核心模組

```
samhook/
├── message.go    # 資料結構定義
├── samhook.go    # 核心功能實現
├── error.go      # 錯誤類型定義
├── client.go     # HTTP 客戶端配置
└── retry.go      # 重試機制
```

### message.go

定義了三個核心資料結構：

1. **Message** - 訊息主體
   - 包含基本訊息屬性（文字、使用者名稱、圖示等）
   - 包含附件列表

2. **Attachment** - 附件結構
   - 支援豐富的格式化選項
   - 支援欄位、圖片、作者資訊等

3. **Field** - 欄位結構
   - 用於在附件中顯示結構化資料
   - 支援短欄位並排顯示

### samhook.go

實現了核心功能：

1. **Send 函數**
   - 接收 Message 結構
   - 序列化為 JSON
   - 發送 HTTP POST 請求
   - 檢查 HTTP 狀態碼並返回詳細錯誤

2. **SendReader 函數**
   - 接收 io.Reader
   - 直接發送 HTTP POST 請求
   - 適用於已有 JSON 資料的場景
   - 檢查 HTTP 狀態碼並返回詳細錯誤

3. **AddAttachment 方法**
   - 為 Message 添加單個附件
   - 返回 *Message 以支援鏈式調用

4. **AddAttachments 方法**
   - 為 Message 添加多個附件
   - 返回 *Message 以支援鏈式調用

### error.go

實現了錯誤處理功能：

1. **WebhookError 類型**
   - 自訂錯誤類型，提供詳細錯誤資訊
   - 支援錯誤分類（網路錯誤、序列化錯誤、API 錯誤）
   - 提供錯誤代碼和詳細訊息

2. **錯誤構造函數**
   - `NewNetworkError()` - 創建網路錯誤
   - `NewSerializationError()` - 創建序列化錯誤
   - `NewAPIError()` - 創建 API 錯誤

### client.go

實現了 HTTP 客戶端配置功能：

1. **SendWithOptions 函數**
   - 支援自訂 HTTP 客戶端配置
   - 使用選項模式（Options Pattern）
   - 支援超時配置

2. **SendWithContext 函數**
   - 支援 Context 超時和取消
   - 使用 `http.NewRequestWithContext()`

3. **ClientOption 類型**
   - `WithTimeout()` - 設置超時
   - `WithClient()` - 使用自訂客戶端

### retry.go

實現了重試機制：

1. **SendWithRetry 函數**
   - 支援自動重試失敗的請求
   - 可配置重試次數和間隔
   - 智能判斷哪些錯誤可重試

2. **ExponentialBackoff 類型**
   - 實現指數退避策略
   - 支援隨機抖動（jitter）
   - 可配置初始間隔和最大間隔

## 資料流程

### 標準發送流程

```
使用者代碼
    ↓
創建 Message 結構
    ↓
（可選）添加附件
    ↓
調用 Send() 函數
    ↓
JSON 序列化
    ↓
HTTP POST 請求
    ↓
Webhook URL
```

### SendReader 流程

```
使用者代碼
    ↓
準備 JSON 資料（io.Reader）
    ↓
調用 SendReader() 函數
    ↓
HTTP POST 請求（直接使用 Reader）
    ↓
Webhook URL
```

## 依賴關係

### 標準庫依賴

- `bytes` - 用於建立請求主體
- `context` - Context 支援（用於超時和取消）
- `fmt` - 字串格式化
- `io` - Reader 介面
- `math` - 數學運算（用於指數退避）
- `math/rand` - 隨機數生成（用於抖動）
- `net` - 網路錯誤分類
- `net/http` - HTTP 客戶端
- `net/url` - URL 錯誤處理
- `strings` - 字串操作
- `time` - 時間和超時控制

### 外部依賴

- `github.com/bytedance/sonic` - 高性能 JSON 序列化/反序列化庫

專案使用高性能的 JSON 庫，確保：
- 編譯速度快
- 執行效能高
- 記憶體使用優化
- 依賴管理簡單

## 擴展性設計

### 已實現的功能

1. ✅ **可配置 HTTP 客戶端**: 支援自訂超時、自訂客戶端（`SendWithOptions`）
2. ✅ **回應處理**: 檢查 HTTP 狀態碼並返回詳細錯誤
3. ✅ **重試機制**: 支援自動重試，包含指數退避策略
4. ✅ **Context 支援**: 支援超時和取消控制
5. ✅ **詳細錯誤處理**: 自訂錯誤類型，提供錯誤分類和詳細資訊

### 未來擴展方向

1. **批次發送**: 支援一次發送多個訊息
2. **非同步發送**: 提供輔助函數簡化非同步發送
3. **速率限制處理**: 自動處理 429 錯誤和請求佇列

## 相容性

### Slack Webhook API

samhook 完全相容 Slack Incoming Webhooks API，支援：
- 基本訊息格式
- 附件格式
- 欄位格式
- 所有標準欄位

### Mattermost Webhook API

Mattermost 的 webhook API 與 Slack 高度相容，因此 samhook 也可以直接用於 Mattermost。

## 錯誤處理策略

### 當前實現

- ✅ 所有函數返回 `error`（可能是 `*WebhookError`）
- ✅ 錯誤分類：網路錯誤、序列化錯誤、API 錯誤
- ✅ 詳細錯誤資訊：包含狀態碼、回應體、錯誤代碼等
- ✅ 錯誤分類方法：`IsNetworkError()`, `IsAPIError()`, `IsSerializationError()`
- ✅ 詳細錯誤訊息：`DetailedMessage()` 提供多行格式

### 錯誤類型

1. **網路錯誤** (`ErrorTypeNetwork`): HTTP 請求失敗、連線錯誤等
2. **序列化錯誤** (`ErrorTypeSerialization`): JSON 序列化失敗
3. **API 錯誤** (`ErrorTypeAPI`): HTTP 狀態碼非 200
4. **未知錯誤** (`ErrorTypeUnknown`): 無法分類的錯誤

### 錯誤處理流程

```
發送請求
    ↓
發生錯誤？
    ↓ 是
創建 WebhookError
    ↓
分類錯誤類型
    ↓
返回 WebhookError
    ↓
使用者可以：
- 檢查錯誤類型
- 獲取狀態碼
- 獲取詳細訊息
- 實現智能處理
```

## 效能考量

### 當前實現

- 使用 `http.DefaultClient`（共用連線池）
- 同步發送（阻塞直到完成）
- 高性能序列化（使用 sonic JSON 庫）

### 效能特點

- **低延遲**: 直接 HTTP 請求，無額外開銷
- **低記憶體**: 最小化記憶體分配
- **高吞吐**: 可並發使用（每個 goroutine 獨立）

## 安全性考量

### 當前實現

- 不驗證 webhook URL
- 不加密傳輸（依賴 HTTPS）
- 不驗證回應

### 建議

1. **URL 驗證**: 驗證 webhook URL 格式
2. **HTTPS 強制**: 在生產環境強制使用 HTTPS
3. **回應驗證**: 檢查 HTTP 狀態碼

## 測試策略

### 已實現的測試

1. ✅ **單元測試**: 
   - 資料結構序列化測試（`message_test.go`）
   - 錯誤類型測試（`samhook_test.go`）
   - 客戶端配置測試（`client_test.go`）
   - 重試機制測試（`retry_test.go`）

2. ✅ **整合測試**: 
   - 使用 `httptest` 套件建立 mock HTTP 伺服器
   - 測試各種 HTTP 狀態碼場景
   - 測試錯誤處理邏輯

3. **端到端測試**: 
   - 可選的整合測試（需要實際 webhook URL）
   - 使用 build tags 標記為整合測試

### 測試覆蓋率

當前測試覆蓋率約為 62%，涵蓋：
- 所有資料結構的序列化
- 核心發送功能
- 錯誤處理邏輯
- 客戶端配置
- 重試機制

## 維護性

### 程式碼組織

- **清晰的模組劃分**: 資料結構與功能分離
- **一致的命名**: 遵循 Go 命名慣例
- **完整的註解**: 所有公開 API 都有中文註解

### 可維護性特點

- **簡單的程式碼結構**: 易於理解和修改
- **最小化複雜度**: 避免過度設計
- **標準化格式**: 使用 Go 標準格式
