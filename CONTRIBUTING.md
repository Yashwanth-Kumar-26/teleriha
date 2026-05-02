# Contributing to TeleRiHa

Thank you for your interest in contributing to TeleRiHa!

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/teleriha.git
   cd TeleRiHa
   ```
3. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/Yashwanth-Kumar-26/teleriha.git
   ```

## Development Setup

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build the project
cd cmd/riha && go build -o riha .
```

## Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write code following Go conventions
   - Add tests for new functionality
   - Update documentation as needed

3. **Run tests**
   ```bash
   go test -v ./...
   go test -cover ./...
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

## Code Style

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Add godoc comments for exported functions
- Keep lines under 100 characters

## Pull Request Process

1. **Update documentation** if needed
2. **Ensure tests pass** locally
3. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```
4. **Open a Pull Request** against `main`

### PR Description

Include:
- What the change does
- Why it's needed
- How to test it
- Any breaking changes

## Reporting Issues

- Use GitHub Issues for bugs
- Include Go version
- Include minimal reproducible example
- Include error messages and stack traces

## Questions?

- Open a GitHub Discussion
- Check existing issues and discussions

---

Thank you for contributing to TeleRiHa!