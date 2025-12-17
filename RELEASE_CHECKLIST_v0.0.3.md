# Release Checklist for v0.0.3

## ✅ Pre-Release Checklist

### Files Updated
- [x] **CHANGELOG.md** - Added v0.0.3 release notes
- [x] **README.md** - Updated version badge to v0.0.3
- [x] **makefile** - Updated default VERSION to v0.0.3

### Homebrew Formula Verification

The Homebrew formula (`homebrew-tap/Formula/wooak.rb`) is configured correctly:

✅ **Current Configuration:**
- Points to v0.0.2 (will be auto-updated by GoReleaser)
- Uses `version` variable in build command
- Has proper test block
- License set to "MIT" (should match LICENSE file)

✅ **GoReleaser Configuration:**
- `brews` section configured correctly
- Repository: `aryasoni98/homebrew-tap`
- Homepage and description set

### What GoReleaser Will Do

When you create the v0.0.3 release, GoReleaser will:

1. **Build binaries** for all platforms (linux, windows, darwin)
2. **Create archives** (tar.gz for Unix, zip for Windows)
3. **Generate checksums** (checksums.txt)
4. **Create GitHub release** with tag v0.0.3
5. **Update Homebrew formula** automatically:
   - Update URL to `https://github.com/aryasoni98/wooak/archive/v0.0.3.tar.gz`
   - Calculate and update SHA256 hash
   - Commit changes to `homebrew-tap` repository

## 🚀 Release Steps

### 1. Verify Everything is Ready

```bash
# Ensure all tests pass
cd wooak
make test

# Ensure build works
make build

# Check for any uncommitted changes
git status
```

### 2. Create and Push Tag

```bash
# Create annotated tag
git tag -a v0.0.3 -m "Release v0.0.3: Enhanced error handling, improved test coverage, and code refactoring"

# Push tag to GitHub
git push origin v0.0.3
```

### 3. Create Release with GoReleaser

```bash
# If using GitHub Actions (recommended):
# Just push the tag - GitHub Actions will run GoReleaser automatically

# If running locally:
goreleaser release --clean
```

### 4. Verify Homebrew Formula Update

After release, check that the Homebrew formula was updated:

```bash
# Check the homebrew-tap repository
cd ../homebrew-tap
git log --oneline -1
cat Formula/wooak.rb | grep -A 2 "url"
```

Expected changes:
- URL should point to v0.0.3
- SHA256 should be updated

### 5. Test Homebrew Installation

```bash
# Test the formula locally
brew install --build-from-source ./Formula/wooak.rb

# Or test from tap (after pushing)
brew tap aryasoni98/homebrew-tap
brew install wooak
brew test wooak
```

## 🔍 Post-Release Verification

### GitHub Release
- [ ] Release created with tag v0.0.3
- [ ] Release notes populated from CHANGELOG
- [ ] All binaries attached
- [ ] Checksums file attached

### Homebrew Formula
- [ ] Formula updated in `homebrew-tap` repository
- [ ] URL points to v0.0.3
- [ ] SHA256 hash is correct
- [ ] Formula builds successfully
- [ ] Formula test passes

### Installation Test
```bash
# Test fresh installation
brew tap aryasoni98/homebrew-tap
brew install wooak
wooak --version  # Should show v0.0.3
```

## 📝 Notes

### Homebrew Formula Details

The formula uses:
- **Source-based installation**: Builds from source using Go
- **Version variable**: Uses `version` in ldflags for version info
- **Test block**: Verifies version command works

### Potential Issues

1. **License Mismatch**: Formula says "MIT" but LICENSE file is Apache 2.0
   - **Fix**: Update formula license to "Apache-2.0"

2. **Branch Name**: Formula uses "master" branch
   - **Verify**: Ensure your default branch is "master" or update to "main"

3. **SHA256 Verification**: GoReleaser calculates this automatically
   - **No action needed** - it's handled by GoReleaser

## ✅ Summary

All files are updated and ready for v0.0.3 release. The Homebrew formula will be automatically updated by GoReleaser when you create the GitHub release.

**Next Steps:**
1. Commit all changes
2. Create and push tag v0.0.3
3. Let GoReleaser handle the release and Homebrew update
4. Verify installation works
