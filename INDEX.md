# Documentation Index

Complete guide to TelemetryOS Video Generator documentation.

## 📖 Documentation Files

### 🚀 Getting Started

| File | Description | Use When |
|------|-------------|----------|
| **[QUICKSTART.md](QUICKSTART.md)** | Fast 60-second setup guide | You're new and want to get started immediately |
| **[README.md](README.md)** | Main project documentation | You want comprehensive feature overview |
| **[example.toml](example.toml)** | Configuration template | You need to configure manually |

### 🔧 Building & Distribution

| File | Description | Use When |
|------|-------------|----------|
| **[BUILD.md](BUILD.md)** | Complete build and release guide | You want to build from source or create releases |
| **[Makefile](Makefile)** | Build automation | You want to see available build commands |

### 🐛 Troubleshooting

| File | Description | Use When |
|------|-------------|----------|
| **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** | Complete troubleshooting guide | You encounter errors or issues |

### 📝 Reference

| File | Description | Use When |
|------|-------------|----------|
| **[SUMMARY.md](SUMMARY.md)** | Project summary and status | You want a high-level overview |
| **[LICENSE](LICENSE)** | MIT License | You need licensing information |

## 🎯 Quick Navigation

### I want to...

**...get started quickly**
→ [QUICKSTART.md](QUICKSTART.md)

**...understand all features**
→ [README.md](README.md)

**...build from source**
→ [BUILD.md](BUILD.md) or run `make help`

**...create a release**
→ [BUILD.md](BUILD.md) → "Release Workflow Summary"

**...fix an error**
→ [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

**...see what's included**
→ [SUMMARY.md](SUMMARY.md)

**...configure manually**
→ [example.toml](example.toml)

## 📚 Documentation by User Type

### First-Time Users
1. Start with [QUICKSTART.md](QUICKSTART.md)
2. Explore features in [README.md](README.md)
3. Keep [TROUBLESHOOTING.md](TROUBLESHOOTING.md) handy

### Developers
1. Read [README.md](README.md) for architecture
2. Follow [BUILD.md](BUILD.md) for development setup
3. Review [SUMMARY.md](SUMMARY.md) for technical specs

### Contributors
1. Review [SUMMARY.md](SUMMARY.md) for project structure
2. Study [BUILD.md](BUILD.md) for build system
3. Check [README.md](README.md) → Contributing section

### Distributors
1. Follow [BUILD.md](BUILD.md) → "Release Workflow Summary"
2. Use `make dist` to create archives
3. Reference [README.md](README.md) for feature list

## 🗂️ Project Structure

```
telemetry-video-gen/
├── 📄 Documentation
│   ├── README.md              # Main documentation
│   ├── QUICKSTART.md          # Fast setup guide
│   ├── BUILD.md               # Build & release guide
│   ├── TROUBLESHOOTING.md     # Error solutions
│   ├── SUMMARY.md             # Project overview
│   ├── INDEX.md               # This file
│   └── example.toml           # Config template
│
├── 📦 Application Code
│   ├── main.go                # Entry point
│   ├── internal/
│   │   ├── api/sora.go        # Sora API client
│   │   ├── config/config.go   # Config management
│   │   └── tui/model.go       # TUI interface
│   ├── go.mod                 # Dependencies
│   └── go.sum                 # Checksums
│
├── 🔨 Build System
│   ├── Makefile               # Build automation
│   ├── .gitignore             # Git exclusions
│   └── LICENSE                # MIT license
│
├── 📦 Build Artifacts (generated)
│   ├── dist/                  # Optimized binaries
│   └── releases/              # Distribution archives
│
└── 📝 Configuration (user)
    └── ~/.config/telemetryos-video-gen.toml
```

## 🎓 Learning Path

### Beginner
```
QUICKSTART.md → Run app → README.md (Usage section)
                    ↓
              (if errors)
                    ↓
           TROUBLESHOOTING.md
```

### Intermediate
```
README.md (full) → BUILD.md (build section) → make build
                        ↓
                  Try customizing
                        ↓
              SUMMARY.md (specs)
```

### Advanced
```
SUMMARY.md → BUILD.md (advanced) → Source code
     ↓
Contribute changes
     ↓
CREATE.md → make dist → Release
```

## 📊 Documentation Stats

| File | Lines | Size | Purpose |
|------|-------|------|---------|
| README.md | 204 | 5.1KB | Main docs |
| QUICKSTART.md | 195 | 5.7KB | Fast start |
| BUILD.md | 240 | 5.2KB | Build guide |
| TROUBLESHOOTING.md | 350 | 7.5KB | Error help |
| SUMMARY.md | 260 | 6.6KB | Overview |
| INDEX.md | 150 | 3.5KB | This file |

**Total Documentation**: ~1,400 lines, ~33KB

## 🔍 Search Tips

### Find Specific Information

**API Errors**
→ TROUBLESHOOTING.md → "API Errors"

**Build Commands**
→ Run `make help` or see BUILD.md

**Example Prompts**
→ QUICKSTART.md → "Example Prompts"

**Configuration Options**
→ README.md → "Configuration" or example.toml

**Retry Logic**
→ SUMMARY.md → "Error Handling & Reliability"

**Cross-Platform Builds**
→ BUILD.md → "Cross-Platform Builds"

**Release Process**
→ BUILD.md → "Release Workflow Summary"

## 🆘 Getting Help

1. **Check docs first**: Use this index to find relevant file
2. **Search issues**: https://github.com/telemetry/video-gen/issues
3. **Ask community**: GitHub Discussions
4. **Contact support**: Open a new issue

## 🤝 Contributing to Docs

Documentation improvements are welcome!

When updating docs:
- Keep this INDEX.md in sync
- Update relevant quick navigation links
- Maintain consistent formatting
- Test all commands and examples
- Update line counts in stats

## 📌 Version

- **Documentation Version**: 1.0.0
- **Last Updated**: 2025-10-06
- **Status**: Complete

---

**Need help?** Start with [QUICKSTART.md](QUICKSTART.md) or [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
