# samhook

English | [繁體中文](README_zh_TW.md)

A lightweight Go library for sending Slack and Mattermost webhook messages.

## Features

- ✅ **Lightweight**: Minimal external dependencies, uses high-performance JSON library
- ✅ **Type-safe**: Complete Go type definitions
- ✅ **Method chaining**: Supports method chaining for better developer experience
- ✅ **Flexible input**: Supports both struct and Reader input methods
- ✅ **Error handling**: Detailed error types and classifications
- ✅ **Configurable**: Supports custom HTTP client and timeout settings
- ✅ **Retry mechanism**: Optional retry functionality with exponential backoff

## Installation

```bash
go get github.com/circleyu/samhook
```

## Quick Start

### Basic Usage

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

### Using Attachments

```go
msg := samhook.Message{
    Text: "System Notification",
}

attachment := samhook.Attachment{
    Color: samhook.Good,
    Title: "Operation Successful",
    Text:  "All tasks completed",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### Error Handling

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        if webhookErr.IsNetworkError() {
            // Handle network error, can retry
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            if statusCode == 429 {
                // Handle rate limiting
            }
        }
    }
}
```

### Using Custom Client

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

// Use custom timeout
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)

// Use Context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := samhook.SendWithContext(ctx, webhookURL, msg)
```

### Using Retry Mechanism

```go
opts := samhook.DefaultRetryOptions
opts.MaxRetries = 5

err := samhook.SendWithRetry(webhookURL, msg, opts)
```

## Documentation

- [API Documentation](docs/api.md) - Complete API reference
- [Architecture Documentation](docs/architecture.md) - Project architecture and design
- [Usage Examples](docs/examples.md) - Detailed usage examples

## Project Structure

```
samhook/
├── go.mod          # Go module definition
├── message.go      # Message data structure definitions
├── samhook.go      # Core sending functionality
├── error.go        # Error type definitions
├── client.go       # HTTP client configuration
├── retry.go        # Retry mechanism
├── message_test.go # Data structure tests
├── samhook_test.go # Core functionality tests
└── README.md       # Project documentation
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
