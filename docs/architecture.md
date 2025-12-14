# Architecture Documentation

English | [繁體中文](architecture_zh_TW.md)

## Project Architecture Overview

samhook is a lightweight Go library designed to provide a simple and easy-to-use API for sending Slack and Mattermost webhook messages.

## Design Principles

### 1. Lightweight Design

- **Minimal external dependencies**: Uses only high-performance JSON library
- **Minimal abstraction**: Direct mapping to Slack/Mattermost webhook API
- **Simple and easy to use**: Provides intuitive API interface

### 2. Type Safety

- **Strong typing**: All data structures have explicit type definitions
- **JSON tags**: Uses standard JSON tags for serialization
- **Optional fields**: Uses `omitempty` tags to serialize only fields with values

### 3. Flexibility

- **Multiple input methods**: Supports both struct and Reader input methods
- **Method chaining**: Supports method chaining for better developer experience
- **Extensible**: Easy to add new features

## Module Structure

### Core Modules

```
samhook/
├── message.go    # Data structure definitions
├── samhook.go   # Core functionality implementation
├── error.go     # Error type definitions
├── client.go    # HTTP client configuration
└── retry.go     # Retry mechanism
```

### message.go

Defines three core data structures:

1. **Message** - Message body
   - Contains basic message attributes (text, username, icon, etc.)
   - Contains attachment list

2. **Attachment** - Attachment structure
   - Supports rich formatting options
   - Supports fields, images, author information, etc.

3. **Field** - Field structure
   - Used to display structured data in attachments
   - Supports short fields for side-by-side display

### samhook.go

Implements core functionality:

1. **Send function**
   - Receives Message structure
   - Serializes to JSON
   - Sends HTTP POST request
   - Checks HTTP status code and returns detailed errors

2. **SendReader function**
   - Receives io.Reader
   - Directly sends HTTP POST request
   - Suitable for scenarios with existing JSON data
   - Checks HTTP status code and returns detailed errors

3. **AddAttachment method**
   - Adds a single attachment to Message
   - Returns *Message to support method chaining

4. **AddAttachments method**
   - Adds multiple attachments to Message
   - Returns *Message to support method chaining

### error.go

Implements error handling functionality:

1. **WebhookError type**
   - Custom error type providing detailed error information
   - Supports error classification (network errors, serialization errors, API errors)
   - Provides error codes and detailed messages

2. **Error constructor functions**
   - `NewNetworkError()` - Creates network error
   - `NewSerializationError()` - Creates serialization error
   - `NewAPIError()` - Creates API error

### client.go

Implements HTTP client configuration functionality:

1. **SendWithOptions function**
   - Supports custom HTTP client configuration
   - Uses options pattern
   - Supports timeout configuration

2. **SendWithContext function**
   - Supports Context timeout and cancellation
   - Uses `http.NewRequestWithContext()`

3. **ClientOption type**
   - `WithTimeout()` - Sets timeout
   - `WithClient()` - Uses custom client

### retry.go

Implements retry mechanism:

1. **SendWithRetry function**
   - Supports automatic retry for failed requests
   - Configurable retry count and interval
   - Intelligently determines which errors can be retried

2. **ExponentialBackoff type**
   - Implements exponential backoff strategy
   - Supports random jitter
   - Configurable initial interval and maximum interval

## Data Flow

### Standard Send Flow

```
User Code
    ↓
Create Message structure
    ↓
(Optional) Add attachments
    ↓
Call Send() function
    ↓
JSON serialization
    ↓
HTTP POST request
    ↓
Webhook URL
```

### SendReader Flow

```
User Code
    ↓
Prepare JSON data (io.Reader)
    ↓
Call SendReader() function
    ↓
HTTP POST request (directly uses Reader)
    ↓
Webhook URL
```

## Dependencies

### Standard Library Dependencies

- `bytes` - For creating request body
- `context` - Context support (for timeout and cancellation)
- `fmt` - String formatting
- `io` - Reader interface
- `math` - Mathematical operations (for exponential backoff)
- `math/rand` - Random number generation (for jitter)
- `net` - Network error classification
- `net/http` - HTTP client
- `net/url` - URL error handling
- `strings` - String operations
- `time` - Time and timeout control

