# Quick Start Guide

Get up and running with TelemetryOS Video Generator in 60 seconds.

## Installation

### Option A: Download Pre-built Binary (Fastest)

1. **Download** the binary for your platform:
   - macOS (Apple Silicon): `video-gen-darwin-arm64.tar.gz`
   - macOS (Intel): `video-gen-darwin-amd64.tar.gz`
   - Linux (x64): `video-gen-linux-amd64.tar.gz`
   - Windows: `video-gen-windows-amd64.zip`

2. **Extract** the archive:
   ```bash
   # macOS/Linux
   tar -xzf video-gen-*.tar.gz
   cd video-gen

   # Windows
   # Extract using Windows Explorer or: unzip video-gen-*.zip
   ```

3. **Run** the application:
   ```bash
   ./video-gen
   ```

### Option B: Build from Source

```bash
git clone https://github.com/telemetry/video-gen.git
cd video-gen
go build -o video-gen
./video-gen
```

## First Run

1. **Enter API Key** (first time only)
   ```
   Enter your OpenAI API key: sk-proj-...
   ```

   Get your key from: https://platform.openai.com/api-keys

2. **Video Management** (optional)
   ```
   Found 3 recent videos. Delete all? [Yes/No]
   > Yes (default)
   ```
   Clean up old video jobs on startup

3. **Enter Prompt**
   ```
   Describe the video you want to generate:
   > A serene sunset over the ocean with gentle waves
   ```
   Press Enter with empty prompt to exit

4. **Model Selection**
   ```
   Select model (use arrow keys):
   → sora - Fast (sora-2)
     sora-pro - Quality (sora-2-pro)
   ```
   Use arrow keys to choose

5. **Reference Image** (optional)
   ```
   Reference image path (optional):
   > ~/Desktop/sunset.jpg
   ```
   Supports tilde expansion - images auto-resize to match video size

6. **Video Duration**
   ```
   Select video duration (use arrow keys):
   → 4 - 4 seconds
     8 - 8 seconds
     12 - 12 seconds
   ```
   Use arrow keys to select

7. **Video Size**
   ```
   Select video size (use arrow keys):
   → 1280x720 - Landscape (HD)
     720x1280 - Portrait (HD)
     1792x1024 - Landscape Wide
     1024x1792 - Portrait Wide
   ```
   Use arrow keys to choose

8. **Output Directory**
   ```
   Output directory:
   > /Users/you/Desktop
   ```
   Confirm or change the location

9. **Wait for Generation**
   ```
   ⣟  Generating video (45s) in_progress (73% complete)
   Polling API every 10s (attempt 5/120)
   ```
   - Video generation takes 1-5 minutes
   - Smart polling: 10s (first 2min) → 5s (at 100%) → 30s (thereafter)
   - Progress updates shown in real-time
   - Automatic retries if needed

10. **Done!**
   ```
   ✓ Video generated successfully!
   Saved to: /Users/you/Desktop/sora_video_20251006_120000.mp4

   Deleting video from service...
   ✓ Video deleted from service
   ```
   Video is automatically cleaned up from OpenAI after download

## Example Prompts

### Good Prompts (Descriptive & Specific)
- "A cat playing with a red ball of yarn in a cozy living room, warm afternoon lighting"
- "Time-lapse of clouds moving across a bright blue sky, 4K quality"
- "A person walking through a forest in autumn, leaves falling, cinematic"
- "Close-up of rain drops on a window with blurred city lights in the background"
- "A steaming cup of coffee on a wooden table, morning sunlight"

### Tips for Better Results
- Be specific about lighting, mood, and setting
- Mention camera movement if desired (pan, zoom, static)
- Include style keywords (cinematic, realistic, artistic)
- Keep prompts under 500 characters
- Use reference images to guide style

## Common Commands

### Run Interactively
```bash
./video-gen
```

### Run in CLI Mode (Non-Interactive)
```bash
# Quick generation with prompt only
./video-gen -prompt "A sunset over the ocean"

# Full control with all parameters
./video-gen \
  -prompt "A cat playing with yarn" \
  -model sora-pro \
  -duration 8 \
  -size 720x1280 \
  -reference ~/Desktop/cat.jpg \
  -output ~/Videos
```

### Check All CLI Flags
```bash
./video-gen -h
```

