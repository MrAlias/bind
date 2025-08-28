# Contributing

Thank you for your interest in contributing to the `bind` package!
We welcome contributions from the community.

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Run tests and ensure they pass
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.23 or later
- golangci-lint (for linting)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Code Quality

Before submitting a pull request, please ensure:

1. All tests pass
2. Code is properly formatted: `go fmt ./...`
3. Code passes linting: `golangci-lint run`
4. New code includes appropriate tests
5. Public APIs include documentation comments

## Pull Request Guidelines

- Keep pull requests focused on a single change
- Include tests for new functionality
- Update documentation as needed
- Write clear commit messages
- Reference any related issues

## Code Style

- Follow standard Go formatting and conventions
- Use meaningful variable and function names
- Include documentation comments for all public APIs
- Write tests that are clear and maintainable

## Questions?

If you have questions about contributing, please open an issue or start a discussion.

## License

By contributing to this project, you agree that your contributions will be licensed under the Apache 2.0 License.
