# Contributing to E-commerce Engine

Thank you for your interest in contributing to the E-commerce Engine! We welcome contributions from the community.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/masumrpg/ecommerce-engine.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit your changes: `git commit -m "Add your commit message"`
7. Push to your branch: `git push origin feature/your-feature-name`
8. Create a Pull Request

## Development Setup

### Prerequisites

- Go 1.25 or later
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/masumrpg/ecommerce-engine.git
cd ecommerce-engine

# Download dependencies
go mod download

# Run tests
go test ./...

# Run examples
go run examples/main.go
```

## Code Guidelines

### Go Code Style

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use `golint` to check for style issues
- Use `go vet` to check for suspicious constructs

### Testing

- Write unit tests for all new functionality
- Ensure all tests pass before submitting a PR
- Aim for high test coverage
- Use table-driven tests where appropriate

### Documentation

- Add godoc comments for all exported functions and types
- Update README.md if you add new features
- Include examples in your documentation

## Pull Request Process

1. Ensure your code follows the style guidelines
2. Add or update tests as needed
3. Update documentation
4. Ensure all tests pass
5. Create a clear and descriptive PR title
6. Provide a detailed description of your changes
7. Link any related issues

## Reporting Issues

When reporting issues, please include:

- Go version
- Operating system
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Any relevant code snippets or error messages

## Feature Requests

We welcome feature requests! Please:

- Check if the feature already exists
- Search existing issues to avoid duplicates
- Provide a clear description of the feature
- Explain the use case and benefits
- Consider providing a basic implementation proposal

## Code of Conduct

Please be respectful and professional in all interactions. We are committed to providing a welcoming and inclusive environment for all contributors.

## Questions?

If you have questions about contributing, feel free to:

- Open an issue with the "question" label
- Contact the maintainers

Thank you for contributing to the E-commerce Engine!