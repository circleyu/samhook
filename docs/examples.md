# Usage Examples

English | [繁體中文](examples_zh_TW.md)

This document provides various usage examples for the samhook library.

## Basic Usage

### Sending a Simple Text Message

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

### Sending a Message with Emoji Icon

```go
msg := samhook.Message{
    Text:      "Notification",
    Username:  "bot",
    IconEmoji: ":robot_face:",
}
samhook.Send(webhookURL, msg)
```

### Sending a Message with Custom Icon

```go
msg := samhook.Message{
    Text:     "Notification",
    Username: "bot",
    IconURL:  "https://example.com/icon.png",
}
samhook.Send(webhookURL, msg)
```

## Using Attachments

### Basic Attachment

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

### Attachment with Fields

```go
msg := samhook.Message{
    Text: "Deployment Notification",
}

attachment := samhook.Attachment{
    Color: samhook.Good,
    Title: "Deployment Successful",
    Fields: []samhook.Field{
        {
            Title: "Environment",
            Value: "Production",
            Short: true,
        },
        {
            Title: "Version",
            Value: "v1.2.3",
            Short: true,
        },
        {
            Title: "Deployment Time",
            Value: "2024-01-15 10:30:00",
            Short: false,
        },
    },
    Footer: "Deployment System",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### Warning Message

```go
msg := samhook.Message{
    Text: "Warning Notification",
}

attachment := samhook.Attachment{
    Color:   samhook.Warning,
    Title:   "High Resource Usage",
    Pretext: "Please note the following resource usage:",
    Fields: []samhook.Field{
        {
            Title: "CPU Usage",
            Value: "85%",
            Short: true,
        },
        {
            Title: "Memory Usage",
            Value: "78%",
            Short: true,
        },
    },
    Footer: "Monitoring System",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### Error Message

```go
msg := samhook.Message{
    Text: "Error Notification",
}

attachment := samhook.Attachment{
    Color: samhook.Danger,
    Title: "System Error",
    Text:  "Database connection failed",
    Fields: []samhook.Field{
        {
            Title: "Error Code",
            Value: "DB_CONN_001",
            Short: true,
        },
        {
            Title: "Occurred At",
            Value: "2024-01-15 10:30:00",
            Short: true,
        },
    },
    Footer: "Error Tracking System",
}

msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

## Method Chaining

### Adding Multiple Attachments

```go
msg := samhook.Message{Text: "Multiple Notifications"}
    .AddAttachment(samhook.Attachment{
        Color: samhook.Good,
        Title: "Task 1 Completed",
    })
    .AddAttachment(samhook.Attachment{
        Color: samhook.Warning,
        Title: "Task 2 Needs Attention",
    })

samhook.Send(webhookURL, msg)
```

### Batch Adding Attachments

```go
attachments := []samhook.Attachment{
    {
        Color: samhook.Good,
        Title: "Service A Normal",
        Text:  "All checks passed",
    },
    {
        Color: samhook.Good,
        Title: "Service B Normal",
        Text:  "All checks passed",
    },
}

msg := samhook.Message{Text: "Health Check Report"}
    .AddAttachments(attachments)

samhook.Send(webhookURL, msg)
```

## Advanced Usage

### Attachment with Author Information

```go
attachment := samhook.Attachment{
    Color:      samhook.Good,
    AuthorName: "Deployment Bot",
    AuthorLink: "https://example.com/bot",
    AuthorIcon: "https://example.com/bot-icon.png",
    Title:      "Deployment Complete",
    TitleLink:  "https://example.com/deployment/123",
    Text:       "Application successfully deployed to production",
}

msg := samhook.Message{}
msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### Attachment with Images

```go
attachment := samhook.Attachment{
    Color:    samhook.Good,
    Title:    "Chart Report",
    Text:     "Weekly data analysis",
    ImageURL: "https://example.com/chart.png",
    ThumbURL: "https://example.com/thumb.png",
}

msg := samhook.Message{}
msg.AddAttachment(attachment)
samhook.Send(webhookURL, msg)
```

### Using SendReader

When you already have a JSON-formatted message:

```go
package main

import (
    "bytes"
    "github.com/circleyu/samhook"
)

func main() {
    webhookURL := "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
    
    // Assume you already have a JSON-formatted message
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

## Real-World Application Scenarios

### Deployment Notification

```go
func notifyDeployment(webhookURL, env, version string, success bool) error {
    msg := samhook.Message{
        Username: "Deployment System",
        IconEmoji: ":rocket:",
    }
    
    attachment := samhook.Attachment{
        Title: "Deployment Notification",
        Fields: []samhook.Field{
            {Title: "Environment", Value: env, Short: true},
            {Title: "Version", Value: version, Short: true},
        },
        Footer: "CI/CD System",
    }
    
    if success {
        attachment.Color = samhook.Good
        attachment.Text = "Deployment Successful"
    } else {
        attachment.Color = samhook.Danger
        attachment.Text = "Deployment Failed"
    }
    
    msg.AddAttachment(attachment)
    return samhook.Send(webhookURL, msg)
}
```

## Error Handling

### Basic Error Handling

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    log.Printf("Send failed: %v", err)
}
```

### Detailed Error Handling with WebhookError

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        // Check error type
        if webhookErr.IsNetworkError() {
            log.Printf("Network error: %v", webhookErr)
            // Can retry
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            log.Printf("API error (status code %d): %v", statusCode, webhookErr)
            
            // Handle specific status codes
            switch statusCode {
            case 401:
                log.Println("Authentication failed, please check webhook URL")
            case 429:
                log.Println("Rate limited, please retry later")
            case 500, 502, 503, 504:
                log.Println("Server error, can retry")
            }
        } else if webhookErr.IsSerializationError() {
            log.Printf("Serialization error: %v", webhookErr)
        }
        
        // Get detailed error message
        log.Println(webhookErr.DetailedMessage())
    } else {
        log.Printf("Unknown error: %v", err)
    }
}
```

## Advanced Configuration

### Using Custom Timeout

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "Important Notification"}

// Set 30 second timeout
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)
```

