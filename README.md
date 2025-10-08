# Video Generator

CLI tool for generating videos using OpenAI's Sora API. Interactive TUI and non-interactive modes supported.

## Requirements

- Go 1.21+ (for building from source)
- OpenAI API key with Sora access

## Build & Run

```bash
# Build from source
git clone https://github.com/gersham/go-sora-video-gen.git
cd go-sora-video-gen
make build
./dist/video-gen

# Or run directly
go run .
```

## Build Commands

```bash
make build              # Build for current platform
make build-all         # Build for all platforms (macOS, Linux, Windows)
make dist              # Create distribution archives
make clean             # Clean build artifacts
```

Cross-platform binaries are created in `./dist/`, archives in `./releases/`.

## Usage

### Interactive Mode

```bash
./video-gen
```

The TUI guides you through video generation. On first run, you'll enter your OpenAI API key which is saved to `~/.config/telemetryos-video-gen.toml`.

**Keyboard Shortcuts:**
- `Ctrl+U` - Clear the current input field
- `Ctrl+C` / `Esc` - Quit the application
- `Enter` - Submit input or retry after error

**Smart Features:**
- Your last prompt is automatically saved and pre-filled on the next run
- After an error (e.g., moderation block), press Enter to retry with the previous prompt pre-filled for easy editing

### Non-Interactive CLI Mode

```bash
# Minimal
./video-gen -p "A sunset over the ocean"

# Full options
./video-gen \
  -p "A cat playing with yarn" \
  -m sora-pro \
  -t 8 \
  -s 720x1280 \
  -r ~/Desktop/cat.jpg \
  -o ~/Desktop

# With debug output
./video-gen -p "Mountain landscape" -d
```

## CLI Flags

| Flag | Options | Default |
|------|---------|---------|
| `-p` | Prompt text (triggers non-interactive mode) | - |
| `-m` | `sora` or `sora-pro` | `sora` |
| `-t` | `4`, `8`, or `12` seconds | `4` |
| `-s` | `1280x720`, `720x1280`, `1792x1024`, `1024x1792` | `1280x720` |
| `-r` | Path to image file (auto-resizes to match size) | - |
| `-o` | Output directory | `~/Desktop` |
| `-d` | Enable debug mode | `false` |

## Reference Images

When using the `-r` flag to provide a reference image, the image is automatically processed to match your selected video dimensions:

**How It Works:**
- **Automatic Resizing & Cropping** - Your image is resized and center-cropped to exactly match the video dimensions
- **Cover Strategy** - The image is scaled to cover the entire frame, then cropped to fit (similar to CSS `background-size: cover`)
- **Preserves Quality** - Images are processed at 95% JPEG quality to maintain visual fidelity

**Tips for Best Results:**
- **Match Aspect Ratios** - For best results, use images with similar aspect ratios to your target video size:
  - `1280x720` or `1792x1024` → Use landscape images (16:9 or 16:10)
  - `720x1280` or `1024x1792` → Use portrait images (9:16 or 10:16)
- **Avoid Important Edge Content** - Keep important subjects centered, as edges may be cropped
- **Use High Resolution** - Start with images at least as large as your target dimensions
- **Supported Formats** - JPEG, PNG, and GIF

**Examples:**
```bash
# Landscape video with landscape reference image (best match)
./video-gen -p "Ocean waves" -s 1280x720 -r ~/Photos/beach.jpg

# Portrait video with portrait reference image (best match)
./video-gen -p "City street" -s 720x1280 -r ~/Photos/portrait.jpg

# Landscape video with portrait image (will crop sides)
./video-gen -p "Sunset" -s 1920x1080 -r ~/Photos/tall-image.jpg
```

## Configuration

Config file: `~/.config/telemetryos-video-gen.toml`

```toml
openai_api_key = "sk-..."
output_dir = "/Users/username/Desktop"
model = "sora-2"
duration = "4"
size = "1280x720"
last_prompt = "A sunset over the ocean"
```

The `last_prompt` field is automatically saved after each video generation and is pre-filled when you restart the application.

## License

MIT License - see LICENSE file for details.
