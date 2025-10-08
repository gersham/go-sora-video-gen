# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a **Go CLI application** for generating videos using OpenAI's Sora API. It provides both an interactive Terminal User Interface (TUI) and a non-interactive CLI mode.

**Tech Stack:** Go 1.21+, Bubble Tea (TUI), OpenAI Sora API

## Project Structure

```
/Users/gersham/Sources/telemetry/telemetry-video-gen/
├── main.go                      # Entry point, routes to TUI or CLI mode
├── internal/
│   ├── api/
│   │   ├── sora.go             # OpenAI Sora API client
│   │   └── image.go            # Image resizing utilities
│   ├── tui/
│   │   └── model.go            # Bubble Tea TUI implementation
│   ├── cli/
│   │   └── cli.go              # Non-interactive CLI mode
│   └── config/
│       └── config.go           # Config management (~/.config/telemetryos-video-gen.toml)
├── Makefile                     # Build commands
├── README.md                    # User documentation
├── CHANGELOG.md                 # Version history
└── LICENSE                      # MIT License
```

## Architecture

### Entry Point (main.go)
- Parses CLI flags (`-p`, `-m`, `-t`, `-s`, `-r`, `-o`, `-d`)
- Routes to CLI mode if `-p` flag provided (non-interactive)
- Routes to TUI mode otherwise (interactive)

### API Client (internal/api/sora.go)
- Implements OpenAI Sora API v1 endpoints
- Key methods:
  - `CreateVideo()` - Initiate video generation (with retry logic)
  - `GetVideo()` - Poll video status
  - `DownloadVideoContent()` - Download completed video via `/content` endpoint
  - `ListVideos()` - List recent video jobs
  - `DeleteVideo()` - Delete video job
- Error handling uses custom `ErrorObject` struct for API errors
- Includes debug logging capability
- Auto-resizes reference images to match target video dimensions

### TUI (internal/tui/model.go)
**Framework:** Bubble Tea (elm-architecture pattern)

**States:**
- `stateAPIKey` - First-run API key entry
- `stateListVideos` - Show recent videos, option to delete
- `stateDeletingVideos` - Batch delete in progress
- `statePrompt` - Enter video prompt
- `stateModel` - Select model (sora-2 or sora-2-pro)
- `stateReferenceImage` - Optional reference image
- `stateDuration` - Select duration (4s, 8s, 12s)
- `stateSize` - Select dimensions
- `stateOutputDir` - Set output directory
- `stateGenerating` - Submitting to API
- `statePolling` - Waiting for completion
- `stateDownloading` - Downloading video
- `stateComplete` - Success, ready for next
- `stateError` - Error occurred

**Key Features:**
- Smart prompt management: Last prompt saved to config and pre-filled on restart
- Error recovery: Press Enter after error to retry with previous prompt pre-filled
- Keyboard shortcuts:
  - `Ctrl+U` - Clear input field
  - `Ctrl+C` / `Esc` - Quit
  - `Enter` - Submit/retry
  - Arrow keys - Navigate selections
- Progress tracking with elapsed time and percentage
- Auto-cleanup of old videos on startup (optional)

**Messages (Bubble Tea commands):**
- `videoCreatedMsg` - Video job started
- `videoReadyMsg` - Video ready for download
- `videoDownloadedMsg` - Download complete
- `errorMsg` - Error occurred
- `pollMsg` - Status update during polling
- `videosListedMsg` - Recent videos fetched
- `videoDeletedMsg` - Single video deleted
- `videosDeletedMsg` - All videos deleted
- `tickMsg` - Timer tick for elapsed time

### CLI (internal/cli/cli.go)
- Non-interactive mode triggered by `-p` flag
- Outputs only essential info to stdout (for automation)
- Polls video status until complete or failed
- Returns exit code 0 on success, 1 on error

### Config (internal/config/config.go)
**File:** `~/.config/telemetryos-video-gen.toml`

**Fields:**
- `openai_api_key` - OpenAI API key (required)
- `output_dir` - Default output directory
- `model` - Default model (sora-2 or sora-2-pro)
- `duration` - Default duration (4, 8, or 12)
- `size` - Default dimensions
- `last_prompt` - Last used prompt (auto-saved)

