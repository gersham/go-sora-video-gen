package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/telemetry/video-gen/internal/api"
	"github.com/telemetry/video-gen/internal/config"
)

type state int

const (
	stateAPIKey state = iota
	stateListVideos
	stateDeletingVideos
	statePrompt
	stateModel
	stateReferenceImage
	stateDuration
	stateSize
	stateOutputDir
	stateGenerating
	statePolling
	stateDownloading
	stateComplete
	stateError
)

type videoCreatedMsg struct {
	id string
}

type videoReadyMsg struct {
	videoID string // Changed to use video ID instead of URL
}

type videoDownloadedMsg struct {
	path string
}

type errorMsg struct {
	err error
}

type pollMsg struct {
	progress int    // Progress percentage from API
	status   string // Status from API
}

type debugMsg struct {
	entry string
}

type videosListedMsg struct {
	videos []api.VideoResponse
}

type videoDeletedMsg struct {
	videoID string
	current int
	total   int
}

type videosDeletedMsg struct{}

type tickMsg time.Time

type Model struct {
	state          state
	textInput      textinput.Model
	spinner        spinner.Model
	cfg            *config.Config
	client         *api.SoraClient
	prompt         string
	model          string
	modelSelection int // 0 = sora-2, 1 = sora-2-pro
	referenceImg   string
	duration          string
	durationSelection int // 0 = 4s, 1 = 8s, 2 = 12s
	size              string
	sizeSelection     int // 0 = 1280x720, 1 = 720x1280, 2 = 1792x1024, 3 = 1024x1792
	outputDir      string
	videoID        string
	outputPath     string
	err            error
	message        string
	pollAttempts   int
	elapsedSeconds int
	progress       int    // Video generation progress percentage (0-100)
	videoStatus    string // Current video status from API
	skipReference  bool
	debug          bool
	debugLogs           []string
	recentVideos        []api.VideoResponse
	deleteVideos        bool // Whether to delete listed videos
	deletingVideoID     string
	deletingVideoIndex  int
	deletingVideoTotal  int
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	debugRequestStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Bold(true)

	debugResponseStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("213")).
				Bold(true)

	debugJSONStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))
)

// CLIOptions holds command-line options
type CLIOptions struct {
	Debug          bool
	Prompt         string
	Model          string
	ReferenceImage string
	Duration       string
	Size           string
	OutputDir      string
}

func NewModel(opts CLIOptions) (*Model, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := &Model{
		textInput: ti,
		spinner:   s,
		cfg:       cfg,
		debug:     opts.Debug,
		debugLogs: make([]string, 0),
	}

	// Check API key first
	if cfg.OpenAIAPIKey == "" {
		m.state = stateAPIKey
		m.textInput.Placeholder = "sk-..."
		return m, nil
	}

	// Create debug callback that appends directly to the slice
	debugCallback := func(entry string) {
		if m.debug {
			m.debugLogs = append(m.debugLogs, entry)
			if len(m.debugLogs) > 50 {
				m.debugLogs = m.debugLogs[len(m.debugLogs)-50:]
			}
		}
	}
	m.client = api.NewClient(cfg.OpenAIAPIKey, m.debug, debugCallback)

	// Determine initial state based on CLI options
	if opts.Prompt != "" {
		// CLI mode: all required params provided, start generation
		m.prompt = opts.Prompt
		m.state = stateGenerating
	} else {
		// Interactive mode: start by listing recent videos
		m.state = stateListVideos
		m.deleteVideos = true // Default to yes for deletion
		m.textInput.Placeholder = ""
	}

	// Apply CLI options or fall back to config/defaults
	// Output directory
	if opts.OutputDir != "" {
		m.outputDir = opts.OutputDir
	} else if cfg.OutputDir != "" {
		m.outputDir = cfg.OutputDir
	} else {
		homeDir, _ := os.UserHomeDir()
		m.outputDir = filepath.Join(homeDir, "Desktop")
	}

	// Model
	if opts.Model != "" {
		modelName := opts.Model
		if modelName == "sora" {
			modelName = "sora-2"
		} else if modelName == "sora-pro" {
			modelName = "sora-2-pro"
		}
		m.model = modelName
		if modelName == "sora-2" {
			m.modelSelection = 0
		} else {
			m.modelSelection = 1
		}
	} else if cfg.Model != "" {
		m.model = cfg.Model
		if cfg.Model == "sora-2" {
			m.modelSelection = 0
		} else {
			m.modelSelection = 1
		}
	} else {
		m.model = "sora-2"
		m.modelSelection = 0
	}

	// Duration
	if opts.Duration != "" {
		m.duration = opts.Duration
		m.durationSelection = getDurationSelection(opts.Duration)
	} else if cfg.Duration != "" {
		m.duration = cfg.Duration
		m.durationSelection = getDurationSelection(cfg.Duration)
	} else {
		m.duration = "4"
		m.durationSelection = 0
	}

	// Size
	if opts.Size != "" {
		m.size = opts.Size
		m.sizeSelection = getSizeSelection(opts.Size)
	} else if cfg.Size != "" {
		m.size = cfg.Size
		m.sizeSelection = getSizeSelection(cfg.Size)
	} else {
		m.size = "1280x720"
		m.sizeSelection = 0
	}

	// Reference image
	if opts.ReferenceImage != "" {
		m.referenceImg = opts.ReferenceImage
	}

	return m, nil
}

