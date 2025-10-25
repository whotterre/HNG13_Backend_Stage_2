package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"task_2/models"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// GenerateSummaryImage creates an image with country statistics
func GenerateSummaryImage(totalCountries int, topCountries []models.Country, lastRefreshed time.Time, outputPath string) error {
	// Create cache directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Image dimensions
	width := 800
	height := 400

	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background with white
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw border
	drawBorder(img, color.RGBA{R: 50, G: 50, B: 50, A: 255})

	// Prepare text content
	lines := []string{
		"Country Data Summary",
		"",
		fmt.Sprintf("Total Countries: %d", totalCountries),
		"",
		"Top 5 Countries by Estimated GDP:",
	}

	for i, country := range topCountries {
		if i >= 5 {
			break
		}
		gdpStr := "N/A"
		if country.EstimatedGDP != nil {
			gdpStr = fmt.Sprintf("$%.2f", *country.EstimatedGDP)
		}
		lines = append(lines, fmt.Sprintf("  %d. %s - %s", i+1, country.Name, gdpStr))
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Last Refreshed: %s", lastRefreshed.Format(time.RFC3339)))

	// Draw text
	y := 40
	for _, line := range lines {
		drawText(img, 30, y, line, color.Black)
		y += 25
	}

	// Save to file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

func drawText(img *image.RGBA, x, y int, label string, col color.Color) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func drawBorder(img *image.RGBA, col color.Color) {
	bounds := img.Bounds()
	// Top border
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		img.Set(x, bounds.Min.Y, col)
		img.Set(x, bounds.Min.Y+1, col)
	}
	// Bottom border
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		img.Set(x, bounds.Max.Y-1, col)
		img.Set(x, bounds.Max.Y-2, col)
	}
	// Left border
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		img.Set(bounds.Min.X, y, col)
		img.Set(bounds.Min.X+1, y, col)
	}
	// Right border
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		img.Set(bounds.Max.X-1, y, col)
		img.Set(bounds.Max.X-2, y, col)
	}
}