### Check Configuration
```bash
cat ~/.config/telemetryos-video-gen.toml
```

### Reset Configuration
```bash
rm ~/.config/telemetryos-video-gen.toml
./video-gen
```

### Build from Source
```bash
git clone https://github.com/telemetry/video-gen.git
cd video-gen
make build
./dist/video-gen
```

### Create Distribution Package
```bash
make dist
# Archives created in ./releases/
```

## Keyboard Shortcuts

- `Enter` - Confirm and proceed to next step
- `Enter` (empty prompt) - Exit application
- `↑↓` or `←→` - Navigate selections (model, duration, size, delete option)
- `Ctrl+C` - Exit application (at any time)
- `Esc` - Exit application (at any time)

## Debug Mode

Enable debug mode to see detailed API communication:

```bash
./dist/video-gen -d
# or
go run . -d
```

Shows:
- API requests (blue) with method, URL, headers, body
- API responses (magenta) with status code and JSON
- Last 10 interactions displayed below main UI
- Helpful for troubleshooting 400/500 errors

## Troubleshooting

### "API error (500): Server error"
- **Common** during Sora beta
- App automatically retries 3 times
- If it persists, check https://status.openai.com
- Verify your account has Sora access

### "API error (401): Invalid API key"
- Check your API key in `~/.config/telemetryos-video-gen.toml`
- Generate new key at https://platform.openai.com/api-keys
- Reset: `rm ~/.config/telemetryos-video-gen.toml`

### "File does not exist" (reference image)
- Use tilde expansion: `~/Desktop/image.jpg`
- Or absolute path: `/Users/you/Pictures/image.jpg`
- Press Enter to skip reference image
- Verify file exists: `ls -la ~/Desktop/image.jpg`

### "Invalid duration" error
- Only 4, 8, or 12 seconds are supported
- Use arrow keys in interactive mode
- Use `-duration 4`, `-duration 8`, or `-duration 12` in CLI mode

### Video Generation Taking Too Long
- Normal generation time: 1-5 minutes
- Smart polling intervals: 10s → 5s → 30s
- Timeout after 60 minutes (120 attempts)
- Check OpenAI dashboard for job status

### More Help
- See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for complete guide
- Open issue: https://github.com/telemetry/video-gen/issues

## Configuration File

Location: `~/.config/telemetryos-video-gen.toml`

```toml
# Your OpenAI API key
openai_api_key = "sk-proj-..."

# Output directory for generated videos
output_dir = "/Users/you/Desktop"
```

**Note:** This file is auto-created on first run. You can edit it manually to change settings.

## Video Output

Videos are saved with timestamp naming:
```
sora_video_20251006_120000.mp4
           ^^^^^^^^_^^^^^^
           YYYYMMDD_HHMMSS
```

**Default Properties:**
- Format: MP4
- Resolution: 1280x720 (landscape HD)
- Duration: 4 seconds (4, 8, or 12 seconds supported)
- Model: sora-2
- Auto-deleted from service after download

## System Requirements

- **Operating System**: macOS 10.15+, Linux (any modern distro), Windows 10+
- **Internet**: Required for API access
- **Disk Space**: ~500MB free (for temporary files and videos)
- **OpenAI Account**: With Sora API access and credits

## Getting Sora Access

Sora API is in limited beta as of early 2025:

1. **Sign up**: https://platform.openai.com
2. **Request access**: Through OpenAI's beta program
3. **Add credits**: Add billing to your account
4. **Generate key**: https://platform.openai.com/api-keys
5. **Start creating**: Use this app!

## Support & Resources

- **Documentation**: See [README.md](README.md)
- **Build Guide**: See [BUILD.md](BUILD.md)
- **Troubleshooting**: See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- **Issues**: https://github.com/telemetry/video-gen/issues
- **OpenAI Docs**: https://platform.openai.com/docs/
- **OpenAI Status**: https://status.openai.com

## Next Steps

1. ✅ Install and run the app
2. ✅ Generate your first video
3. ✅ Experiment with different prompts
4. ✅ Try reference images for style guidance
5. 📚 Read [README.md](README.md) for full documentation
6. 🔧 Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md) if you hit issues
7. 🛠️ See [BUILD.md](BUILD.md) to build from source

---

**Happy video generating! 🎬**
