# TeleRiHa Project Context

TeleRiHa is a high-performance, modular Telegram Bot framework for Go. It focuses on simplicity, speed, and extensibility through a robust plugin system and standard Go paradigms.

## Core Concepts

- **Bot**: The central orchestrator that manages the connection to Telegram.
- **Router**: Dispatches incoming updates to appropriate handlers based on matching rules.
- **Context**: A per-update object containing request data and helper methods for responses.
- **Plugin**: A modular component that extends bot functionality (e.g., rate limiting, i18n).
- **Store**: An interface for persistent or volatile state storage (e.g., Memory, Redis).
- **Handler**: A function that processes an update and returns an error.
- **Middleware**: Functions that wrap handlers to provide cross-cutting concerns.

## Technical Stack

- **Language**: Go 1.24+
- **Logging**: zerolog
- **CLI**: cobra
- **Storage**: Redis (go-redis), In-memory
- **Testing**: testify, httptest
