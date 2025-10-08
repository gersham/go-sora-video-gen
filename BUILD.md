# Build and Release Guide

This guide covers building and packaging the TelemetryOS Video Generator for distribution.

## Prerequisites

- Go 1.21 or higher
- Make (included on macOS/Linux, use WSL on Windows)
- Git (for version management)

## Quick Build

### Development Build

For local testing and development:

```bash
go build -o video-gen
```

### Production Build (Current Platform)

For an optimized binary on your current platform:

```bash
make build
```

Output: `./dist/video-gen` (smaller, stripped binary ~6.5MB vs ~9.4MB)

## Cross-Platform Builds

### Build for All Platforms

```bash
make build-all
```

This creates optimized binaries in `./dist/` for:
- macOS Intel (darwin-amd64)
- macOS Apple Silicon (darwin-arm64)
- Linux x64 (linux-amd64)
- Linux ARM64 (linux-arm64)
- Windows x64 (windows-amd64)

### Build for Specific Platform

```bash
# macOS Intel
make build-darwin-amd64

# macOS Apple Silicon
make build-darwin-arm64

# Linux x64
make build-linux-amd64

# Linux ARM64
make build-linux-arm64

# Windows x64
make build-windows-amd64
```

## Creating Distribution Archives

### Package All Platforms

```bash
make dist
```

This creates compressed archives in `./releases/`:
- `video-gen-darwin-amd64.tar.gz` (macOS Intel)
- `video-gen-darwin-arm64.tar.gz` (macOS Apple Silicon)
- `video-gen-linux-amd64.tar.gz` (Linux x64)
- `video-gen-linux-arm64.tar.gz` (Linux ARM64)
- `video-gen-windows-amd64.zip` (Windows x64)

Each archive contains:
- The optimized binary
- README.md
- LICENSE
- example.toml (configuration template)

### Package Specific Platform

```bash
# For your current platform (macOS ARM in this example)
make dist-darwin-arm64
```

## Size Optimization

The Makefile uses several optimizations to reduce binary size:

1. **`-ldflags "-s -w"`** - Strip debugging symbols and DWARF info
2. **`-trimpath`** - Remove file system paths from binary

Results:
- Development build: ~9.4 MB
- Production build: ~6.5 MB (~30% reduction)
- Compressed archive: ~2.6 MB (~72% reduction)

## Versioning

Set version during build:

```bash
make build VERSION=1.2.3
```

Or for distribution:

```bash
make dist VERSION=1.2.3
```

## Testing the Build

After building, test the binary:

```bash
# Run the binary
./dist/video-gen

# Test that it shows the TUI and prompts for API key
# Press Ctrl+C to exit
```

## Distribution Checklist

Before creating a release:

- [ ] Update version in `README.md` badges
- [ ] Update `CHANGELOG.md` (if exists)
- [ ] Run `make clean` to remove old artifacts
- [ ] Run `make dist` to create fresh archives
- [ ] Test binaries on target platforms
- [ ] Verify archive contents with `tar -tzf` or `unzip -l`
- [ ] Create GitHub release with archives
- [ ] Update release notes

## Publishing to GitHub Releases

1. Create a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. Build distribution archives:
   ```bash
   make clean
   make dist VERSION=1.0.0
   ```

3. Create GitHub release:
   - Go to https://github.com/telemetry/video-gen/releases/new
   - Select the tag (v1.0.0)
   - Fill in release title and notes
   - Upload archives from `./releases/`
   - Publish release

## Cleaning Build Artifacts

Remove all build artifacts:

```bash
make clean
```

This removes:
- `./dist/` directory (binaries)
- `./releases/` directory (archives)
- Root-level `video-gen` binary

## Troubleshooting

### Build Errors

**"command not found: make"**
- Install make: `brew install make` (macOS) or `sudo apt install make` (Linux)

**"Go version too old"**
- Update Go: Download from https://go.dev/dl/

**Cross-compilation fails**
- Ensure you have latest Go toolchain
- Run `go mod download` first

### Size Issues

**Binary too large**
- Use `make build` instead of `go build`
- Ensure `-ldflags "-s -w"` flags are applied
- Consider using UPX for additional compression (optional)

### Archive Issues

**tar/zip command not found**
- tar is standard on macOS/Linux
- zip may need: `sudo apt install zip` (Linux)

## Advanced: Custom Builds

### Building with Custom Flags

Edit the Makefile `LDFLAGS` variable:

```makefile
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.buildDate=$(shell date -u +%Y-%m-%d)"
```

### Adding Build Information

In `main.go`, add variables:

```go
var (
    version   = "dev"
    buildDate = "unknown"
)

func main() {
    fmt.Printf("Version: %s, Built: %s\n", version, buildDate)
    // ... rest of code
}
```

### Platform-Specific Optimizations

For production deployments, consider:

```bash
# Linux with static linking
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -extldflags '-static'" -o dist/video-gen-linux-amd64 .
```

## Release Workflow Summary

Complete release workflow:

```bash
# 1. Ensure code is committed
git status

# 2. Clean previous builds
make clean

# 3. Build and package all platforms
make dist VERSION=1.0.0

# 4. Verify archives
ls -lh releases/

# 5. Test the binary (current platform)
./dist/video-gen

# 6. Create git tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 7. Upload to GitHub Releases
# (Manual step via GitHub web interface)

# 8. Announce release
# (Update documentation, notify users, etc.)
```

---

**Questions?** Open an issue on GitHub or contact the maintainers.
