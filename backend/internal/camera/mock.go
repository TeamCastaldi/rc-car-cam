package camera

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"time"
)

// MockSource is a Source that repeatedly serves one generated JPEG test
// image, paced by interval. It needs no camera hardware and exists for
// development and testing before real capture code is written.
type MockSource struct {
	frame    []byte
	interval time.Duration
}

var _ Source = (*MockSource)(nil)

// NewMockSource builds a MockSource that yields a frame roughly every
// interval. Pass 0 to return frames immediately on every call.
func NewMockSource(interval time.Duration) (*MockSource, error) {
	frame, err := generateTestFrame()
	if err != nil {
		return nil, err
	}
	return &MockSource{frame: frame, interval: interval}, nil
}

// NextFrame implements Source. It always returns the same generated frame —
// the "loop" is repeated delivery of one image, paced by interval. The
// returned slice is a copy, safe for the caller to hold or mutate.
func (m *MockSource) NextFrame(ctx context.Context) ([]byte, error) {
	if m.interval <= 0 {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return m.frameCopy(), nil
	}

	timer := time.NewTimer(m.interval)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return m.frameCopy(), nil
	}
}

func (m *MockSource) frameCopy() []byte {
	frame := make([]byte, len(m.frame))
	copy(frame, m.frame)
	return frame
}

// generateTestFrame renders a small solid-color JPEG in memory, so the mock
// needs no binary image asset checked into the repo.
func generateTestFrame() ([]byte, error) {
	const size = 320
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 30, G: 144, B: 255, A: 255}}, image.Point{}, draw.Src)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
