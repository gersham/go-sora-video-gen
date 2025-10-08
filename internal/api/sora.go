package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	baseURL        = "https://api.openai.com/v1"
	createEndpoint = "/videos"
)

type SoraClient struct {
	apiKey     string
	httpClient *http.Client
	debug      bool
	debugLog   func(string)
}

type CreateVideoRequest struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model,omitempty"`
	Seconds        string `json:"seconds,omitempty"`
	Size           string `json:"size,omitempty"`
	InputReference string `json:"-"` // File path, handled separately
}

type CreateVideoResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Object string `json:"object"`
}

type ErrorObject struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

type VideoResponse struct {
	ID                 string       `json:"id"`
	Status             string       `json:"status"`
	Error              *ErrorObject `json:"error,omitempty"`
	CreatedAt          int64        `json:"created_at"`
	CompletedAt        int64        `json:"completed_at,omitempty"`
	ExpiresAt          int64        `json:"expires_at,omitempty"`
	Progress           int          `json:"progress,omitempty"`
	Model              string       `json:"model,omitempty"`
	Seconds            string       `json:"seconds,omitempty"`
	Size               string       `json:"size,omitempty"`
	Object             string       `json:"object,omitempty"`
	RemixedFromVideoID string       `json:"remixed_from_video_id,omitempty"`
}

type ListVideosResponse struct {
	Data   []VideoResponse `json:"data"`
	Object string          `json:"object"`
}

