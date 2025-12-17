# Release Workflow Documentation

## Overview

The release workflow automates the complete release process for Wooak, including comprehensive testing across all supported platforms, building artifacts, and publishing GitHub releases. The workflow is triggered automatically when you push a commit with a release tag in the commit message.

## How to Trigger a Release

To trigger a release, include a release tag in your commit message using the format:

```
Releases: v0.0.3
```

### Examples

```bash
# Simple release commit
git commit -m "Releases: v0.0.3"

# Release with feature description
git commit -m "feat: Add new SSH key validation Releases: v0.0.3"

# Release with multiple changes
git commit -m "fix: Resolve connection timeout issue
feat: Add connection retry logic
Releases: v0.0.3"
```

### Version Format

The version must follow semantic versioning:
- Format: `vMAJOR.MINOR.PATCH`
- Examples: `v0.0.3`, `v1.2.0`, `v2.0.1`
- The `v` prefix is required

## Workflow Steps

### 1. Release Detection

The workflow first checks the commit message for a release tag pattern (`Releases: vX.X.X`). If found, it extracts the version and proceeds with the release process.

### 2. Comprehensive Testing

The workflow runs end-to-end tests across all supported platforms and architectures:

#### Linux Platforms
- **AMD64** (x86_64): Tested on Ubuntu with Go 1.21 and 1.23
- **ARM64**: Tested on Ubuntu
- **ARM**: Tested on Ubuntu
- **386** (i386): Tested on Ubuntu

#### Windows Platforms
- **AMD64**: Tested on Windows Latest

#### macOS Platforms
- **AMD64** (Intel): Tested on macOS 13
- **ARM64** (Apple Silicon): Tested on macOS 14

Each test job:
- Runs unit tests with race detection
- Attempts to run integration tests (if available)
- Builds the binary for the target platform
- Verifies the binary was created successfully

### 3. Code Quality Checks

Before building release artifacts, the workflow performs:
- `go vet` - Static analysis
- `go fmt` - Code formatting verification
- `staticcheck` - Advanced static analysis
- `golangci-lint` - Comprehensive linting

### 4. Security Scan

The workflow runs `govulncheck` to identify known vulnerabilities in dependencies.

### 5. Build Release Artifacts

After all tests pass, the workflow:
1. Creates a Git tag with the version from the commit message
2. Runs GoReleaser to build artifacts for all platforms
3. Creates a GitHub release with the artifacts

### 6. Release Verification

The workflow verifies that:
- The Git tag was created successfully
- The GitHub release was published
- All release artifacts are available

## Supported Platforms and Architectures

The workflow builds and tests the following combinations:

| OS | Architecture | Status |
|---|---|---|
| Linux | amd64 | ✅ |
| Linux | arm64 | ✅ |
| Linux | arm | ✅ |
| Linux | 386 | ✅ |
| Windows | amd64 | ✅ |
| macOS | amd64 | ✅ |
| macOS | arm64 | ✅ |

## Artifacts Created

GoReleaser creates the following artifacts for each release:

- **Linux**: `wooak_Linux_x86_64.tar.gz`, `wooak_Linux_arm64.tar.gz`, etc.
- **Windows**: `wooak_Windows_x86_64.zip`
- **macOS**: `wooak_Darwin_x86_64.tar.gz`, `wooak_Darwin_arm64.tar.gz`
- **Checksums**: `checksums.txt` for verification

## Workflow Configuration

The workflow is configured in `.github/workflows/release.yml` and uses:

- **GoReleaser**: For building and publishing releases (configured in `.goreleaser.yaml`)
- **GitHub Actions**: For CI/CD automation
- **Codecov**: For test coverage reporting (optional)

## Troubleshooting

### Workflow Not Triggering

1. **Check commit message format**: Ensure it contains `Releases: vX.X.X`
2. **Check branch**: Workflow only runs on `main` or `master` branches
3. **Check workflow file**: Verify `.github/workflows/release.yml` exists

### Tests Failing

1. **Check test logs**: Review the test output in GitHub Actions
2. **Run tests locally**: Use `make test` to reproduce issues
3. **Check Go version**: Ensure your local Go version matches the workflow

### Release Not Created

1. **Check tag creation**: Verify the tag was created in the repository
2. **Check GoReleaser logs**: Review the build-release job output
3. **Check permissions**: Ensure `GITHUB_TOKEN` has release permissions
4. **Check for existing tag**: If the tag already exists, the workflow will skip tag creation

### Version Extraction Failed

1. **Check commit message**: Ensure it follows the format `Releases: vX.X.X`
2. **Check version format**: Must be semantic versioning with `v` prefix
3. **Check workflow logs**: Review the detect-release job output

## Best Practices

1. **Version Management**: Always use semantic versioning
2. **Commit Messages**: Include release tags at the end of commit messages
3. **Testing**: Run tests locally before pushing release commits
4. **Changelog**: Update CHANGELOG.md before releasing
5. **Documentation**: Update documentation for significant releases

## Manual Release Process

If you need to create a release manually:

1. **Create and push a tag**:
   ```bash
   git tag -a v0.0.3 -m "Release v0.0.3"
   git push origin v0.0.3
   ```

2. **Run GoReleaser locally** (requires `GITHUB_TOKEN`):
   ```bash
   export GITHUB_TOKEN=your_token_here
   goreleaser release --clean
   ```

## Related Files

- `.github/workflows/release.yml` - Workflow definition
- `.goreleaser.yaml` - GoReleaser configuration
- `makefile` - Local build and test commands

## Support

For issues or questions about the release workflow:
- Open an issue on GitHub
- Check the [GitHub Actions documentation](https://docs.github.com/en/actions)
- Review [GoReleaser documentation](https://goreleaser.com/)
