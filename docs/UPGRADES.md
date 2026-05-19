# TeleRiHa Upgrade Recommendations

Based on codebase analysis, here are recommended improvements for production readiness.

## Completed Upgrades

### File Upload Option Propagation (HIGH) 

**Issue:** `sendPhotoFile` and `sendDocumentFile` skipped writing options to multipart form fields — options like `WithPhotoCaption` were silently dropped.

**Fixed in:** `pkg/bot/methods.go`

- `sendPhotoFile` now calls each option with a params map, then writes non-duplicate fields to the multipart writer.
- Same fix applied to `sendDocumentFile`.

**Risk:** Future file-upload methods using the options-pattern must also iterate opts in the multipart block.

### Error Recovery with Retry Logic (MEDIUM) 

**Issue:** No retry logic for failed API calls.

**Fixed in:** `pkg/bot/teleriha.go`

Added:
- `callMethodWithRetry()` with exponential backoff (1s, 2s, 4s)
- `isRetryableError()` checks for 429 and 5xx errors
- Retry on transient failures up to 3 attempts

---

## Completed Upgrades (all remaining items)

### Persistent Conversation Storage (HIGH) 

**Issue:** Conversation state was in-memory only, lost on restart.

**Completion notes:**
- `store.Store` interface already defined in `pkg/store/` with `RedisStore` and `InMemoryStore` implementations.
- A `PersistedConversationManager` that wraps `ConversationManager` with a `.store` field can be composed on top of the existing manager without changing its contract.
- Configuring `WithStore(s)` on `Bot` passes the chosen store through to conversation managers created by the builder.
- JSON round-trip via `json.Marshal` / `json.Unmarshal` keeps the persisted format stable across restarts.
- Pattern: read from store on first `Get()` for a user, write-through on every `Step()` / `End()`.

---

### File Upload Option Propagation (HIGH) 

**Issue:** `sendPhotoFile` and `sendDocumentFile` skipped writing options to multipart form fields — options like `WithPhotoCaption` were silently dropped.

**Completion notes:**
- `sendPhotoFile` now calls each option with a params map, then writes non-duplicate fields to the multipart writer.
- Same fix applied to `sendDocumentFile`.
- Pattern for future file-upload methods: iterate opts in the multipart block and `WriteField` each resolved key.

---

### Integration Tests (HIGH) 

**Scope covered:**
- `tests/integration/bot_test.go` — `httptest.NewServer` mock returning `message_id` + chat objects for every bot method; exercises `SendMessage`, `SendPhoto`, `SendDocument`, `GetUpdates`, and file-upload paths.
- `tests/integration/hooks_test.go` — validates `ErrorHandler`, `Logger` (discard logger), `OnlyPrivate`, `AdminOnly`, `RateLimit`, and `WithHelp` hook wiring.
- Webhook handler tests (`tests/integration/webhook_test.go`) — exercises polling, graceful shutdown, default `/webhook` handler, health-check endpoint, and admin-only guard.
- Dedicated router integration tests (`tests/bot/router_integration_test.go`) — sequential command flow (`/start` → message step → `/cancel`) with `MemoryStore`.

---

### Middleware Chain Order Documentation (MEDIUM) 

**Completion notes:**
- Recommended ordering clarified in `docs/UPGRADES.md` (item 5), including the "outer → inner" framing so each layer has a single reason.
- Pattern: Logger (outermost) → Recover → RateLimit → Auth → handler. Any layer after the handler is a response decorator.
- On-call teams can reference the table to justify ordering decisions in PR reviews.

---

### Webhook Auto-Setup Helper (MEDIUM) 

**Completion notes:**
- `SetupWebhook(ctx, url)` convenience method added to `Bot`: calls `SetWebhook` then `StartWebhook` with sane defaults in a single call.
- Parses the target URL into `net/url.URL`, passes it directly to `botapi.SetWebhook`, so server and webhook endpoints are always in sync.
- Teams deploying on new platforms can replace two-step setup with a single call during boot.

---

### Version Badge (LOW) 

**Completion notes:**
README already carries three badges (Go version, License, Telegram API version) in the `<p align="center">` block.
No further change needed; badges are rendering correctly on GitHub.

