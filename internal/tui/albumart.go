package tui

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"strings"
	"sync"
)

// imageProtocol represents the terminal's image rendering capability.
type imageProtocol int

const (
	protoHalfBlock imageProtocol = iota
	protoITerm2
	protoKitty
)

var (
	detectedProtocol     imageProtocol
	detectedProtocolOnce sync.Once
)

// detectImageProtocol checks environment variables to determine the best
// image rendering protocol the terminal supports.
func detectImageProtocol() imageProtocol {
	detectedProtocolOnce.Do(func() {
		termProgram := os.Getenv("TERM_PROGRAM")
		term := os.Getenv("TERM")

		// Kitty
		if term == "xterm-kitty" {
			detectedProtocol = protoKitty
			return
		}

		// WezTerm supports the Kitty graphics protocol.
		if termProgram == "WezTerm" {
			detectedProtocol = protoKitty
			return
		}

		// Ghostty supports the Kitty graphics protocol.
		if termProgram == "ghostty" {
			detectedProtocol = protoKitty
			return
		}

		// iTerm2
		if termProgram == "iTerm.app" || os.Getenv("LC_TERMINAL") == "iTerm2" {
			detectedProtocol = protoITerm2
			return
		}

		detectedProtocol = protoHalfBlock
	})
	return detectedProtocol
}

// decodeAlbumArt decodes JPEG or PNG image data into an image.Image.
func decodeAlbumArt(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// renderAlbumArt renders album art using the best available terminal protocol.
// Returns a placeholder if data is nil/empty or cannot be decoded.
func renderAlbumArt(data []byte, maxWidth int) string {
	if len(data) == 0 {
		return albumArtPlaceholder(maxWidth)
	}

	proto := detectImageProtocol()

	switch proto {
	case protoITerm2:
		return renderITerm2(data, maxWidth)
	case protoKitty:
		return renderKitty(data, maxWidth)
	default:
		return renderHalfBlockArt(data, maxWidth)
	}
}

// --- iTerm2 inline image protocol ---

func renderITerm2(data []byte, widthCols int) string {
	img, err := decodeAlbumArt(data)
	if err != nil {
		return renderHalfBlockArt(data, widthCols)
	}
	bounds := img.Bounds()
	rows := imageRowCount(bounds.Dx(), bounds.Dy(), widthCols)

	b64 := base64.StdEncoding.EncodeToString(data)

	// The escape sequence renders the image inline. It is emitted standalone
	// (not inside lipgloss containers) to avoid layout measurement issues.
	return fmt.Sprintf("\x1b]1337;File=inline=1;width=%d;height=%d;preserveAspectRatio=1:%s\a",
		widthCols, rows, b64)
}

// --- Kitty graphics protocol ---

func renderKitty(data []byte, widthCols int) string {
	img, err := decodeAlbumArt(data)
	if err != nil {
		return renderHalfBlockArt(data, widthCols)
	}
	bounds := img.Bounds()
	rows := imageRowCount(bounds.Dx(), bounds.Dy(), widthCols)

	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return renderHalfBlockArt(data, widthCols)
	}

	b64 := base64.StdEncoding.EncodeToString(pngBuf.Bytes())

	// Kitty protocol transmits base64 data in chunks of up to 4096 bytes.
	// Emitted standalone, not inside lipgloss containers.
	const chunkSize = 4096

	var sb strings.Builder
	for i := 0; i < len(b64); i += chunkSize {
		end := i + chunkSize
		if end > len(b64) {
			end = len(b64)
		}
		chunk := b64[i:end]
		isLast := end >= len(b64)

		if i == 0 {
			more := 1
			if isLast {
				more = 0
			}
			fmt.Fprintf(&sb, "\x1b_Ga=T,f=100,c=%d,r=%d,m=%d;%s\x1b\\",
				widthCols, rows, more, chunk)
		} else {
			more := 1
			if isLast {
				more = 0
			}
			fmt.Fprintf(&sb, "\x1b_Gm=%d;%s\x1b\\", more, chunk)
		}
	}

	return sb.String()
}

// imageRowCount computes how many terminal rows an image will occupy when
// displayed at widthCols columns, preserving aspect ratio. Terminal cells
// are roughly twice as tall as they are wide, so each row ≈ 2 character widths
// of vertical space.
func imageRowCount(imgW, imgH, widthCols int) int {
	if imgW == 0 {
		return widthCols / 2
	}
	rows := (widthCols * imgH) / (imgW * 2)
	if rows < 1 {
		rows = 1
	}
	return rows
}

// --- Half-block fallback ---

func renderHalfBlockArt(data []byte, maxWidth int) string {
	img, err := decodeAlbumArt(data)
	if err != nil {
		return albumArtPlaceholder(maxWidth)
	}

	scaled := resizeImage(img, maxWidth)
	return renderHalfBlocks(scaled)
}

// resizeImage downscales img to targetWidth using area-average (box filter)
// sampling for smooth results. Preserves aspect ratio and ensures an even row
// count for half-block rendering.
func resizeImage(img image.Image, targetWidth int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW == 0 || srcH == 0 {
		return img
	}

	targetHeight := (targetWidth * srcH) / srcW
	if targetHeight%2 != 0 {
		targetHeight++
	}
	if targetHeight < 2 {
		targetHeight = 2
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	for y := 0; y < targetHeight; y++ {
		sy0 := y * srcH / targetHeight
		sy1 := (y + 1) * srcH / targetHeight
		if sy1 == sy0 {
			sy1 = sy0 + 1
		}
		for x := 0; x < targetWidth; x++ {
			sx0 := x * srcW / targetWidth
			sx1 := (x + 1) * srcW / targetWidth
			if sx1 == sx0 {
				sx1 = sx0 + 1
			}

			var rSum, gSum, bSum uint64
			count := uint64((sx1 - sx0) * (sy1 - sy0))
			for sy := sy0; sy < sy1; sy++ {
				for sx := sx0; sx < sx1; sx++ {
					r, g, b, _ := img.At(bounds.Min.X+sx, bounds.Min.Y+sy).RGBA()
					rSum += uint64(r)
					gSum += uint64(g)
					bSum += uint64(b)
				}
			}
			dst.Set(x, y, color.RGBA{
				R: uint8(rSum / count >> 8),
				G: uint8(gSum / count >> 8),
				B: uint8(bSum / count >> 8),
				A: 255,
			})
		}
	}
	return dst
}

// renderHalfBlocks converts an image to a string using Unicode upper-half-block
// characters (▀) with truecolor ANSI foreground/background colors.
func renderHalfBlocks(img image.Image) string {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	var sb strings.Builder
	for y := bounds.Min.Y; y < bounds.Min.Y+h; y += 2 {
		for x := bounds.Min.X; x < bounds.Min.X+w; x++ {
			ur, ug, ub, _ := img.At(x, y).RGBA()
			var lr, lg, lb uint32
			if y+1 < bounds.Min.Y+h {
				lr, lg, lb, _ = img.At(x, y+1).RGBA()
			}
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

func albumArtPlaceholder(width int) string {
	if width < 4 {
		width = 4
	}
	height := width / 2
	if height < 3 {
		height = 3
	}

	gray := color.RGBA{R: 60, G: 60, B: 60, A: 255}
	img := image.NewRGBA(image.Rect(0, 0, width, height*2))
	for y := 0; y < height*2; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, gray)
		}
	}
	art := renderHalfBlocks(img)

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
