// Package camera defines how video frames are produced for the stream —
// the seam between a frame producer (mock, or real camera hardware) and a
// consumer such as the internal/stream MJPEG handler.
package camera

import "context"

// Source produces a continuous sequence of JPEG-encoded video frames.
type Source interface {
	// NextFrame returns the next frame as JPEG-encoded bytes. It blocks
	// until a frame is available, an error occurs, or ctx is canceled.
	NextFrame(ctx context.Context) ([]byte, error)
}
