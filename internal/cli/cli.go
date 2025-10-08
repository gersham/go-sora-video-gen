package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/telemetry/video-gen/internal/api"
	"github.com/telemetry/video-gen/internal/config"
)

type Options struct {
	Debug          bool
	Prompt         string
	Model          string
	ReferenceImage string
	Duration       string
	Size           string
	OutputDir      string
}

// RunNonInteractive runs the video generation in non-interactive mode
func RunNonInteractive(opts Options) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check API key
	if cfg.OpenAIAPIKey == "" {
		return fmt.Errorf("OpenAI API key not found. Please run interactively first or set key in config")
	}

	// Set defaults from config
	model := opts.Model
	if model == "" {
		if cfg.Model != "" {
			model = cfg.Model
		} else {
			model = "sora-2"
		}
	} else {
		// Normalize model name
		if model == "sora" {
			model = "sora-2"
		} else if model == "sora-pro" {
			model = "sora-2-pro"
		}
	}

	duration := opts.Duration
	if duration == "" {
		if cfg.Duration != "" {
			duration = cfg.Duration
		} else {
			duration = "4"
		}
	}
	// Validate duration (must be 4, 8, or 12)
	if duration != "4" && duration != "8" && duration != "12" {
		return fmt.Errorf("invalid duration '%s'. Supported values are: '4', '8', and '12'", duration)
	}

	size := opts.Size
	if size == "" {
		if cfg.Size != "" {
			size = cfg.Size
		} else {
			size = "1280x720"
		}
	}

	outputDir := opts.OutputDir
	if outputDir == "" {
		if cfg.OutputDir != "" {
			outputDir = cfg.OutputDir
		} else {
			homeDir, _ := os.UserHomeDir()
			outputDir = filepath.Join(homeDir, "Desktop")
		}
	}

	// Expand tilde in reference image path
	referenceImage := opts.ReferenceImage
	if referenceImage != "" && strings.HasPrefix(referenceImage, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			referenceImage = filepath.Join(homeDir, referenceImage[2:])
		}
	}

	// Create debug callback
	debugCallback := func(entry string) {
		if opts.Debug {
			fmt.Println(entry)
		}
	}

	// Create API client
	client := api.NewClient(cfg.OpenAIAPIKey, opts.Debug, debugCallback)

	// Step 1: Create video
	fmt.Printf("Creating video generation job...\n")
	fmt.Printf("  Prompt: %s\n", opts.Prompt)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Duration: %ss\n", duration)
	fmt.Printf("  Size: %s\n", size)
	if referenceImage != "" {
		fmt.Printf("  Reference: %s\n", referenceImage)
	}
	fmt.Println()

	createReq := api.CreateVideoRequest{
		Prompt:         opts.Prompt,
		Model:          model,
		InputReference: referenceImage,
		Seconds:        duration,
		Size:           size,
	}

	createResp, err := client.CreateVideo(createReq)
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}

	fmt.Printf("✓ Video job created: %s\n", createResp.ID)
	fmt.Println()

	// Step 2: Poll for completion
	videoID := createResp.ID
	pollAttempts := 0
	maxAttempts := 200
	startTime := time.Now()

	fmt.Println("Polling for completion...")
	fmt.Println("(This may take several minutes)")
	fmt.Println()

	for pollAttempts < maxAttempts {
		pollAttempts++
		elapsed := int(time.Since(startTime).Seconds())

		// Determine poll interval: 10s for first 2 minutes, 30s thereafter
		var pollInterval time.Duration
		if pollAttempts == 1 {
			// First check is immediate
			pollInterval = 0
		} else if elapsed < 120 {
			pollInterval = 10 * time.Second
		} else {
			pollInterval = 30 * time.Second
		}

		if pollInterval > 0 {
			time.Sleep(pollInterval)
		}

		resp, err := client.GetVideo(videoID)
		if err != nil {
			return fmt.Errorf("failed to get video status: %w", err)
		}

		elapsed = int(time.Since(startTime).Seconds())
		progressStr := ""
		if resp.Progress > 0 {
			progressStr = fmt.Sprintf(" (%d%% complete)", resp.Progress)
		}

		fmt.Printf("[%ds] Status: %s%s (attempt %d/%d)\n", elapsed, resp.Status, progressStr, pollAttempts, maxAttempts)

		// Only download when status is "completed"
		if resp.Status == "completed" {
			fmt.Println()
			fmt.Printf("✓ Video generation completed!\n")
			fmt.Println()

			// Step 3: Download video content directly
			timestamp := time.Now().Format("20060102_150405")
			filename := fmt.Sprintf("sora_video_%s.mp4", timestamp)
			outputPath := filepath.Join(outputDir, filename)

			fmt.Printf("Downloading video to: %s\n", outputPath)

			// Retry download with 10s intervals (up to 12 attempts = 2 minutes)
			maxDownloadRetries := 12
			var downloadErr error
			for downloadAttempt := 0; downloadAttempt < maxDownloadRetries; downloadAttempt++ {
				if downloadAttempt > 0 {
					fmt.Printf("  Retrying download (attempt %d/%d)...\n", downloadAttempt+1, maxDownloadRetries)
					time.Sleep(10 * time.Second)
				}

				downloadErr = client.DownloadVideoContent(videoID, outputPath)
				if downloadErr == nil {
					break // Success!
				}

				// Check if it's a 404 (not ready yet) - if so, retry
				if !strings.Contains(downloadErr.Error(), "404") && !strings.Contains(downloadErr.Error(), "not ready") {
					// Other errors, fail immediately
					return fmt.Errorf("failed to download video: %w", downloadErr)
				}
			}

			if downloadErr != nil {
				return fmt.Errorf("video content not available after %d attempts (2 minutes): %w", maxDownloadRetries, downloadErr)
			}

			fmt.Println()
			fmt.Printf("✓ Video saved successfully!\n")
			fmt.Printf("  Location: %s\n", outputPath)

			// Delete the video from the service after successful download
			fmt.Println()
			fmt.Printf("Deleting video from service...\n")
			if err := client.DeleteVideo(videoID); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to delete video from service: %v\n", err)
			} else {
				fmt.Printf("✓ Video deleted from service\n")
			}

			return nil
		}

		if resp.Status == "failed" {
			errMsg := "Video generation failed"
			if resp.Error != nil && resp.Error.Message != "" {
				errMsg += ": " + resp.Error.Message
			}
			return fmt.Errorf(errMsg)
		}

	}

	return fmt.Errorf("timeout waiting for video generation")
}