// Helper function to get size selection index
func getDurationSelection(duration string) int {
	switch duration {
	case "4":
		return 0
	case "8":
		return 1
	case "12":
		return 2
	default:
		return 0
	}
}

func getSizeSelection(size string) int {
	switch size {
	case "1280x720":
		return 0
	case "720x1280":
		return 1
	case "1792x1024":
		return 2
	case "1024x1792":
		return 3
	default:
		return 0
	}
}

func (m *Model) addDebugLog(entry string) {
	if m.debug {
		m.debugLogs = append(m.debugLogs, entry)
		// Keep last 50 entries
		if len(m.debugLogs) > 50 {
			m.debugLogs = m.debugLogs[len(m.debugLogs)-50:]
		}
	}
}

func (m Model) Init() tea.Cmd {
	// Clear screen on startup
	clearScreen := func() tea.Msg {
		return tea.ClearScreen()
	}

	// If we're in CLI mode (generating state), start immediately
	if m.state == stateGenerating {
		return tea.Batch(clearScreen, textinput.Blink, m.spinner.Tick, m.createVideo(), tick())
	}
	// If in interactive mode, list recent videos
	if m.state == stateListVideos {
		return tea.Batch(clearScreen, textinput.Blink, m.spinner.Tick, m.listVideos())
	}
	return tea.Batch(clearScreen, textinput.Blink, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		// Continue ticking during deleting state
		if m.state == stateDeletingVideos {
			return m, tea.Batch(cmd, m.spinner.Tick)
		}
		return m, cmd

	case tickMsg:
		if m.state == statePolling || m.state == stateGenerating {
			m.elapsedSeconds++
			return m, tick()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlU:
			// Clear the input field
			m.textInput.SetValue("")
			return m, nil

		case tea.KeyEnter:
			if m.state == stateListVideos {
				// User confirmed deletion choice
				if m.deleteVideos && len(m.recentVideos) > 0 {
					// Transition to deleting state
					m.state = stateDeletingVideos
					return m, tea.Batch(m.deleteAllVideos(), m.spinner.Tick)
				} else {
					// Skip deletion, go to prompt
					m.state = statePrompt
					m.textInput.SetValue(m.cfg.LastPrompt)
					m.textInput.Placeholder = "Describe the video you want to generate..."
					m.textInput.Focus()
					return m, nil
				}
			}
			if m.state == stateComplete {
				// Restart after completion - preserve prompt and reference image
				previousPrompt := m.prompt
				m.state = statePrompt
				m.videoID = ""
				m.outputPath = ""
				m.err = nil
				m.message = ""
				m.pollAttempts = 0
				m.elapsedSeconds = 0
				m.progress = 0
				m.skipReference = false
				// Keep referenceImg set so it becomes the default
				m.textInput.SetValue(previousPrompt)
				m.textInput.Placeholder = "Describe the video you want to generate..."
				m.textInput.Focus()
				return m, nil
			}
			if m.state == stateError {
				// Retry after error - preserve prompt and allow editing
				previousPrompt := m.prompt
				m.state = statePrompt
				m.videoID = ""
				m.outputPath = ""
				m.err = nil
				m.message = ""
				m.pollAttempts = 0
				m.elapsedSeconds = 0
				m.progress = 0
				m.skipReference = false
				// Pre-fill with previous prompt for easy editing
				m.textInput.SetValue(previousPrompt)
				m.textInput.Placeholder = "Describe the video you want to generate..."
				m.textInput.Focus()
				return m, nil
			}
			if m.state == stateModel {
				// Handle model selection with Enter
				if m.modelSelection == 0 {
					m.model = "sora-2"
				} else {
					m.model = "sora-2-pro"
				}
				m.cfg.Model = m.model
				m.state = stateReferenceImage
				// Set previous reference image as default (if it exists)
				m.textInput.SetValue(m.referenceImg)
				m.textInput.Placeholder = "Path to reference image (or press Enter to skip)..."
				m.message = ""
				return m, nil
			}
			if m.state == stateSize {
				// Handle size selection with Enter
				sizes := []string{"1280x720", "720x1280", "1792x1024", "1024x1792"}
				m.size = sizes[m.sizeSelection]
				m.cfg.Size = m.size
				m.state = stateOutputDir
				m.textInput.SetValue(m.outputDir)
				m.textInput.Placeholder = "Output directory..."
				m.message = ""
				return m, nil
			}
			return m.handleEnter()

		case tea.KeyUp, tea.KeyLeft:
			if m.state == stateListVideos {
				m.deleteVideos = !m.deleteVideos
				return m, nil
			}
			if m.state == stateModel {
				m.modelSelection = (m.modelSelection - 1 + 2) % 2
				return m, nil
			}
			if m.state == stateDuration {
				m.durationSelection = (m.durationSelection - 1 + 3) % 3
				return m, nil
			}
			if m.state == stateSize {
				m.sizeSelection = (m.sizeSelection - 1 + 4) % 4
				return m, nil
			}

		case tea.KeyDown, tea.KeyRight:
			if m.state == stateListVideos {
				m.deleteVideos = !m.deleteVideos
				return m, nil
			}
			if m.state == stateModel {
				m.modelSelection = (m.modelSelection + 1) % 2
				return m, nil
			}
			if m.state == stateDuration {
				m.durationSelection = (m.durationSelection + 1) % 3
				return m, nil
			}
			if m.state == stateSize {
				m.sizeSelection = (m.sizeSelection + 1) % 4
				return m, nil
			}
		}

	case videoCreatedMsg:
		m.videoID = msg.id
		m.state = statePolling
		m.pollAttempts = 0
		m.elapsedSeconds = 0
		m.progress = 0
		return m, tea.Batch(m.checkVideoStatus(), tick())

	case pollMsg:
		if m.state != statePolling {
			return m, nil
		}
		m.pollAttempts++
		m.progress = msg.progress   // Update progress from API
		m.videoStatus = msg.status  // Update status from API
		if m.pollAttempts > 200 {
			return m, func() tea.Msg {
				return errorMsg{err: fmt.Errorf("timeout waiting for video generation")}
			}
		}
		return m, m.pollVideo()

	case videoReadyMsg:
		m.state = stateDownloading
		return m, m.downloadVideo()

	case videoDownloadedMsg:
		m.outputPath = msg.path
		m.state = stateComplete
		return m, nil

	case videosListedMsg:
		m.recentVideos = msg.videos
		// Stay in stateListVideos to show the list
		return m, nil

	case videoDeletedMsg:
		m.deletingVideoID = msg.videoID
		m.deletingVideoIndex = msg.current
		m.deletingVideoTotal = msg.total
		return m, nil

	case videosDeletedMsg:
		m.recentVideos = nil
		m.deletingVideoID = ""
		m.deletingVideoIndex = 0
		m.deletingVideoTotal = 0
		m.state = statePrompt
		m.textInput.SetValue(m.cfg.LastPrompt)
		m.textInput.Placeholder = "Describe the video you want to generate..."
		m.textInput.Focus()
		return m, nil

	case errorMsg:
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	value := strings.TrimSpace(m.textInput.Value())

	switch m.state {
	case stateAPIKey:
		if value == "" {
			m.message = "API key cannot be empty"
			return m, nil
		}
		m.cfg.OpenAIAPIKey = value
		if err := config.Save(m.cfg); err != nil {
			m.err = err
			m.state = stateError
			return m, nil
		}
		// Create debug callback that appends directly to the slice
		debugCallback := func(entry string) {
			if m.debug {
				m.debugLogs = append(m.debugLogs, entry)
				if len(m.debugLogs) > 50 {
					m.debugLogs = m.debugLogs[len(m.debugLogs)-50:]
				}
			}
		}
		m.client = api.NewClient(value, m.debug, debugCallback)
		m.state = statePrompt
		m.textInput.SetValue("")
		m.textInput.Placeholder = "Describe the video you want to generate..."
		m.message = ""
		return m, nil

	case statePrompt:
		if value == "" {
			// Empty prompt means exit
			return m, tea.Quit
		}
		m.prompt = value
		m.cfg.LastPrompt = value
		m.state = stateModel
		// Model selection is now handled by arrow keys, not text input
		m.message = ""
		return m, nil

	case stateReferenceImage:
		if value != "" {
			// Expand tilde to home directory
			if strings.HasPrefix(value, "~/") {
				homeDir, err := os.UserHomeDir()
				if err == nil {
					value = filepath.Join(homeDir, value[2:])
				}
			}
			// Validate file exists
			if _, err := os.Stat(value); os.IsNotExist(err) {
				m.message = "File does not exist"
				return m, nil
			}
			m.referenceImg = value
		} else {
			m.skipReference = true
		}
		m.state = stateDuration
		m.textInput.SetValue(m.duration)
		m.textInput.Placeholder = m.duration
		m.message = ""
		return m, nil

	case stateDuration:
		// Duration selection is confirmed, save and move to size
		durations := []string{"4", "8", "12"}
		m.duration = durations[m.durationSelection]
		m.cfg.Duration = m.duration
		m.state = stateSize
		// Size selection is handled by arrow keys, not text input
		m.message = ""
		return m, nil

	case stateOutputDir:
		if value != "" {
			m.outputDir = value
		}
		m.cfg.OutputDir = m.outputDir
		// Save config with all updates
		if err := config.Save(m.cfg); err != nil {
			m.err = fmt.Errorf("failed to save config: %w", err)
			m.state = stateError
			return m, nil
		}
		m.state = stateGenerating
		return m, m.createVideo()
	}

	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) createVideo() tea.Cmd {
	return func() tea.Msg {
		req := api.CreateVideoRequest{
			Prompt:         m.prompt,
			Model:          m.model,
			InputReference: m.referenceImg,
			Seconds:        m.duration,
			Size:           m.size,
		}

		resp, err := m.client.CreateVideo(req)
		if err != nil {
			return errorMsg{err: err}
		}

		return videoCreatedMsg{id: resp.ID}
	}
}

func (m Model) pollVideo() tea.Cmd {
	return func() tea.Msg {
		// Dynamic polling: 10s for first 2 minutes, 10s when at 100%, 30s thereafter
		var pollInterval time.Duration
		if m.progress >= 100 {
			// Poll every 10s when at 100% waiting for completion
			pollInterval = 10 * time.Second
		} else if m.elapsedSeconds < 120 {
			pollInterval = 10 * time.Second
		} else {
			pollInterval = 30 * time.Second
		}
		time.Sleep(pollInterval)

		// Check video status after sleep
		resp, err := m.client.GetVideo(m.videoID)
		if err != nil {
			return errorMsg{err: err}
		}

		// Only download when status is "completed"
		if resp.Status == "completed" {
			return videoReadyMsg{videoID: m.videoID}
		}

		if resp.Status == "failed" {
			errMsg := "Video generation failed"
			if resp.Error != nil && resp.Error.Message != "" {
				errMsg += ": " + resp.Error.Message
			}
			return errorMsg{err: fmt.Errorf(errMsg)}
		}

		// Continue polling with progress and status update
		return pollMsg{progress: resp.Progress, status: resp.Status}
	}
}

func (m Model) checkVideoStatus() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.GetVideo(m.videoID)
		if err != nil {
			return errorMsg{err: err}
		}

		// Only download when status is "completed"
		if resp.Status == "completed" {
			return videoReadyMsg{videoID: m.videoID}
		}

		if resp.Status == "failed" {
			errMsg := "Video generation failed"
			if resp.Error != nil && resp.Error.Message != "" {
				errMsg += ": " + resp.Error.Message
			}
			return errorMsg{err: fmt.Errorf(errMsg)}
		}

		// Continue polling with progress and status update
		return pollMsg{progress: resp.Progress, status: resp.Status}
	}
}

