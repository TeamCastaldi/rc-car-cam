// Package stream serves camera frames over HTTP as an MJPEG
// (multipart/x-mixed-replace) response.
package stream

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/teamcastaldi/rc-car-cam/backend/internal/camera"
)

// Handler serves a camera.Source as a continuous MJPEG stream.
type Handler struct {
	Source camera.Source
}

var _ http.Handler = (*Handler)(nil)

// NewHandler builds a Handler that streams frames pulled from src.
func NewHandler(src camera.Source) *Handler {
	return &Handler{Source: src}
}

// ServeHTTP writes frames from h.Source as multipart/x-mixed-replace parts
// until the request's context is canceled (the client disconnects) or the
// Source returns an error.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Source == nil {
		http.Error(w, "stream: no camera source configured", http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "stream: response writer does not support flushing", http.StatusInternalServerError)
		return
	}

	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary="+mw.Boundary())
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)

	ctx := r.Context()

	for {
		frame, err := h.Source.NextFrame(ctx)
		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				log.Printf("stream: source error: %v", err)
			}
			return
		}

		header := textproto.MIMEHeader{}
		header.Set("Content-Type", "image/jpeg")
		header.Set("Content-Length", strconv.Itoa(len(frame)))

		part, err := mw.CreatePart(header)
		if err != nil {
			return
		}
		if _, err := part.Write(frame); err != nil {
			return
		}

		flusher.Flush()
	}
}
