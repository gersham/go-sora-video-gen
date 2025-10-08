# Troubleshooting Guide

Common issues and solutions for TelemetryOS Video Generator.

## API Errors

### 500 Server Error

**Error Message:**
```
API error (500 - server_error): Server error
```

**Cause:** OpenAI's Sora API is experiencing temporary issues or high load.

**Solutions:**
1. **Wait and retry** - The app now automatically retries 3 times with exponential backoff (2s, 4s, 8s)
2. **Check API status** - Visit [OpenAI Status](https://status.openai.com) to see if there are known issues
3. **Try again later** - Server errors are typically temporary
4. **Verify your account** - Ensure your OpenAI account has Sora API access
5. **Check your quota** - You may have hit rate limits or quota restrictions

**Note:** As of early 2025, Sora API is in limited beta. Make sure you have:
- An OpenAI account with Sora access
- Valid API credits
- Not exceeded your usage limits

### 401 Authentication Error

**Error Message:**
```
API error (401 - invalid_api_key): Invalid API key
```

**Solutions:**
1. Check your API key is correct:
   ```bash
   cat ~/.config/telemetryos-video-gen.toml
   ```
2. Delete config and re-enter your key:
   ```bash
   rm ~/.config/telemetryos-video-gen.toml
   ./video-gen
   ```
3. Generate a new API key at [OpenAI API Keys](https://platform.openai.com/api-keys)

### 429 Rate Limit Error

**Error Message:**
```
API error (429 - rate_limit_exceeded): Rate limit exceeded
```

**Solutions:**
1. Wait a few minutes before trying again
2. Check your organization's rate limits in the OpenAI dashboard
3. Upgrade your OpenAI plan if you need higher limits

### 400 Bad Request

**Error Message:**
```
API error (400 - invalid_request_error): ...
```

**Solutions:**
1. Check your prompt isn't too long (keep under 500 characters)
2. If using a reference image, ensure:
   - File exists and path is correct
   - Image is a supported format (JPEG, PNG)
   - Image size is reasonable (< 20MB)

## Network Issues

### Connection Timeout

**Error:**
```
failed to execute request: context deadline exceeded
```

**Solutions:**
1. Check your internet connection
2. Try again - network may be temporarily unstable
3. Check if you're behind a firewall or proxy
4. Verify you can access `https://api.openai.com`

### DNS Resolution Failures

**Error:**
```
failed to execute request: dial tcp: lookup api.openai.com: no such host
```

**Solutions:**
1. Check internet connection
2. Try a different DNS server
3. Flush DNS cache:
   - macOS: `sudo dscacheutil -flushcache`
   - Linux: `sudo systemd-resolve --flush-caches`

## File System Issues

### Config Directory Permission Error

**Error:**
```
failed to create config directory: permission denied
```

**Solutions:**
1. Check permissions:
   ```bash
   ls -la ~/.config
   ```
2. Create directory manually:
   ```bash
   mkdir -p ~/.config
   chmod 755 ~/.config
   ```

### Output Directory Not Found

**Error:**
```
failed to create output directory: no such file or directory
```

**Solutions:**
1. Ensure the output directory exists or parent directory exists
2. Use full absolute paths (e.g., `/Users/username/Desktop`)
3. Check directory permissions

### Disk Space Issues

**Error:**
```
failed to write video data: no space left on device
```

**Solutions:**
1. Free up disk space (videos can be 10-100MB+)
2. Change output directory to a drive with more space
3. Check disk usage:
   ```bash
   df -h
   ```

## Application Issues

### Program Won't Start

**Symptoms:** Double-clicking does nothing or immediate crash

**Solutions:**
1. Run from terminal to see error messages:
   ```bash
   ./video-gen
   ```
2. Check if binary has execute permissions:
   ```bash
   chmod +x video-gen
   ```
3. Verify Go version if building from source:
   ```bash
   go version  # Should be 1.21+
   ```

### Text Input Not Working

**Symptoms:** Can't type in the terminal

**Solutions:**
1. Ensure terminal supports interactive input
2. Try a different terminal (iTerm2, Terminal.app, etc.)
3. Check if another process is blocking input

### Display Issues

**Symptoms:** Garbled text, weird characters, broken layout

**Solutions:**
1. Ensure terminal supports UTF-8:
   ```bash
   echo $LANG  # Should show UTF-8
   ```
2. Update terminal to latest version
3. Try a modern terminal (iTerm2, Alacritty, etc.)
4. Set UTF-8 locale:
   ```bash
   export LANG=en_US.UTF-8
   ```

## Sora API Specific Issues

### "Model Not Found" or "Model Not Available"

**Error:**
```
API error (404): model 'sora-2' not found
```

**Cause:** Your account may not have access to Sora, or the model name changed.

**Solutions:**
1. Verify Sora access in your OpenAI account dashboard
2. Check if Sora is available in your region
3. Contact OpenAI support to request access
4. Check [OpenAI Docs](https://platform.openai.com/docs/) for current model names

### Video Generation Stuck at "Generating"

**Symptoms:** Stays at "Generating video (check X)..." for a long time

**Solutions:**
1. **Be patient** - Video generation can take 1-5 minutes or longer
2. The app polls every 2-3 seconds with exponential backoff
3. Maximum wait time is ~3 minutes (60 attempts)
4. Check OpenAI dashboard for job status
5. If it times out, the job may still complete - check your OpenAI dashboard

### Video Quality Issues

**Issue:** Generated video doesn't match expectations

**Solutions:**
1. **Improve your prompt:**
   - Be specific and descriptive
   - Include details about style, motion, lighting
   - Example: "A serene sunset over a calm ocean, golden hour lighting, slow gentle waves"
2. **Use reference images** to guide generation style
3. Try different prompt variations
4. Understand Sora's capabilities and limitations

## Getting Help

### Debug Mode

Run with verbose output (if implemented):
```bash
DEBUG=1 ./video-gen
```

### Check Logs

Application logs are typically in:
- Console output (when running from terminal)
- System logs (macOS: Console.app)

### Collect Information

When reporting issues, include:
1. **Error message** (full text)
2. **Operating system** and version
3. **Application version** (`./video-gen --version` if available)
4. **Steps to reproduce**
5. **What you expected vs. what happened**
6. **API key status** (valid/invalid, not the actual key!)

### Report Bugs

Open an issue at: [GitHub Issues](https://github.com/telemetry/video-gen/issues)

Include:
- Description of the problem
- Error messages
- Steps to reproduce
- System information
- Screenshots if applicable

### Community Support

- Check existing issues on GitHub
- Search OpenAI community forums
- Review OpenAI API documentation
- Contact TelemetryOS support

## Known Limitations

1. **Sora API is in beta** - Expect occasional issues and changes
2. **Rate limits** - OpenAI enforces request limits
3. **Video generation time** - Can take 1-5 minutes per video
4. **File size** - Videos can be large (10-100MB+)
5. **Region availability** - Sora may not be available in all regions
6. **Cost** - Video generation consumes API credits

## FAQs

**Q: How long does video generation take?**
A: Typically 1-5 minutes, but can vary based on complexity and API load.

**Q: How much does it cost?**
A: Check OpenAI's pricing page for current Sora API rates.

**Q: Can I cancel a video generation?**
A: Currently, no. Once submitted, the job will complete or timeout.

**Q: What video formats are supported?**
A: Sora generates MP4 videos.

**Q: Can I generate longer videos?**
A: Default is 4 seconds. Check OpenAI docs for maximum duration.

**Q: Why is my reference image ignored?**
A: Ensure the image is:
- A valid JPEG/PNG file
- Under 20MB
- Accessible (correct path)

---

**Still having issues?** Contact support or open a GitHub issue with details.