### Build System (Makefile)
- `make build` - Build for current platform
- `make build-all` - Build for all platforms
- `make dist` - Create distribution archives
- `make clean` - Clean build artifacts

**Output:**
- Binaries in `./dist/`
- Archives in `./releases/`

**Platforms:**
- macOS (Intel and Apple Silicon)
- Linux (amd64 and ARM64)
- Windows (amd64)

## Development Commands

```bash
# Run in interactive mode
go run .

# Run in CLI mode
go run . -p "A sunset over the ocean"

# Run with debug output
go run . -d

# Build for current platform
make build

# Build distribution archives
make dist

# Test
go test ./...

# Clean build artifacts
make clean
```

## Key Implementation Details

### API Error Handling
The Sora API returns errors as objects, not strings:
```go
type ErrorObject struct {
    Message string `json:"message"`
    Type    string `json:"type"`
    Code    string `json:"code"`
}
```
Always check `resp.Error != nil && resp.Error.Message != ""` when handling errors.

### Video Download
Use the `/videos/{id}/content` endpoint directly, NOT the `video_url` field from the API response.

### Reference Images
Images are automatically resized to match the target video dimensions using a "cover" strategy (resize and crop to fill).

### Polling Strategy
Smart polling intervals based on video status:
- Initial: 10 seconds
- After 50% progress: 5 seconds
- Near completion: 30 seconds

### Config Persistence
Config is saved after:
- API key entry
- Output directory selection
- Prompt submission (saves `last_prompt`)

Model, duration, and size are saved when output directory is confirmed.

### Prompt Management
- Prompts are saved to config after submission
- On app restart or error recovery, last prompt is pre-filled
- User can edit pre-filled prompt before submission
- `Ctrl+U` clears the input field

## Common Development Tasks

### Adding a New TUI State
1. Add state constant to `type state int` enum
2. Add state transition logic in `Update()`
3. Add view rendering in `View()`
4. Add any new message types needed
5. Add keyboard handling for the state

### Adding a New CLI Flag
1. Add flag in `main.go` CLI parsing
2. Add field to `CLIOptions` struct
3. Pass through to `cli.Run()` or TUI model
4. Update README.md flag table

### Adding a New Config Field
1. Add field to `Config` struct in `internal/config/config.go`
2. Update load/save logic if needed
3. Use in TUI or CLI as appropriate
4. Update README.md config example

### Modifying API Client
- All API methods are in `internal/api/sora.go`
- Debug logging is built-in via `debugLog` callback
- Error handling should use `ErrorObject` struct
- Add retry logic for transient errors

## Debugging

### Enable Debug Mode
```bash
go run . -d
```

Debug output includes:
- Full HTTP requests (method, URL, headers, body)
- Full HTTP responses (status, body)
- API polling events
- State transitions

Debug logs are shown in TUI at bottom of screen (last 50 entries).

### Common Issues

**"Error: json: cannot unmarshal object into Go struct field VideoResponse.error of type string"**
- API returned error object but code expected string
- Ensure `VideoResponse.Error` is `*ErrorObject`, not `string`

**"Video content not available after N attempts"**
- Video may still be processing
- Check API status with debug mode
- Verify using correct endpoint (`/content` not `video_url`)

**"Your request was blocked by our moderation system"**
- Prompt violated OpenAI content policy
- User can press Enter to retry with edited prompt

**Reference image dimension mismatch**
- Image will be auto-resized to match video dimensions
- Uses "cover" strategy (may crop image)

## Testing

Run tests:
```bash
go test ./...
```

Key areas to test:
- Config load/save
- API client methods
- Image resizing logic
- CLI flag parsing
- State transitions in TUI

## Release Process

1. Update version in `Makefile` (`VERSION` variable)
2. Update `CHANGELOG.md` with changes
3. Run `make dist` to create distribution archives
4. Upload archives from `./releases/` to GitHub release
5. Tag release with version (e.g., `v1.0.0`)

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Add comments for complex logic
- Use descriptive variable names
- Handle all errors explicitly

## Dependencies

Main dependencies:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/charmbracelet/lipgloss` - TUI styling
- `github.com/BurntSushi/toml` - Config file parsing

## License

MIT License - see LICENSE file for details.

---

**Remember:** This is a CLI tool first, TUI second. The non-interactive mode (`-p` flag) should always work for automation purposes. The TUI is for convenience and exploration.
