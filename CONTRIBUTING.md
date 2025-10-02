# Contributing to Wooak

Thank you for your interest in contributing to Wooak! This document provides guidelines and information for contributors, especially those participating in Hacktoberfest.

## 🎃 Hacktoberfest Participation

Wooak is participating in [Hacktoberfest 2025](https://hacktoberfest.com/participation/)! We welcome contributions from developers of all skill levels.

### Hacktoberfest Guidelines

- **Quality over Quantity**: We value meaningful contributions that improve the project
- **Follow our standards**: All contributions must meet our quality and testing requirements
- **Respect the community**: Follow our Code of Conduct and be respectful to all contributors
- **Valid contributions**: Ensure your PRs are not spam and provide real value

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Git** - [Download Git](https://git-scm.com/downloads)
- **Make** - Usually pre-installed on Unix systems
- **Ollama** (optional) - For AI features testing

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/aryasoni98/wooak.git
   cd wooak
   ```

2. **Setup Development Environment**
   ```bash
   make dev-setup    # Install tools and dependencies
   make deps         # Download Go dependencies
   ```

3. **Verify Installation**
   ```bash
   make test         # Run tests
   make build        # Build the project
   ```

## 🛠️ Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### 2. Make Your Changes

- Write clean, readable code
- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed

### 3. Quality Checks

Before submitting, ensure your code passes all quality checks:

```bash
make quality        # Run all quality checks
make test           # Run unit tests
make coverage       # Generate coverage report
make lint           # Run linter
```

### 4. Commit Your Changes

Use semantic commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```bash
git add .
git commit -m "feat(ui): add new keyboard shortcut for search"
```

**Commit Message Format:**
- `feat(scope): description` - New features
- `fix(scope): description` - Bug fixes
- `docs: description` - Documentation
- `test(scope): description` - Tests
- `refactor(scope): description` - Code refactoring
- `chore: description` - Maintenance tasks
- `upgrade: description` - Dependency upgrades
- `security: description` - Security improvements
- `perf: description` - Performance improvements

**Available Scopes:**
- `ui` - User interface changes
- `ai` - AI-related features
- `security` - Security features
- `config` - Configuration handling
- `parser` - SSH config parsing
- `core` - Core business logic
- `deps` - Dependencies
- `workflow` - GitHub Actions
- `makefile` - Build system

### 5. Pull Request Title Format

**Important:** Your PR title must follow the same conventional commit format as your commit messages:

✅ **Good PR Titles:**
- `feat: add dark mode support to TUI`
- `fix(ui): resolve SSH connection timeout`
- `docs: update installation instructions`
- `upgrade: bump Go version to 1.21`
- `chore: update dependencies`

❌ **Bad PR Titles:**
- `Upgrade Codebase` (missing type prefix)
- `Fix bug` (too vague)
- `New feature` (missing colon)

**For Hacktoberfest PRs:** Use `[Hacktoberfest]: <type>: <description>`
- `[Hacktoberfest]: feat: add macOS support for Ollama`
- `[Hacktoberfest]: fix: resolve connection timeout on Windows`

**Need help?** Check our [Conventional Commits Guide](.github/CONVENTIONAL_COMMITS.md) for detailed examples.

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:
- Clear title following conventional commit format
- Reference any related issues
- Include screenshots for UI changes
- Ensure all CI checks pass

## 📋 Contribution Types

### 🐛 Bug Fixes

- Reproduce the bug with a test case
- Fix the issue with minimal changes
- Add tests to prevent regression
- Update documentation if needed

### ✨ New Features

- Discuss large features in an issue first
- Follow the existing architecture patterns
- Add comprehensive tests
- Update documentation and help text
- Consider backward compatibility

### 📚 Documentation

- Fix typos and improve clarity
- Add examples and use cases
- Update API documentation
- Improve README sections

### 🧪 Tests

- Add unit tests for new functionality
- Improve test coverage
- Add integration tests
- Add performance benchmarks

### 🎨 UI/UX Improvements

- Follow the existing design patterns
- Ensure keyboard navigation works
- Test on different terminal sizes
- Consider accessibility

## 🏗️ Project Structure

```
wooak/
├── cmd/                    # Application entry point
├── internal/               # Private application code
│   ├── adapters/          # External interface adapters
│   │   ├── data/          # Data layer adapters
│   │   └── ui/            # User interface adapters
│   ├── core/              # Business logic
│   │   ├── domain/        # Domain models
│   │   ├── ports/         # Interface definitions
│   │   └── services/      # Business services
│   └── logger/            # Logging utilities
├── .github/               # GitHub configuration
│   └── ISSUE_TEMPLATE/    # Issue templates
├── docs/                  # Documentation
└── tests/                 # Test files
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific test packages
go test ./internal/core/services/...

# Run benchmarks
make benchmark
```

### Writing Tests

- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Test both success and error cases
- Aim for high test coverage

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "success case",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:     "error case",
            input:    invalidInput,
            expected: nil,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("FunctionName() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## 📝 Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `golint`
- Write clear, self-documenting code
- Use meaningful variable and function names

### Code Organization

- Keep functions small and focused
- Use interfaces for testability
- Follow the existing architectural patterns
- Add comments for complex logic

### Error Handling

```go
// Good
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// Avoid
result, _ := someFunction() // Don't ignore errors
```

## 🔍 Code Review Process

### For Contributors

1. **Self-Review**: Review your own code before submitting
2. **Address Feedback**: Respond to review comments promptly
3. **Update Tests**: Add tests for any new functionality
4. **Update Documentation**: Keep docs in sync with code changes

### For Reviewers

1. **Be Constructive**: Provide helpful feedback
2. **Be Respectful**: Remember the human behind the code
3. **Be Thorough**: Check for bugs, performance issues, and style
4. **Be Prompt**: Respond to PRs in a timely manner

## 🎯 Good First Issues

Looking for your first contribution? Check out these areas:

- **Documentation**: Fix typos, improve examples
- **Tests**: Add test coverage for existing functions
- **UI Polish**: Improve user experience
- **Bug Fixes**: Fix issues labeled as "good first issue"
- **Performance**: Optimize existing code

## 🚫 What Not to Contribute

To maintain quality and avoid spam:

- **Don't** submit PRs that only change whitespace or formatting
- **Don't** submit PRs that only update version numbers
- **Don't** submit PRs without tests for new functionality
- **Don't** submit PRs that break existing functionality
- **Don't** submit PRs with unclear or missing descriptions

## 🆘 Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Code Review**: Ask questions in PR comments
- **Documentation**: Check the README and inline code comments

## 📜 License

By contributing to Wooak, you agree that your contributions will be licensed under the Apache License 2.0.

## 🙏 Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Project documentation

Thank you for contributing to Wooak! 🎉
