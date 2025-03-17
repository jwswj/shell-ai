# Contributing to Shell-AI

Thank you for considering contributing to Shell-AI! This document provides guidelines and instructions for contributing to the project.

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/shell-ai.git`
3. Change to the project directory: `cd shell-ai`
4. Install dependencies: `go mod tidy`
5. Build the project: `make build`
6. Run the binary locally: `./bin/shai <prompt>`

For development, you can also install the binary directly to your GOPATH:

```bash
go install ./cmd/shai
```

## Running Tests

Run tests with:

```bash
make test
```

## Code Style

Please follow the standard Go code style guidelines:

- Format your code with `gofmt` or `go fmt`
- Use meaningful variable and function names
- Write comments for non-obvious code
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines

## Pull Request Process

1. Create a new branch for your feature or bugfix: `git checkout -b feature/your-feature-name`
2. Make your changes
3. Run tests to ensure they pass: `make test`
4. Commit your changes with a descriptive commit message
5. Push to your fork: `git push origin feature/your-feature-name`
6. Open a pull request against the main repository

## Reporting Issues

When reporting issues, please include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Any relevant logs or error messages
- Your environment (OS, Go version, etc.)

## Feature Requests

Feature requests are welcome! Please provide:

- A clear and descriptive title
- A detailed description of the proposed feature
- Any relevant examples or use cases

## License

By contributing to Shell-AI, you agree that your contributions will be licensed under the project's MIT License.

Thank you for contributing to Shell-AI!
