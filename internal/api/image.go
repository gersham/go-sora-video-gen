package api

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

// parseSize parses a size string like "1280x720" into width and height
func parseSize(size string) (int, int, error) {
	parts := strings.Split(size, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("size must be in format WIDTHxHEIGHT")
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %w", err)
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %w", err)
	}

	return width, height, nil
}

// resizeAndCropToFill resizes and crops an image to fill the target dimensions
// using a "cover" strategy (scales to cover the entire target, cropping excess)
func resizeAndCropToFill(src image.Image, targetWidth, targetHeight int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Calculate scale factors to cover target dimensions
	scaleX := float64(targetWidth) / float64(srcWidth)
	scaleY := float64(targetHeight) / float64(srcHeight)

	// Use the larger scale to ensure we cover the entire target
	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY
	}

	// Calculate scaled dimensions
	scaledWidth := int(float64(srcWidth) * scale)
	scaledHeight := int(float64(srcHeight) * scale)

	// Resize using nearest neighbor (fast, simple)
	scaled := resizeImage(src, scaledWidth, scaledHeight)

	// Calculate crop offsets to center the image
	cropX := (scaledWidth - targetWidth) / 2
	cropY := (scaledHeight - targetHeight) / 2

	// Crop to target dimensions
	cropped := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			cropped.Set(x, y, scaled.At(x+cropX, y+cropY))
		}
	}

	return cropped
}

// resizeImage performs simple nearest-neighbor image scaling
func resizeImage(src image.Image, width, height int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	xRatio := float64(srcWidth) / float64(width)
	yRatio := float64(srcHeight) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	return dst
}
