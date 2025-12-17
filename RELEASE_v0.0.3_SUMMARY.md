# Release v0.0.3 - Summary and Verification

## ✅ Files Updated for v0.0.3

### Version Updates
- ✅ **CHANGELOG.md** - Added comprehensive v0.0.3 release notes
- ✅ **README.md** - Updated version badge from v0.0.2 to v0.0.3
- ✅ **makefile** - Updated default VERSION from v0.0.2 to v0.0.3

### Homebrew Formula Status
- ✅ **License Fixed** - Updated from "MIT" to "Apache-2.0" (matches LICENSE file)
- ✅ **Formula Structure** - Correctly configured for GoReleaser auto-update
- ✅ **Test Block** - Includes version verification test

## 🔍 Homebrew Formula Verification

### Current Formula (`homebrew-tap/Formula/wooak.rb`)

```ruby
class Wooak < Formula
  desc "A simple terminal UI for managing SSH connections"
  homepage "https://github.com/aryasoni98/wooak"
  url "https://github.com/aryasoni98/wooak/archive/v0.0.2.tar.gz"  # Will be auto-updated
  sha256 "ad4e841246e8c620673be7e381c4e44cac895213bed8b0a00924d73ce4873cc9"  # Will be auto-updated
  license "Apache-2.0"  # ✅ Fixed
  head "https://github.com/aryasoni98/wooak.git", branch: "master"  # ✅ Correct
  
  depends_on "go" => :build
  
  def install
    system "go", "build", *std_go_args(ldflags: "-X main.version=#{version}"), "./cmd/main.go"
  end
  
  test do
    assert_match "Wooak version", shell_output("#{bin}/wooak --version")
  end
end
```

### ✅ What Will Happen Automatically

When you create the v0.0.3 GitHub release, **GoReleaser will automatically**:

1. **Build binaries** for all platforms:
   - Linux (amd64, arm, arm64, 386)
   - Windows (amd64, arm, arm64, 386)
   - macOS/Darwin (amd64, arm64)

2. **Create release archives**:
   - `.tar.gz` for Unix systems
   - `.zip` for Windows

3. **Generate checksums** (`checksums.txt`)

4. **Create GitHub release** with:
   - Tag: `v0.0.3`
   - Release notes from CHANGELOG
   - All binaries attached
   - Checksums file attached

5. **Update Homebrew formula** in `homebrew-tap` repository:
   - URL: `https://github.com/aryasoni98/wooak/archive/v0.0.3.tar.gz`
   - SHA256: Automatically calculated from the archive
   - Commit and push changes to `homebrew-tap` repo

## 🧪 Testing Homebrew Formula

### Before Release (Current v0.0.2)
```bash
# Test current formula
cd homebrew-tap
brew install --build-from-source ./Formula/wooak.rb
wooak --version  # Should show v0.0.2
```

### After Release (v0.0.3)
```bash
# Test updated formula
brew tap aryasoni98/homebrew-tap
brew install wooak
wooak --version  # Should show v0.0.3
brew test wooak  # Should pass
```

## 📋 GoReleaser Configuration Verification

### ✅ Configuration Check

**`.goreleaser.yaml`** is correctly configured:

```yaml
brews:
  - repository:
      owner: aryasoni98
      name: homebrew-tap
    homepage: "https://github.com/aryasoni98/wooak"
    description: "A simple terminal UI for managing SSH connections."
```

**Status:**
- ✅ Repository owner matches: `aryasoni98`
- ✅ Repository name matches: `homebrew-tap`
- ✅ Homepage URL is correct
- ✅ Description is set

## 🚀 Release Process

### Step 1: Commit Changes
```bash
cd wooak
git add CHANGELOG.md README.md makefile
git commit -m "chore: prepare for v0.0.3 release"
git push origin master
```

### Step 2: Create and Push Tag
```bash
git tag -a v0.0.3 -m "Release v0.0.3: Enhanced error handling, improved test coverage, and code refactoring"
git push origin v0.0.3
```

### Step 3: GitHub Actions (Automatic)
If you have GitHub Actions configured, it will:
- Detect the new tag
- Run GoReleaser
- Create the release
- Update Homebrew formula

### Step 4: Manual Release (If Needed)
```bash
# Install GoReleaser if not installed
brew install goreleaser

# Create release
goreleaser release --clean
```

## ✅ Verification Checklist

### Pre-Release
- [x] CHANGELOG.md updated with v0.0.3 notes
- [x] README.md version badge updated
- [x] makefile version updated
- [x] Homebrew formula license fixed
- [x] All tests passing
- [x] GoReleaser config verified

### Post-Release
- [ ] GitHub release created with v0.0.3 tag
- [ ] All binaries attached to release
- [ ] Checksums file attached
- [ ] Homebrew formula auto-updated in `homebrew-tap`
- [ ] Formula URL points to v0.0.3
- [ ] Formula SHA256 is correct
- [ ] Installation test passes: `brew install aryasoni98/homebrew-tap/wooak`
- [ ] Version command works: `wooak --version` shows v0.0.3

## 🎯 Expected Results

### GitHub Release
- **Tag**: `v0.0.3`
- **Title**: `v0.0.3`
- **Description**: Auto-generated from CHANGELOG.md
- **Assets**: 
  - Multiple platform binaries (.tar.gz and .zip)
  - checksums.txt

### Homebrew Formula Update
The formula in `homebrew-tap/Formula/wooak.rb` will be automatically updated to:
```ruby
url "https://github.com/aryasoni98/wooak/archive/v0.0.3.tar.gz"
sha256 "<auto-calculated-hash>"
```

### Installation Command
After release, users can install with:
```bash
brew tap aryasoni98/homebrew-tap
brew install wooak
```

## 🔧 Troubleshooting

### If Homebrew Formula Doesn't Update

1. **Check GoReleaser logs** in GitHub Actions
2. **Verify repository access** - GoReleaser needs write access to `homebrew-tap`
3. **Check GitHub token** - Ensure GITHUB_TOKEN has proper permissions
4. **Manual update** (if needed):
   ```bash
   cd homebrew-tap
   # Calculate SHA256
   curl -L https://github.com/aryasoni98/wooak/archive/v0.0.3.tar.gz | shasum -a 256
   # Update formula manually
   # Commit and push
   ```

### If Installation Fails

1. **Check formula syntax**: `brew audit --strict Formula/wooak.rb`
2. **Test locally**: `brew install --build-from-source ./Formula/wooak.rb`
3. **Check Go version**: Ensure compatible Go version is available
4. **Check dependencies**: Verify all dependencies are available

## 📝 Summary

**Status**: ✅ **READY FOR RELEASE**

All files are updated and verified. The Homebrew formula is correctly configured and will be automatically updated by GoReleaser when you create the GitHub release.

**Next Action**: Create and push the v0.0.3 tag to trigger the release process.