func (m Model) listVideos() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.ListVideos(10)
		if err != nil {
			return errorMsg{err: err}
		}
		return videosListedMsg{videos: resp.Data}
	}
}

func (m Model) deleteAllVideos() tea.Cmd {
	videos := m.recentVideos

	return func() tea.Msg {
		// Delete all videos
		for _, video := range videos {
			// Ignore errors and continue
			_ = m.client.DeleteVideo(video.ID)
		}

		// All done
		return videosDeletedMsg{}
	}
}

func (m Model) downloadVideo() tea.Cmd {
	return func() tea.Msg {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("sora_video_%s.mp4", timestamp)
		outputPath := filepath.Join(m.outputDir, filename)

		// Retry download up to 12 times (2 minutes with 10s intervals)
		maxRetries := 12
		for attempt := 0; attempt < maxRetries; attempt++ {
			if attempt > 0 {
				time.Sleep(10 * time.Second)
			}

			err := m.client.DownloadVideoContent(m.videoID, outputPath)
			if err == nil {
				// Download successful, now delete the video from the service
				if deleteErr := m.client.DeleteVideo(m.videoID); deleteErr != nil {
					// Log error but don't fail the operation since download succeeded
					// The video will remain on the service but user has their file
					fmt.Fprintf(os.Stderr, "Warning: failed to delete video from service: %v\n", deleteErr)
				}
				return videoDownloadedMsg{path: outputPath}
			}

			// Check if it's a 404 (not ready yet) - if so, retry
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not ready") {
				continue
			}

			// Other errors, fail immediately
			return errorMsg{err: err}
		}

		return errorMsg{err: fmt.Errorf("video content not available after %d attempts (2 minutes)", maxRetries)}
	}
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("TelemetryOS Video Generator (Sora)"))
	sb.WriteString("\n\n")

	switch m.state {
	case stateAPIKey:
		sb.WriteString(promptStyle.Render("Enter your OpenAI API key:"))
		sb.WriteString("\n")
		sb.WriteString(m.textInput.View())
		if m.message != "" {
			sb.WriteString("\n")
			sb.WriteString(errorStyle.Render(m.message))
		}

	case stateListVideos:
		if m.recentVideos == nil {
			sb.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), infoStyle.Render("Loading recent videos...")))
		} else if len(m.recentVideos) == 0 {
			sb.WriteString(promptStyle.Render("No recent videos found."))
			sb.WriteString("\n\n")
			sb.WriteString(promptStyle.Render("Press Enter to continue..."))
		} else {
			sb.WriteString(promptStyle.Render(fmt.Sprintf("Recent videos (%d found):", len(m.recentVideos))))
			sb.WriteString("\n\n")

			for i, video := range m.recentVideos {
				if i >= 10 {
					break
				}
				createdTime := time.Unix(video.CreatedAt, 0).Format("Jan 2, 15:04")
				statusColor := promptStyle
				if video.Status == "completed" {
					statusColor = successStyle
				} else if video.Status == "failed" {
					statusColor = errorStyle
				}
				sb.WriteString(fmt.Sprintf("  %s - %s (%s) - %s\n",
					promptStyle.Render(video.ID[:20]+"..."),
					statusColor.Render(video.Status),
					infoStyle.Render(video.Model),
					promptStyle.Render(createdTime)))
			}

			sb.WriteString("\n")
			sb.WriteString(promptStyle.Render("Delete all listed videos? (use arrow keys to toggle)"))
			sb.WriteString("\n")

			if m.deleteVideos {
				sb.WriteString(successStyle.Render("▶ Yes"))
				sb.WriteString("  ")
				sb.WriteString(promptStyle.Render("No"))
			} else {
				sb.WriteString(promptStyle.Render("  Yes"))
				sb.WriteString("  ")
				sb.WriteString(successStyle.Render("▶ No"))
			}

			sb.WriteString("\n\n")
			sb.WriteString(promptStyle.Render("Press Enter to confirm"))
		}

	case stateDeletingVideos:
		sb.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), infoStyle.Render(fmt.Sprintf("Deleting %d videos...", len(m.recentVideos)))))
		sb.WriteString("\n")
		sb.WriteString(promptStyle.Render("This may take a moment..."))

	case statePrompt:
		sb.WriteString(promptStyle.Render("Enter video generation prompt:"))
		sb.WriteString("\n")
		sb.WriteString(m.textInput.View())
		if m.message != "" {
			sb.WriteString("\n")
			sb.WriteString(errorStyle.Render(m.message))
		}

	case stateModel:
		sb.WriteString(promptStyle.Render("Select model (use arrow keys):"))
		sb.WriteString("\n\n")

		// Option 1: sora-2
		if m.modelSelection == 0 {
			sb.WriteString(successStyle.Render("▶ sora-2"))
		} else {
			sb.WriteString(promptStyle.Render("  sora-2"))
		}
		sb.WriteString(promptStyle.Render("       - Fast generation, good quality"))
		sb.WriteString("\n")

		// Option 2: sora-2-pro
		if m.modelSelection == 1 {
			sb.WriteString(successStyle.Render("▶ sora-2-pro"))
		} else {
			sb.WriteString(promptStyle.Render("  sora-2-pro"))
		}
		sb.WriteString(promptStyle.Render("   - Superior quality, slower"))
		sb.WriteString("\n\n")
		sb.WriteString(promptStyle.Render("Press Enter to confirm"))
		if m.message != "" {
			sb.WriteString("\n")
			sb.WriteString(errorStyle.Render(m.message))
		}

	case stateReferenceImage:
		sb.WriteString(promptStyle.Render("Reference image path (optional):"))
		sb.WriteString("\n")
		sb.WriteString(m.textInput.View())
		if m.message != "" {
			sb.WriteString("\n")
			sb.WriteString(errorStyle.Render(m.message))
		}

	case stateDuration:
		sb.WriteString(promptStyle.Render("Select video duration (use arrow keys):"))
		sb.WriteString("\n\n")

		durations := []struct {
			duration string
			desc     string
		}{
			{"4", "4 seconds"},
			{"8", "8 seconds"},
			{"12", "12 seconds"},
		}

		for i, dur := range durations {
			if i == m.durationSelection {
				sb.WriteString(successStyle.Render(fmt.Sprintf("→ %s - %s", dur.duration, dur.desc)))
			} else {
				sb.WriteString(fmt.Sprintf("  %s - %s", dur.duration, dur.desc))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("\n")
		sb.WriteString(promptStyle.Render("Press Enter to confirm"))

	case stateSize:
		sb.WriteString(promptStyle.Render("Select video size (use arrow keys):"))
		sb.WriteString("\n\n")

		sizes := []struct {
			size string
			desc string
		}{
			{"1280x720", "Landscape (HD)"},
			{"720x1280", "Portrait (HD)"},
			{"1792x1024", "Landscape (Wide)"},
			{"1024x1792", "Portrait (Wide)"},
		}

		for i, s := range sizes {
			if m.sizeSelection == i {
				sb.WriteString(successStyle.Render("▶ " + s.size))
			} else {
				sb.WriteString(promptStyle.Render("  " + s.size))
			}
			sb.WriteString(promptStyle.Render("   - " + s.desc))
			sb.WriteString("\n")
		}

		sb.WriteString("\n")
		sb.WriteString(promptStyle.Render("Press Enter to confirm"))
		if m.message != "" {
			sb.WriteString("\n")
			sb.WriteString(errorStyle.Render(m.message))
		}

	case stateOutputDir:
		sb.WriteString(promptStyle.Render("Output directory:"))
		sb.WriteString("\n")
		sb.WriteString(m.textInput.View())

	case stateGenerating:
		sb.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), infoStyle.Render(fmt.Sprintf("Creating video generation job... (%ds)", m.elapsedSeconds))))
		sb.WriteString("\n")
		sb.WriteString(promptStyle.Render("This may take a moment. Retrying automatically if needed..."))

	case statePolling:
		// Display status after time: "Generating video (17s) queued"
		progressStr := ""
		if m.progress > 0 {
			progressStr = fmt.Sprintf(" (%d%% complete)", m.progress)
		}
		statusDisplay := "unknown"
		if m.videoStatus != "" {
			statusDisplay = m.videoStatus
		}
		sb.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), infoStyle.Render(fmt.Sprintf("Generating video (%ds) %s%s", m.elapsedSeconds, statusDisplay, progressStr))))
		sb.WriteString("\n")
		pollInterval := "10s"
		if m.elapsedSeconds >= 120 {
			pollInterval = "30s"
		}
		sb.WriteString(promptStyle.Render(fmt.Sprintf("Polling API every %s (attempt %d/200)", pollInterval, m.pollAttempts)))

	case stateDownloading:
		sb.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), infoStyle.Render("Downloading video...")))

	case stateComplete:
		sb.WriteString(successStyle.Render("✓ Video generated successfully!"))
		sb.WriteString("\n\n")
		sb.WriteString(infoStyle.Render(fmt.Sprintf("Saved to: %s", m.outputPath)))
		sb.WriteString("\n\n")
		sb.WriteString(promptStyle.Render("Press Enter to generate another video..."))

	case stateError:
		sb.WriteString(errorStyle.Render("✗ Error occurred:"))
		sb.WriteString("\n")
		sb.WriteString(errorStyle.Render(m.err.Error()))
		sb.WriteString("\n\n")
		sb.WriteString(promptStyle.Render("Press Enter to try again with a different prompt..."))
	}

	sb.WriteString("\n\n")
	sb.WriteString(promptStyle.Render("Press Ctrl+C to quit"))

	// Debug logs at the bottom
	if m.debug && len(m.debugLogs) > 0 {
		sb.WriteString("\n\n")
		sb.WriteString(strings.Repeat("─", 80))
		sb.WriteString("\n")
		sb.WriteString(debugRequestStyle.Render("DEBUG MODE"))
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("─", 80))
		sb.WriteString("\n\n")

		// Show last 10 log entries
		start := 0
		if len(m.debugLogs) > 10 {
			start = len(m.debugLogs) - 10
		}

		for i := start; i < len(m.debugLogs); i++ {
			entry := m.debugLogs[i]
			if strings.HasPrefix(entry, "REQUEST:") {
				sb.WriteString(debugRequestStyle.Render("→ "))
				sb.WriteString(debugJSONStyle.Render(entry))
			} else {
				sb.WriteString(debugResponseStyle.Render("← "))
				sb.WriteString(debugJSONStyle.Render(entry))
			}
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}
