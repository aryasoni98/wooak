# Fix Homebrew Error

The error you're seeing is a **local Homebrew installation issue**, not related to the Wooak formula. This is a common Ruby gem compatibility issue with Homebrew.

## Quick Fix

Try these solutions in order:

### Solution 1: Update Homebrew
```bash
brew update
brew doctor
```

### Solution 2: Reinstall Homebrew's Ruby gems
```bash
cd /opt/homebrew/Library/Homebrew
rm -rf vendor/bundle
brew update
```

### Solution 3: Clear Homebrew cache
```bash
brew cleanup
rm -rf ~/Library/Caches/Homebrew
brew update
```

### Solution 4: Reinstall Homebrew (last resort)
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh)"
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## Why Wooak Installation Fails

The `wooak --version` command fails because:

1. **GitHub Release Not Created Yet**: The v0.0.3 release needs to be created on GitHub first
2. **Homebrew Formula Not Updated**: GoReleaser will automatically update the formula when the release is created
3. **No Binaries Available**: The release binaries need to be built and uploaded

## Next Steps

### Option 1: Create Release Manually (Recommended)

1. Go to: https://github.com/aryasoni98/wooak/releases/new
2. Select tag: `v0.0.3`
3. Title: `v0.0.3`
4. Description: Copy from `CHANGELOG.md` v0.0.3 section
5. Click "Publish release"

Then run GoReleaser to build binaries:
```bash
# Install GoReleaser if needed
brew install goreleaser

# Create release with binaries
cd /Users/aryasoni/Documents/CodeBase/project/wooak/wooak
goreleaser release --clean
```

### Option 2: Update GitHub Actions Workflow

Add tag trigger to `.github/workflows/ci.yml`:

```yaml
on:
  push:
    branches: [master, main]
    tags:
      - 'v*'
  pull_request:
```

Then push a commit to trigger the workflow.

### Option 3: Test Locally (Without Homebrew)

Build and test locally:
```bash
cd /Users/aryasoni/Documents/CodeBase/project/wooak/wooak
make build
./bin/wooak --version
```

## After Release is Created

Once the GitHub release is created and GoReleaser runs:

1. **GoReleaser will automatically**:
   - Build binaries for all platforms
   - Upload them to the GitHub release
   - Update the Homebrew formula in `homebrew-tap` repository

2. **Then you can install**:
   ```bash
   # Fix Homebrew first (see solutions above)
   brew tap aryasoni98/homebrew-tap
   brew install wooak
   wooak --version  # Should show v0.0.3
   ```

## Current Status

✅ **Tag pushed**: v0.0.3 is on GitHub  
⏳ **Release pending**: GitHub release needs to be created  
⏳ **Binaries pending**: Need to run GoReleaser  
⏳ **Formula update pending**: Will happen automatically after release  
