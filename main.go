package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/telemetry/video-gen/internal/cli"
	"github.com/telemetry/video-gen/internal/tui"
)

func main() {
	// CLI flags
	debug := flag.Bool("d", false, "Enable debug mode (show API requests/responses)")
	prompt := flag.String("p", "", "Video generation prompt (triggers non-interactive mode)")
	model := flag.String("m", "", "Model: 'sora' or 'sora-pro'")
	referenceImage := flag.String("r", "", "Path to reference image")
	duration := flag.String("t", "", "Duration: 4, 8, or 12 seconds")
	size := flag.String("s", "", "Size: '1280x720', '720x1280', '1792x1024', or '1024x1792'")
	outputDir := flag.String("o", "", "Output directory")

	flag.Parse()

	// If prompt is provided via -p flag, run in non-interactive CLI mode
	if *prompt != "" {
		opts := cli.Options{
			Debug:          *debug,
			Prompt:         *prompt,
			Model:          *model,
			ReferenceImage: *referenceImage,
			Duration:       *duration,
			Size:           *size,
			OutputDir:      *outputDir,
		}

		if err := cli.RunNonInteractive(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Otherwise run interactive TUI mode
	opts := tui.CLIOptions{
		Debug:          *debug,
		Prompt:         *prompt,
		Model:          *model,
		ReferenceImage: *referenceImage,
		Duration:       *duration,
		Size:           *size,
		OutputDir:      *outputDir,
	}

	tuiModel, err := tui.NewModel(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tuiModel)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
