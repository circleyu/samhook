# API Documentation

English | [繁體中文](api_zh_TW.md)

## Overview

This document provides a complete API reference for the samhook library.

## Constants

### Color Constants

Predefined color constants for setting attachment colors:

```go
const Warning string = "#FFBB00"  // Warning color (yellow)
const Danger string = "#FF0000"   // Danger color (red)
const Good string = "#00FF00"      // Normal information color (green)
```

## Data Structures

### Message

Message structure representing the webhook message to be sent.

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

#### Field Descriptions

- `Parse` - Parse mode (optional)
- `Username` - Sender username (optional)
- `IconURL` - Icon URL (optional)
- `IconEmoji` - Icon emoji (optional)
- `Channel` - Target channel (optional)
- `Text` - Message text content (optional)
- `Attachments` - List of attachments (optional)

### Attachment

Attachment structure for enriching message content.

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

#### Field Descriptions

- `Fallback` - Fallback text (for clients that don't support attachments)
- `Color` - Left-side color bar (can use predefined constants)
- `Pretext` - Text before attachment
- `AuthorName` - Author name
- `AuthorLink` - Author link
- `AuthorIcon` - Author icon URL
- `Title` - Attachment title
- `TitleLink` - Title link
- `Text` - Attachment text content
- `ImageURL` - Image URL
- `Fields` - List of fields
- `Footer` - Footer text
- `FooterIcon` - Footer icon URL
- `ThumbURL` - Thumbnail URL

### Field

Field structure for displaying key-value information in attachments.

```go
type Field struct {
    Title string `json:"title,omitempty"`
    Value string `json:"value,omitempty"`
    Short bool   `json:"short,omitempty"`
}
```

#### Field Descriptions

- `Title` - Field title
- `Value` - Field value
- `Short` - Whether it's a short field (for side-by-side display)

## Functions

### Send

Sends a message to the specified webhook URL.

```go
func Send(url string, msg Message) error
```

#### Parameters

- `url` - Webhook URL (string)
- `msg` - Message structure to send

#### Return Value

- `error` - Returns an error if sending fails, otherwise returns nil

#### Behavior

1. Serializes the message structure to JSON
2. Creates an HTTP POST request
3. Sets Content-Type to application/json
4. Sends the request and closes the response body

### SendReader

Sends a message from io.Reader to the specified webhook URL.

```go
func SendReader(url string, r io.Reader) error
```

#### Parameters

- `url` - Webhook URL (string)
- `r` - Reader containing JSON message

#### Return Value

- `error` - Returns an error if sending fails, otherwise returns nil

#### Use Cases

Use this function when you already have JSON-formatted message data (e.g., read from a file or obtained from another source).

## Methods

### AddAttachment

Adds a single attachment to the message. Returns a message pointer to support method chaining.

```go
func (m *Message) AddAttachment(attachment Attachment) *Message
```

#### Parameters

- `attachment` - Attachment to add

#### Return Value

- `*Message` - Message pointer (for method chaining)

#### Example

```go
msg := samhook.Message{Text: "Alert!"}
msg.AddAttachment(samhook.Attachment{
    Color: samhook.Danger,
    Title: "Error",
    Text:  "Something went wrong",
})
```

### AddAttachments

Adds multiple attachments to the message. Returns a message pointer to support method chaining.

```go
func (m *Message) AddAttachments(attachments []Attachment) *Message
```

#### Parameters

- `attachments` - List of attachments to add

#### Return Value

- `*Message` - Message pointer (for method chaining)

#### Example

```go
msg := samhook.Message{Text: "Multiple alerts"}
attachments := []samhook.Attachment{
    {Color: samhook.Warning, Title: "Warning 1"},
    {Color: samhook.Danger, Title: "Error 1"},
}
msg.AddAttachments(attachments)
```

### SendWithOptions

Sends a message using the options pattern, supporting custom HTTP client configuration.

```go
func SendWithOptions(url string, msg Message, opts ...ClientOption) error
```

#### Parameters

- `url` - Webhook URL (string)
- `msg` - Message structure to send
- `opts` - Client options (variadic parameters)

#### Return Value

- `error` - Returns an error if sending fails, otherwise returns nil

#### Example

```go
err := samhook.SendWithOptions(webhookURL, msg,
    samhook.WithTimeout(30 * time.Second),
)
```

### SendWithContext

Sends a message using Context, supporting timeout and cancellation.

```go
func SendWithContext(ctx context.Context, url string, msg Message, opts ...ClientOption) error
```

#### Parameters

- `ctx` - Context for timeout and cancellation control
- `url` - Webhook URL (string)
- `msg` - Message structure to send
- `opts` - Client options (variadic parameters)

#### Return Value

- `error` - Returns an error if sending fails, otherwise returns nil

#### Example

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := samhook.SendWithContext(ctx, webhookURL, msg)
```

### SendWithRetry

Sends a message with retry mechanism, supporting custom client configuration.

```go
func SendWithRetry(url string, msg Message, opts RetryOptions, clientOpts ...ClientOption) error
```

#### Parameters

- `url` - Webhook URL (string)
- `msg` - Message structure to send
- `opts` - Retry options
- `clientOpts` - Optional client configuration options (variadic)

#### Return Value

- `error` - Returns an error if all retries fail, otherwise returns nil

#### Retry Conditions

- Network errors are automatically retried
- 5xx server errors are automatically retried
- 429 rate limit errors are automatically retried
- 4xx client errors are not retried

#### Example

```go
opts := samhook.DefaultRetryOptions
opts.MaxRetries = 5
err := samhook.SendWithRetry(webhookURL, msg, opts)

// With custom timeout
err := samhook.SendWithRetry(webhookURL, msg, opts, samhook.WithTimeout(30*time.Second))
```

## Option Types

### ClientOption

HTTP client option function type.

```go
type ClientOption func(*http.Client)
```

#### Available Options

- `WithTimeout(timeout time.Duration)` - Sets timeout duration
- `WithClient(client *http.Client)` - Uses a custom HTTP client

### RetryOptions

Retry options structure.

```go
type RetryOptions struct {
    MaxRetries int
    Interval   time.Duration
    Backoff    *ExponentialBackoff
}
```

#### Field Descriptions

- `MaxRetries` - Maximum number of retries
- `Interval` - Fixed retry interval (if Backoff is not set)
- `Backoff` - Exponential backoff configuration (optional)

### ExponentialBackoff

Exponential backoff configuration.

```go
type ExponentialBackoff struct {
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
    Jitter          bool
}
```

#### Field Descriptions

- `InitialInterval` - Initial retry interval
- `MaxInterval` - Maximum retry interval
- `Multiplier` - Backoff multiplier (typically 2.0)
- `Jitter` - Whether to add random jitter

## Error Types

### WebhookError

Custom error type providing detailed error information.

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

#### Methods

- `Error() string` - Implements the error interface
- `Unwrap() error` - Returns the original error
- `IsNetworkError() bool` - Checks if it's a network error
- `IsSerializationError() bool` - Checks if it's a serialization error
- `IsAPIError() bool` - Checks if it's an API error
- `GetStatusCode() int` - Returns HTTP status code
- `GetResponseBody() string` - Returns API response body
- `GetErrorCode() string` - Returns error code
- `DetailedMessage() string` - Returns detailed multi-line error message

#### Error Type Constants

```go
const (
    ErrorTypeNetwork       = "network"
    ErrorTypeSerialization = "serialization"
    ErrorTypeAPI           = "api"
    ErrorTypeUnknown       = "unknown"
)
```

#### Error Code Constants

```go
const (
    ErrorCodeNetworkTimeout     = "NETWORK_TIMEOUT"
    ErrorCodeNetworkConnection = "NETWORK_CONNECTION"
    ErrorCodeNetworkDNS         = "NETWORK_DNS"
    ErrorCodeSerializationJSON = "SERIALIZATION_JSON"
    ErrorCodeAPIUnauthorized   = "API_UNAUTHORIZED"
    ErrorCodeAPIForbidden      = "API_FORBIDDEN"
    ErrorCodeAPINotFound       = "API_NOT_FOUND"
    ErrorCodeAPIRateLimit      = "API_RATE_LIMIT"
    ErrorCodeAPIServerError    = "API_SERVER_ERROR"
)
```

**Note**: Network errors are now automatically classified into more specific types:
- `NETWORK_TIMEOUT` - Request timeout
- `NETWORK_DNS` - DNS resolution failure
- `NETWORK_CONNECTION` - Connection failure

#### Error Constructor Functions

- `NewNetworkError(url string, err error) *WebhookError` - Creates a network error
- `NewSerializationError(err error) *WebhookError` - Creates a serialization error
- `NewAPIError(url string, statusCode int, responseBody string) *WebhookError` - Creates an API error

## Error Handling

All functions return `error` when an error occurs. Errors returned from `Send()` and `SendReader()` may be of type `*WebhookError`, providing more detailed error information.

### Common Error Scenarios

- **Network errors**: Connection failures, timeouts, etc.
- **Serialization errors**: JSON serialization failures
- **API errors**: HTTP status codes other than 200 (e.g., 401, 429, 500, etc.)

### Error Handling Example

```go
err := samhook.Send(webhookURL, msg)
if err != nil {
    if webhookErr, ok := err.(*samhook.WebhookError); ok {
        if webhookErr.IsNetworkError() {
            // Handle network error
        } else if webhookErr.IsAPIError() {
            statusCode := webhookErr.GetStatusCode()
            // Handle based on status code
        }
    }
}
```

It is recommended to handle these errors appropriately in production environments and implement corresponding retry or fallback strategies based on error types.

## Utility Functions

### ValidateWebhookURL

Validates webhook URL format and protocol.

```go
func ValidateWebhookURL(webhookURL string) error
```

#### Parameters

- `webhookURL` - Webhook URL to validate (string)

#### Return Value

- `error` - Returns an error if validation fails, otherwise returns nil

#### Validation Rules

- URL must not be empty
- URL must be a valid URL format
- URL scheme must be `http` or `https`
- URL must have a host

#### Example

```go
err := samhook.ValidateWebhookURL(webhookURL)
if err != nil {
    log.Fatalf("Invalid webhook URL: %v", err)
}
```

## Logging

### Logger Interface

Optional logging interface for request tracking.

```go
type Logger interface {
    LogRequest(url string, method string, duration time.Duration, err error)
}
```

### SetLogger

Sets a package-level logger for request logging.

```go
func SetLogger(logger Logger)
```

### SetLoggerWriter

Sets a package-level logger using an `io.Writer`.

```go
func SetLoggerWriter(w io.Writer)
```

#### Example

```go
import (
    "os"
    "github.com/circleyu/samhook"
)

// Use standard output for logging
samhook.SetLoggerWriter(os.Stdout)

// Or use a custom logger
customLogger := &MyLogger{}
samhook.SetLogger(customLogger)
```