type APIError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func NewClient(apiKey string, debug bool, debugLog func(string)) *SoraClient {
	return &SoraClient{
		apiKey:   apiKey,
		debug:    debug,
		debugLog: debugLog,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// CreateVideo initiates video generation with the Sora API with retry logic
func (c *SoraClient) CreateVideo(req CreateVideoRequest) (*CreateVideoResponse, error) {
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 2s, 4s, 8s
			waitTime := time.Duration(1<<uint(attempt)) * time.Second
			time.Sleep(waitTime)
		}

		result, err := c.createVideoAttempt(req)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Don't retry on authentication or validation errors
		if isClientError(err) {
			break
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

func (c *SoraClient) createVideoAttempt(req CreateVideoRequest) (*CreateVideoResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add text fields
	if err := writer.WriteField("prompt", req.Prompt); err != nil {
		return nil, fmt.Errorf("failed to write prompt: %w", err)
	}

	if req.Model != "" {
		if err := writer.WriteField("model", req.Model); err != nil {
			return nil, fmt.Errorf("failed to write model: %w", err)
		}
	}

	if req.Seconds != "" {
		if err := writer.WriteField("seconds", req.Seconds); err != nil {
			return nil, fmt.Errorf("failed to write seconds: %w", err)
		}
	}

	if req.Size != "" {
		if err := writer.WriteField("size", req.Size); err != nil {
			return nil, fmt.Errorf("failed to write size: %w", err)
		}
	}

	// Add reference file if provided
	if req.InputReference != "" {
		file, err := os.Open(req.InputReference)
		if err != nil {
			return nil, fmt.Errorf("failed to open reference file: %w", err)
		}
		defer file.Close()

		// Decode image
		img, format, err := image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %w", err)
		}

		// Parse target dimensions from size string (e.g., "1280x720")
		targetWidth, targetHeight, err := parseSize(req.Size)
		if err != nil {
			return nil, fmt.Errorf("invalid size format: %w", err)
		}

		// Resize and crop image to match target dimensions
		img = resizeAndCropToFill(img, targetWidth, targetHeight)

		// Detect MIME type from format
		filename := filepath.Base(req.InputReference)
		contentType := "application/octet-stream"
		switch format {
		case "jpeg":
			contentType = "image/jpeg"
		case "png":
			contentType = "image/png"
		case "gif":
			contentType = "image/gif"
		}

		// Create form file with proper Content-Type header
		h := make(map[string][]string)
		h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="input_reference"; filename="%s"`, filename)}
		h["Content-Type"] = []string{contentType}
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		// Encode resized image to part
		if format == "png" {
			if err := png.Encode(part, img); err != nil {
				return nil, fmt.Errorf("failed to encode PNG: %w", err)
			}
		} else {
			// Default to JPEG for other formats
			if err := jpeg.Encode(part, img, &jpeg.Options{Quality: 95}); err != nil {
				return nil, fmt.Errorf("failed to encode JPEG: %w", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", baseURL+createEndpoint, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// Debug log request
	if c.debug && c.debugLog != nil {
		reqJSON, _ := json.MarshalIndent(map[string]interface{}{
			"method":  "POST",
			"url":     baseURL + createEndpoint,
			"headers": map[string]string{"Content-Type": writer.FormDataContentType()},
			"body": map[string]string{
				"prompt": req.Prompt,
				"model":  req.Model,
				"seconds": req.Seconds,
				"size": req.Size,
			},
		}, "", "  ")
		c.debugLog(fmt.Sprintf("REQUEST:\n%s", string(reqJSON)))
	}

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Debug log response
	if c.debug && c.debugLog != nil {
		var prettyJSON bytes.Buffer
		if json.Indent(&prettyJSON, respBody, "", "  ") == nil {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, prettyJSON.String()))
		} else {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, string(respBody)))
		}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Try to parse structured error
		var apiErr APIError
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.Error.Message != "" {
			errMsg := apiErr.Error.Message
			// Add helpful context for dimension mismatch errors
			if strings.Contains(errMsg, "must match the requested width and height") {
				errMsg += fmt.Sprintf("\n\nHint: Your reference image must be exactly %s pixels to match the requested video size.", req.Size)
				errMsg += "\nPlease resize your image or choose a different video size that matches your image dimensions."
			}
			return nil, &httpError{
				statusCode: resp.StatusCode,
				message:    errMsg,
				errorType:  apiErr.Error.Type,
			}
		}
		return nil, &httpError{
			statusCode: resp.StatusCode,
			message:    string(respBody),
		}
	}

	var result CreateVideoResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

type httpError struct {
	statusCode int
	message    string
	errorType  string
}

func (e *httpError) Error() string {
	if e.errorType != "" {
		return fmt.Sprintf("API error (%d - %s): %s", e.statusCode, e.errorType, e.message)
	}
	return fmt.Sprintf("API error (%d): %s", e.statusCode, e.message)
}

func isClientError(err error) bool {
	if httpErr, ok := err.(*httpError); ok {
		// 4xx errors are client errors - don't retry
		return httpErr.statusCode >= 400 && httpErr.statusCode < 500
	}
	return false
}

// ListVideos retrieves a list of video jobs
func (c *SoraClient) ListVideos(limit int) (*ListVideosResponse, error) {
	url := fmt.Sprintf("%s%s?limit=%d&order=desc", baseURL, createEndpoint, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Debug log request
	if c.debug && c.debugLog != nil {
		reqJSON, _ := json.MarshalIndent(map[string]interface{}{
			"method": "GET",
			"url":    url,
		}, "", "  ")
		c.debugLog(fmt.Sprintf("REQUEST:\n%s", string(reqJSON)))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Debug log response
	if c.debug && c.debugLog != nil {
		var prettyJSON bytes.Buffer
		if json.Indent(&prettyJSON, body, "", "  ") == nil {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, prettyJSON.String()))
		} else {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, string(body)))
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ListVideosResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetVideo retrieves the status and URL of a video generation job
func (c *SoraClient) GetVideo(videoID string) (*VideoResponse, error) {
	url := fmt.Sprintf("%s%s/%s", baseURL, createEndpoint, videoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Debug log request
	if c.debug && c.debugLog != nil {
		reqJSON, _ := json.MarshalIndent(map[string]interface{}{
			"method": "GET",
			"url":    url,
		}, "", "  ")
		c.debugLog(fmt.Sprintf("REQUEST:\n%s", string(reqJSON)))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Debug log response
	if c.debug && c.debugLog != nil {
		var prettyJSON bytes.Buffer
		if json.Indent(&prettyJSON, body, "", "  ") == nil {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, prettyJSON.String()))
		} else {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, string(body)))
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result VideoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// DownloadVideo downloads the video from the provided URL to the specified path
func (c *SoraClient) DownloadVideo(videoURL, outputPath string) error {
	resp, err := c.httpClient.Get(videoURL)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download video (status %d)", resp.StatusCode)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write video data: %w", err)
	}

	return nil
}

// DeleteVideo deletes a video job
func (c *SoraClient) DeleteVideo(videoID string) error {
	url := fmt.Sprintf("%s%s/%s", baseURL, createEndpoint, videoID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Debug log request
	if c.debug && c.debugLog != nil {
		reqJSON, _ := json.MarshalIndent(map[string]interface{}{
			"method": "DELETE",
			"url":    url,
		}, "", "  ")
		c.debugLog(fmt.Sprintf("REQUEST:\n%s", string(reqJSON)))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Debug log response
	if c.debug && c.debugLog != nil {
		var prettyJSON bytes.Buffer
		if json.Indent(&prettyJSON, body, "", "  ") == nil {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, prettyJSON.String()))
		} else {
			c.debugLog(fmt.Sprintf("RESPONSE [%d]:\n%s", resp.StatusCode, string(body)))
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// DownloadVideoContent downloads the video content directly from the /content endpoint
func (c *SoraClient) DownloadVideoContent(videoID, outputPath string) error {
	url := fmt.Sprintf("%s%s/%s/content", baseURL, createEndpoint, videoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Debug log request
	if c.debug && c.debugLog != nil {
		reqJSON, _ := json.MarshalIndent(map[string]interface{}{
			"method": "GET",
			"url":    url,
		}, "", "  ")
		c.debugLog(fmt.Sprintf("REQUEST:\n%s", string(reqJSON)))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download video content: %w", err)
	}
	defer resp.Body.Close()

	// Debug log response
	if c.debug && c.debugLog != nil {
		c.debugLog(fmt.Sprintf("RESPONSE [%d]: Streaming video content (Content-Type: %s, Content-Length: %s)",
			resp.StatusCode,
			resp.Header.Get("Content-Type"),
			resp.Header.Get("Content-Length")))
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to download video content (status %d): %s", resp.StatusCode, string(body))
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write video data: %w", err)
	}

	return nil
}
