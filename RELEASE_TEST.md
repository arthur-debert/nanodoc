# Testing Homebrew Release

## Prerequisites

1. Ensure you have a GitHub personal access token with repo access to `homebrew-tools`
2. Set it as `HOMEBREW_TAP_TOKEN` secret in the repository settings

## Steps to Test

### 1. Create a Test Release

```bash
# After PR is merged to main
git checkout main
git pull origin main

# Create a test tag
git tag v0.1.0-beta.1
git push origin v0.1.0-beta.1
```

This will trigger the release workflow.

### 2. Monitor the Release

1. Go to Actions tab on GitHub
2. Watch the "Release" workflow
3. Check that it completes successfully

### 3. Verify Homebrew Formula

The workflow should create/update a formula in https://github.com/arthur-debert/homebrew-tools

Check that:
- Formula file exists at `Formula/nanodoc.rb`
- It includes completion installation
- It includes man page installation

### 4. Install via Homebrew

```bash
# Add the tap
brew tap arthur-debert/tools

# Install nanodoc
brew install nanodoc

# Verify installation
nanodoc --version

# Test completions
nanodoc completion bash

# Test man page
man nanodoc
```

### 5. Test Shell Completions

#### Bash
```bash
source <(nanodoc completion bash)
# Now type: nanodoc <TAB>
```

#### Zsh
```bash
source <(nanodoc completion zsh)
# Now type: nanodoc <TAB>
```

#### Fish
```bash
nanodoc completion fish | source
# Now type: nanodoc <TAB>
```

## Troubleshooting

1. **Formula not updated**: Check that `HOMEBREW_TAP_TOKEN` has write access to the homebrew-tools repo
2. **Completions not working**: Ensure shell completion is enabled in your shell
3. **Man page not found**: Check that man path includes Homebrew's man directory

## Debug Mode

To test formula generation in debug mode:

```bash
DEBUG=1 goreleaser release --snapshot --clean
```

This will create the formula in a `debug/` directory instead of `Formula/`.