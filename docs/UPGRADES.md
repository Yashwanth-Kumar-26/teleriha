# TeleRiHa Upgrade Recommendations

Based on codebase analysis, here are recommended improvements for production readiness.

## High Priority

### 1. Persistent Conversation Storage

**Issue:** Conversation state is in-memory only, lost on restart.

**Current:**
```go
// conversation.go - in-memory map
conversations map[int64]*Conversation
```

**Recommendation:**
```go
// Add Redis-backed conversation storage
type PersistedConversationManager struct {
    *ConversationManager
    store store.Store
}

func (p *PersistedConversationManager) Start(userID, chatID int64, id string) *Conversation {
    conv := p.ConversationManager.Start(userID, chatID, id)
    // Persist to Redis
    data, _ := json.Marshal(conv)
    p.store.Set("conv:"+key, data, 24*time.Hour)
    return conv
}
```

### 2. File Upload Option Propagation

**Issue:** In `sendPhotoFile`, options are not properly applied.

**Current (methods.go:199-203):**
```go
for _, opt := range opts {
    opt(map[string]interface{}{})
    // Options are not applied here as we're using multipart
}
```

**Recommendation:**
```go
func (b *Bot) sendPhotoFile(...) (*Message, error) {
    // Add all options to form fields
    for _, opt := range opts {
        params := make(map[string]interface{})
        opt(params)
        for k, v := range params {
            if err := writer.WriteField(k, fmt.Sprintf("%v", v)); err != nil {
                return nil, err
            }
        }
    }
}
```

### 3. Add Integration Tests

**Missing:**
- Real Telegram API mock for testing
- File upload integration tests
- Webhook stress tests
- Redis integration tests

**Recommendation:**
```go
// tests/integration/bot_test.go
func TestBotIntegration(t *testing.T) {
    // Mock Telegram server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Return mock Update response
    }))
    defer mockServer.Close()

    b := bot.New("test-token", bot.WithBaseURL(mockServer.URL))
    // Test methods
}
```

---

## Medium Priority

### 4. Error Recovery Patterns

**Issue:** No retry logic for failed API calls.

**Recommendation:**
```go
func (b *Bot) callMethodWithRetry(method string, params map[string]interface{}, result interface{}) ([]byte, error) {
    var lastErr error
    for attempt := 0; attempt < 3; attempt++ {
        body, err := b.callMethod(method, params, result)
        if err == nil {
            return body, nil
        }
        // Check if retryable
        if !isRetryableError(err) {
            return nil, err
        }
        lastErr = err
        time.Sleep(time.Duration(attempt+1) * time.Second)
    }
    return nil, lastErr
}

func isRetryableError(err error) bool {
    // 429 Too Many Requests, 5xx errors
    // ...
}
```

### 5. Document Middleware Chain Order

**Issue:** Middleware order not clearly documented.

**Recommendation:** Add examples showing correct ordering:
```go
// CORRECT ORDER (outer to inner):
// 1. Logger (outermost - logs first, last)
// 2. Recover (catches panics after logging)
// 3. RateLimit (rejects before processing)
// 4. Auth (validates after rate limiting)

b.Router.Use(bot.Logger(log.Logger))
b.Router.Use(bot.Recover(log.Logger))
b.Router.Use(bot.RateLimitMiddleware(limiter))
b.Router.Use(bot.OnlyPrivate())
```

### 6. Add Webhook Auto-Setup Helper

**Issue:** No convenience method for common webhook setup.

**Recommendation:**
```go
func (b *Bot) SetupWebhook(ctx context.Context, url string) error {
    // Set webhook
    if err := b.SetWebhook(url, 100, nil); err != nil {
        return err
    }

    // Start server
    return b.StartWebhook("/webhook", 8080)
}
```

---

## Low Priority

### 7. Add CONTRIBUTING.md

**Missing:** Contribution guidelines for open source.

### 8. Add LICENSE

**Missing:** MIT license file.

### 9. Add Version Badge

**Current README:** No version badge.

**Recommendation:**
```markdown
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8)](https://golang.org)
[![Version](https://img.shields.io/github/v/release/Yashwanth-Kumar-26/teleriha)](https://github.com/Yashwanth-Kumar-26/teleriha/releases)
```

### 10. Performance Benchmarking

**Missing:** No performance benchmarks.

**Recommendation:**
```go
// tests/benchmarks/benchmark_test.go
func BenchmarkSendMessage(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bot.SendMessage(chatID, "test")
    }
}
```

---

## Summary

| Priority | Issue | Effort |
|----------|-------|--------|
| High | Persistent conversation storage | Medium |
| High | File upload option propagation | Low |
| High | Integration tests | High |
| Medium | Error recovery patterns | Medium |
| Medium | Middleware documentation | Low |
| Medium | Webhook helper | Low |
| Low | CONTRIBUTING.md | Low |
| Low | LICENSE | Low |
| Low | Version badge | Low |
| Low | Benchmarks | Medium |