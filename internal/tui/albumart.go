package tui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

// decodeAlbumArt decodes JPEG or PNG image data into an image.Image.
func decodeAlbumArt(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// resizeImage performs nearest-neighbor downscale to targetWidth, preserving
// aspect ratio. The returned image always has an even number of rows so that
// half-block rendering pairs evenly.
func resizeImage(img image.Image, targetWidth int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW == 0 || srcH == 0 {
		return img
	}

	// Each character cell is roughly twice as tall as it is wide, so the
	// target height in pixels should account for 2 pixels per row of text.
	targetHeight := (targetWidth * srcH) / srcW
	// Ensure even row count for half-block pairing.
	if targetHeight%2 != 0 {
		targetHeight++
	}
	if targetHeight < 2 {
		targetHeight = 2
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	for y := 0; y < targetHeight; y++ {
		srcY := y * srcH / targetHeight
		for x := 0; x < targetWidth; x++ {
			srcX := x * srcW / targetWidth
			dst.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}
	return dst
}

// renderHalfBlocks converts an image to a string using Unicode upper-half-block
// characters (▀) with truecolor ANSI foreground/background colors. Each text
// row represents two pixel rows.
func renderHalfBlocks(img image.Image) string {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	var sb strings.Builder
	for y := bounds.Min.Y; y < bounds.Min.Y+h; y += 2 {
		for x := bounds.Min.X; x < bounds.Min.X+w; x++ {
			// Upper pixel → foreground, lower pixel → background.
			ur, ug, ub, _ := img.At(x, y).RGBA()
			var lr, lg, lb uint32
			if y+1 < bounds.Min.Y+h {
				lr, lg, lb, _ = img.At(x, y+1).RGBA()
			}
			// RGBA returns 16-bit values; shift to 8-bit.
			fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀",
				ur>>8, ug>>8, ub>>8,
				lr>>8, lg>>8, lb>>8,
			)
		}
		sb.WriteString("\x1b[0m")
		if y+2 < bounds.Min.Y+h {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// renderAlbumArt renders album art data as a half-block character art string.
// Returns a placeholder if data is nil/empty or cannot be decoded.
func renderAlbumArt(data []byte, maxWidth int) string {
	if len(data) == 0 {
		return albumArtPlaceholder(maxWidth)
	}

	img, err := decodeAlbumArt(data)
	if err != nil {
		return albumArtPlaceholder(maxWidth)
	}

	scaled := resizeImage(img, maxWidth)
	return renderHalfBlocks(scaled)
}

func albumArtPlaceholder(width int) string {
	if width < 4 {
		width = 4
	}
	// Create a simple bordered box with "No artwork" centered.
	height := width / 2
	if height < 3 {
		height = 3
	}

	// Build a placeholder with a muted color.
	gray := color.RGBA{R: 60, G: 60, B: 60, A: 255}
	img := image.NewRGBA(image.Rect(0, 0, width, height*2))
	for y := 0; y < height*2; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, gray)
		}
	}
	art := renderHalfBlocks(img)

	// Overlay "No artwork" text centered (we just put it below the block).
	label := "No artwork"
	if len(label) > width {
		label = label[:width]
	}
	pad := (width - len(label)) / 2
	if pad < 0 {
		pad = 0
	}
	return art + "\n" + strings.Repeat(" ", pad) + label
}
