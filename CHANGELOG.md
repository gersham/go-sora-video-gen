# Changelog

All notable changes to TelemetryOS Video Generator will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-06

### Added
- ✨ Initial release of TelemetryOS Video Generator
- 🎨 Beautiful TUI interface built with Bubble Tea
- 🔐 Automatic API key management with TOML config
- 🖼️ Reference image support for guided generation
- ⏱️ **Configurable video duration (1-60 seconds, default 4s)**
- 📁 Flexible output directory configuration
- 📥 Automatic video download with timestamp naming
- 🔄 Exponential backoff retry logic (3 attempts)
- 🛡️ Robust error handling and validation
- 📦 Cross-platform build support (5 platforms)
- 🚀 GitHub Actions workflow for automated releases
- 📚 Comprehensive documentation:
  - README.md - Main documentation
  - QUICKSTART.md - Fast setup guide
  - BUILD.md - Build and release guide
  - TROUBLESHOOTING.md - Complete troubleshooting
  - SUMMARY.md - Project overview
  - INDEX.md - Documentation index
  - CHANGELOG.md - This file

### Features in Detail

#### Core Functionality
- Interactive prompt-based workflow
- OpenAI Sora API integration (sora-2 model)
- Real-time progress updates during generation
- Automatic polling until video is ready (max 3 minutes)
- MP4 video output with configurable duration
- Default resolution: 720x1280 (portrait)

#### Error Handling
- Automatic retry with exponential backoff (2s, 4s, 8s delays)
- Structured API error parsing
- Client vs server error differentiation
- Graceful timeout handling
- User-friendly error messages

#### Build System
- Comprehensive Makefile with 15+ targets
- Size optimization (~30% reduction: 9.4MB → 6.5MB)
- Cross-compilation for:
  - macOS Intel (darwin-amd64)
  - macOS Apple Silicon (darwin-arm64)
  - Linux x64 (linux-amd64)
  - Linux ARM64 (linux-arm64)
  - Windows x64 (windows-amd64)
- Distribution archive creation (.tar.gz, .zip)
- Version injection support

#### Developer Experience
- Clean project structure with internal packages
- Well-documented code
- Example configuration file
- Comprehensive .gitignore

#### CI/CD
- GitHub Actions workflow for releases
- Automatic checksum generation
- Multi-platform builds in CI
- Release notes generation

### Configuration

Default config location: `~/.config/telemetryos-video-gen.toml`

```toml
openai_api_key = "sk-..."
output_dir = "/Users/username/Desktop"
```

### Usage Flow

1. Enter OpenAI API key (first time only)
2. Describe video in natural language prompt
3. Optionally provide reference image path
4. **Specify video duration (1-60s, default 4s)**
5. Confirm output directory
6. Wait for generation (1-5 minutes)
7. Video automatically downloaded to specified directory

### Technical Specifications

- **Language**: Go 1.21+
- **TUI Framework**: Bubble Tea
- **Config Format**: TOML
- **API**: OpenAI Sora `/videos` endpoint
- **Dependencies**:
  - github.com/BurntSushi/toml v1.3.2
  - github.com/charmbracelet/bubbles v0.18.0
  - github.com/charmbracelet/bubbletea v0.25.0
  - github.com/charmbracelet/lipgloss v0.9.1

### Known Limitations

- Sora API is in limited beta (requires access)
- Video generation can take 1-5 minutes
- Maximum duration: 60 seconds
- Fixed resolution: 720x1280
- No batch processing support
- No history tracking

### Installation

#### Pre-built Binaries
Download from [GitHub Releases](https://github.com/telemetry/video-gen/releases)

#### Build from Source
```bash
git clone https://github.com/telemetry/video-gen.git
cd video-gen
make build
./dist/video-gen
```

### Requirements

- OpenAI API key with Sora access
- Internet connection
- ~500MB free disk space
- macOS 10.15+, Linux (any modern distro), or Windows 10+

---

## [Unreleased]

### Planned Features
- [ ] CLI flags for automation
- [ ] Batch processing from file
- [ ] Progress bar for downloads
- [ ] History tracking
- [ ] Custom resolution support
- [ ] Video preview before download
- [ ] Verbose/debug mode
- [ ] Configuration presets

---

## Version History

- **1.0.0** (2025-10-06) - Initial release with full feature set

## Release Process

1. Update version in this CHANGELOG
2. Update version in Makefile if needed
3. Commit changes: `git commit -am "Bump version to X.Y.Z"`
4. Create tag: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
5. Push: `git push origin main --tags`
6. GitHub Actions automatically builds and publishes release

## Contributing

See [README.md](README.md) for contribution guidelines.

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

---

**Questions?** Open an issue on [GitHub](https://github.com/telemetry/video-gen/issues)