### External Dependencies

- `github.com/bytedance/sonic` - High-performance JSON serialization/deserialization library

The project uses a high-performance JSON library, ensuring:
- Fast compilation
- High execution performance
- Optimized memory usage
- Simple dependency management

## Extensibility Design

### Implemented Features

1. ✅ **Configurable HTTP client**: Supports custom timeout, custom client (`SendWithOptions`)
2. ✅ **Response handling**: Checks HTTP status code and returns detailed errors
3. ✅ **Retry mechanism**: Supports automatic retry with exponential backoff strategy
4. ✅ **Context support**: Supports timeout and cancellation control
5. ✅ **Detailed error handling**: Custom error type providing error classification and detailed information

### Future Extension Directions

1. **Batch sending**: Support sending multiple messages at once
2. **Async sending**: Provide helper functions to simplify async sending
3. **Rate limiting handling**: Automatically handle 429 errors and request queuing

## Compatibility

### Slack Webhook API

samhook is fully compatible with Slack Incoming Webhooks API, supporting:
- Basic message format
- Attachment format
- Field format
- All standard fields

### Mattermost Webhook API

Mattermost's webhook API is highly compatible with Slack, so samhook can also be used directly with Mattermost.

## Error Handling Strategy

### Current Implementation

- ✅ All functions return `error` (may be `*WebhookError`)
- ✅ Error classification: network errors, serialization errors, API errors
- ✅ Detailed error information: includes status code, response body, error code, etc.
- ✅ Error classification methods: `IsNetworkError()`, `IsAPIError()`, `IsSerializationError()`
- ✅ Detailed error messages: `DetailedMessage()` provides multi-line format

### Error Types

1. **Network errors** (`ErrorTypeNetwork`): HTTP request failures, connection errors, etc.
2. **Serialization errors** (`ErrorTypeSerialization`): JSON serialization failures
3. **API errors** (`ErrorTypeAPI`): HTTP status codes other than 200
4. **Unknown errors** (`ErrorTypeUnknown`): Errors that cannot be classified

### Error Handling Flow

```
Send Request
    ↓
Error occurred?
    ↓ Yes
Create WebhookError
    ↓
Classify error type
    ↓
Return WebhookError
    ↓
Users can:
- Check error type
- Get status code
- Get detailed message
- Implement intelligent handling
```

## Performance Considerations

### Current Implementation

- Uses `http.DefaultClient` (shared connection pool)
- Synchronous sending (blocks until completion)
- High-performance serialization (uses sonic JSON library)

### Performance Characteristics

- **Low latency**: Direct HTTP requests, no additional overhead
- **Low memory**: Minimizes memory allocation
- **High throughput**: Can be used concurrently (each goroutine is independent)

## Security Considerations

### Current Implementation

- Does not validate webhook URL
- Does not encrypt transmission (relies on HTTPS)
- Does not validate response

### Recommendations

1. **URL validation**: Validate webhook URL format
2. **HTTPS enforcement**: Enforce HTTPS in production environments
3. **Response validation**: Check HTTP status code

## Testing Strategy

### Implemented Tests

1. ✅ **Unit tests**: 
   - Data structure serialization tests (`message_test.go`)
   - Error type tests (`samhook_test.go`)
   - Client configuration tests (`client_test.go`)
   - Retry mechanism tests (`retry_test.go`)

2. ✅ **Integration tests**: 
   - Uses `httptest` package to create mock HTTP servers
   - Tests various HTTP status code scenarios
   - Tests error handling logic

3. **End-to-end tests**: 
   - Optional integration tests (requires actual webhook URL)
   - Marked as integration tests using build tags

### Test Coverage

Current test coverage is approximately 62%, covering:
- Serialization of all data structures
- Core sending functionality
- Error handling logic
- Client configuration
- Retry mechanism

## Maintainability

### Code Organization

- **Clear module separation**: Data structures separated from functionality
- **Consistent naming**: Follows Go naming conventions
- **Complete comments**: All public APIs have comments

### Maintainability Features

- **Simple code structure**: Easy to understand and modify
- **Minimized complexity**: Avoids over-engineering
- **Standardized format**: Uses Go standard format