### Using Custom HTTP Client

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

### Using Context for Timeout and Cancellation

```go
import (
    "context"
    "time"
    "github.com/circleyu/samhook"
)

// Use Context timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := samhook.SendWithContext(ctx, webhookURL, msg)

// Use Context cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(5 * time.Second)
    cancel() // Cancel request after 5 seconds
}()

err := samhook.SendWithContext(ctx, webhookURL, msg)
```

## Retry Mechanism

### Basic Retry

```go
import "github.com/circleyu/samhook"

msg := samhook.Message{Text: "Important Notification"}

// Use default retry options (max 3 retries, 1 second interval)
opts := samhook.DefaultRetryOptions
err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### Custom Retry Configuration

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "Important Notification"}

opts := samhook.RetryOptions{
    MaxRetries: 5,
    Interval:   2 * time.Second,
    Backoff: &samhook.ExponentialBackoff{
        InitialInterval: 1 * time.Second,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,
        Jitter:          true, // Add random jitter to avoid thundering herd
    },
}

err := samhook.SendWithRetry(webhookURL, msg, opts)
```

### Retry Scenario Example

```go
// Handle temporary network failures
func sendWithRetry(webhookURL string, msg samhook.Message) error {
    opts := samhook.DefaultRetryOptions
    opts.MaxRetries = 3
    
    err := samhook.SendWithRetry(webhookURL, msg, opts)
    if err != nil {
        // Check if it's a retryable error
        if webhookErr, ok := err.(*samhook.WebhookError); ok {
            if webhookErr.IsNetworkError() {
                log.Printf("Network error, retried %d times: %v", opts.MaxRetries, err)
            } else if webhookErr.IsAPIError() && webhookErr.GetStatusCode() >= 500 {
                log.Printf("Server error, retried %d times: %v", opts.MaxRetries, err)
            }
        }
    }
    return err
}
```

### Monitoring Alert

```go
func sendAlert(webhookURL, alertType, message string, severity int) error {
    msg := samhook.Message{
        Username: "Monitoring System",
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
        Footer: "Monitoring System",
    }
    
    msg.AddAttachment(attachment)
    return samhook.Send(webhookURL, msg)
}
```

### Task Completion Notification

```go
func notifyTaskCompletion(webhookURL, taskName, duration string) error {
    msg := samhook.Message{
        Text:     "Task Completion Notification",
        Username: "Task System",
    }
    
    attachment := samhook.Attachment{
        Color: samhook.Good,
        Title: "Task Completed",
        Fields: []samhook.Field{
            {Title: "Task Name", Value: taskName, Short: false},
            {Title: "Execution Time", Value: duration, Short: true},
        },
    }
    
    msg.AddAttachment(attachment)
    return samhook.Send(webhookURL, msg)
}
```

## URL Validation

### Validating Webhook URL

```go
import "github.com/circleyu/samhook"

err := samhook.ValidateWebhookURL(webhookURL)
if err != nil {
    log.Fatalf("Invalid webhook URL: %v", err)
}

// Proceed with sending
msg := samhook.Message{Text: "Hello"}
samhook.Send(webhookURL, msg)
```

### Validating Before Sending

```go
func sendWithValidation(webhookURL string, msg samhook.Message) error {
    if err := samhook.ValidateWebhookURL(webhookURL); err != nil {
        return fmt.Errorf("webhook URL validation failed: %w", err)
    }
    return samhook.Send(webhookURL, msg)
}
```

## Logging

### Basic Logging

```go
import (
    "os"
    "github.com/circleyu/samhook"
)

// Enable logging to standard output
samhook.SetLoggerWriter(os.Stdout)

msg := samhook.Message{Text: "Hello"}
samhook.Send(webhookURL, msg)
// Output: [samhook] POST https://hooks.slack.com/... - success - duration: 123ms
```

### Custom Logger

```go
import (
    "log"
    "os"
    "time"
    "github.com/circleyu/samhook"
)

type MyLogger struct {
    logger *log.Logger
}

func (l *MyLogger) LogRequest(url string, method string, duration time.Duration, err error) {
    status := "OK"
    if err != nil {
        status = "ERROR"
    }
    l.logger.Printf("[%s] %s %s took %v", status, method, url, duration)
}

// Use custom logger
logger := &MyLogger{
    logger: log.New(os.Stderr, "[webhook] ", log.LstdFlags),
}
samhook.SetLogger(logger)

msg := samhook.Message{Text: "Hello"}
samhook.Send(webhookURL, msg)
```

### Retry with Custom Client Options

```go
import (
    "time"
    "github.com/circleyu/samhook"
)

msg := samhook.Message{Text: "Important Notification"}

opts := samhook.DefaultRetryOptions
opts.MaxRetries = 5

// Use retry with custom timeout
err := samhook.SendWithRetry(webhookURL, msg, opts, 
    samhook.WithTimeout(30*time.Second),
)
```
