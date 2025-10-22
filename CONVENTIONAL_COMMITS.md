# 📝 Conventional Commits Guide

This project uses [Conventional Commits](https://www.conventionalcommits.org/) to ensure consistent and meaningful commit messages and pull request titles.

## 🎯 **Why Use Conventional Commits?**

- **Automated changelog generation**
- **Semantic versioning support**
- **Better project organization**
- **Clear communication of changes**
- **Automated release processes**

## 📋 **PR Title Format**

All pull request titles must follow this format:

```
<type>[optional scope]: <description>
```

### **Examples:**

✅ **Good PR Titles:**
- `feat: add dark mode support to TUI`
- `fix: resolve SSH connection timeout issue`
- `docs: update installation instructions`
- `chore: update dependencies to latest versions`
- `upgrade: bump Go version to 1.21`
- `security: fix vulnerability in key validation`
- `perf: optimize server connection pooling`

❌ **Bad PR Titles:**
- `Upgrade Codebase` (missing type prefix)
- `Fix bug` (too vague)
- `Update stuff` (not descriptive)
- `New feature` (missing colon and scope)

## 🏷️ **Available Types**

| Type | Description | Example |
|------|-------------|---------|
| `feat` | A new feature | `feat: add AI-powered server search` |
| `fix` | A bug fix | `fix: resolve connection timeout` |
| `docs` | Documentation changes | `docs: update README with new features` |
| `style` | Code style changes (formatting, etc.) | `style: format code with gofumpt` |
| `refactor` | Code refactoring | `refactor: simplify SSH connection logic` |
| `perf` | Performance improvements | `perf: optimize server list rendering` |
| `test` | Adding or updating tests | `test: add unit tests for AI service` |
| `chore` | Maintenance tasks | `chore: update GitHub Actions` |
| `ci` | CI/CD changes | `ci: add automated testing workflow` |
| `build` | Build system changes | `build: update Makefile targets` |
| `upgrade` | Dependency upgrades | `upgrade: bump golangci-lint to v1.64.2` |
| `update` | General updates | `update: improve error handling` |
| `bump` | Version bumps | `bump: release v0.2.0` |
| `security` | Security-related changes | `security: fix SSH key validation` |
| `deps` | Dependency changes | `deps: add new Go modules` |

## 🎯 **Optional Scopes**

You can add a scope to provide more context:

```
<type>(<scope>): <description>
```

### **Available Scopes:**

| Scope | Description | Example |
|-------|-------------|---------|
| `ui` | User interface changes | `feat(ui): add dark mode toggle` |
| `cli` | Command-line interface | `fix(cli): resolve argument parsing` |
| `config` | Configuration changes | `feat(config): add SSH key validation` |
| `parser` | SSH config parsing | `fix(parser): handle malformed config` |
| `deps` | Dependencies | `upgrade(deps): update Go modules` |
| `security` | Security features | `feat(security): add audit logging` |
| `workflow` | GitHub Actions | `ci(workflow): add auto-labeling` |
| `makefile` | Build system | `chore(makefile): add new targets` |
| `docker` | Containerization | `feat(docker): add multi-stage build` |

## 🎃 **Hacktoberfest PRs**

For Hacktoberfest contributions, use this format:

```
[Hacktoberfest]: <type>[optional scope]: <description>
```

### **Examples:**
- `[Hacktoberfest]: feat: add macOS support for Ollama`
- `[Hacktoberfest]: fix: resolve connection timeout on Windows`
- `[Hacktoberfest]: docs: improve installation guide`
- `[Hacktoberfest]: feat(ui): add keyboard shortcuts help`

## 📝 **Description Guidelines**

### **Good Descriptions:**
- **Clear and concise**: `fix: resolve SSH connection timeout after 30 seconds`
- **Action-oriented**: `feat: add support for SSH key passphrase prompts`
- **Specific**: `docs: update macOS installation instructions for Homebrew`

### **Bad Descriptions:**
- **Too vague**: `fix: bug fix`
- **Too long**: `feat: add a really cool new feature that does amazing things and will revolutionize how users interact with SSH servers`
- **No action**: `SSH connection improvements`

## 🚫 **What to Avoid**

1. **Uppercase first letter**: Use lowercase for the description
   - ❌ `feat: Add new feature`
   - ✅ `feat: add new feature`

2. **Ending with period**: Don't end descriptions with a period
   - ❌ `fix: resolve connection issue.`
   - ✅ `fix: resolve connection issue`

3. **Imperative mood**: Use imperative mood (like git commit messages)
   - ❌ `feat: adds new feature`
   - ✅ `feat: add new feature`

## 🔧 **Automated Checks**

The project uses GitHub Actions to validate PR titles:

- **Semantic PR Action**: Validates conventional commit format
- **Auto-labeling**: Adds labels based on PR type
- **Changelog Generation**: Automatically generates changelog entries

## 🆘 **Need Help?**

If you're unsure about the format:

1. **Check existing PRs** for examples
2. **Use the PR template** which includes guidance
3. **Ask in discussions** if you need clarification
4. **Use the `skip-semantic-pr` label** to bypass validation (use sparingly)

## 📚 **Resources**

- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Angular Commit Guidelines](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#commit)
- [Commitizen](https://github.com/commitizen/cz-cli) - Tool to help write conventional commits

## 🎯 **Quick Reference**

```bash
# Common patterns:
feat: add new feature
fix: resolve bug
docs: update documentation
chore: maintenance task
upgrade: dependency upgrade
security: security improvement
perf: performance optimization
test: add or update tests
ci: CI/CD changes
build: build system changes

# With scope:
feat(ui): add dark mode
fix(cli): resolve argument parsing
docs(install): update macOS guide

# Hacktoberfest:
[Hacktoberfest]: feat: add new feature
[Hacktoberfest]: fix: resolve bug
```

Remember: **Good commit messages and PR titles make the project history more readable and help with automated tooling!** 🚀
