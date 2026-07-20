package stream

import (
	"bytes"
	"context"
	"errors"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teamcastaldi/rc-car-cam/backend/internal/camera"
)

// splitParts extracts each part's body from a multipart/x-mixed-replace
// stream. It doesn't use mime/multipart.Reader because that requires a
// closing boundary ("--boundary--"), which Handler never writes: the
// stream is meant to run until the client disconnects, not to terminate
// itself. Real MJPEG consumers (a browser's <img> tag) parse the same way
// this does — boundary-delimited, no closing marker required.
func splitParts(t *testing.T, body []byte, boundary string) [][]byte {
	t.Helper()

	delim := []byte("--" + boundary + "\r\n")
	segments := bytes.Split(body, delim)
	if len(segments) == 0 || len(segments[0]) != 0 {
		t.Fatalf("expected body to start with boundary %q", delim)
	}
	segments = segments[1:] // drop the empty pre-boundary segment

	parts := make([][]byte, 0, len(segments))
	for i, seg := range segments {
		if i < len(segments)-1 {
			seg = bytes.TrimSuffix(seg, []byte("\r\n"))
		}
		headerEnd := bytes.Index(seg, []byte("\r\n\r\n"))
		if headerEnd == -1 {
			t.Fatalf("part %d missing header/body separator", i+1)
		}
		parts = append(parts, seg[headerEnd+4:])
	}
	return parts
}

// fakeSource is a camera.Source whose behavior is controlled per-test.
type fakeSource struct {
	next func(ctx context.Context) ([]byte, error)
}

func (f *fakeSource) NextFrame(ctx context.Context) ([]byte, error) {
	return f.next(ctx)
}

var _ camera.Source = (*fakeSource)(nil)

// countingSource wraps a real Source and reports context.Canceled after max
// successful calls, ending the handler's loop deterministically. It uses
// context.Canceled (rather than an arbitrary sentinel error) purely so the
// handler treats it as an expected disconnect and doesn't log it as a real
// Source error on every test run.
type countingSource struct {
	src   camera.Source
	max   int
	calls int
}

func (c *countingSource) NextFrame(ctx context.Context) ([]byte, error) {
	if c.calls >= c.max {
		return nil, context.Canceled
	}
	c.calls++
	return c.src.NextFrame(ctx)
}

var _ camera.Source = (*countingSource)(nil)

func TestHandler_ServeHTTP_WritesFramesFromMockSource(t *testing.T) {
	mock, err := camera.NewMockSource(0)
	if err != nil {
		t.Fatalf("NewMockSource: %v", err)
	}
	wantFrame, err := mock.NextFrame(context.Background())
	if err != nil {
		t.Fatalf("NextFrame (reference): %v", err)
	}

	src := &countingSource{src: mock, max: 3}
	h := NewHandler(src)

	req := httptest.NewRequest("GET", "/stream", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !rec.Flushed {
		t.Error("expected response to be flushed")
	}

	_, params, err := mime.ParseMediaType(rec.Header().Get("Content-Type"))
	if err != nil {
		t.Fatalf("parse Content-Type: %v", err)
	}
	parts := splitParts(t, rec.Body.Bytes(), params["boundary"])

	if len(parts) != 3 {
		t.Fatalf("got %d parts, want 3", len(parts))
	}
	for i, body := range parts {
		if !bytes.Equal(body, wantFrame) {
			t.Errorf("part %d body did not match the mock source's frame", i+1)
		}
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("Content-Type: image/jpeg")) {
		t.Error("expected parts to carry a Content-Type: image/jpeg header")
	}
}

func TestHandler_ServeHTTP_SetsMultipartContentType(t *testing.T) {
	src := &fakeSource{next: func(ctx context.Context) ([]byte, error) {
		return nil, context.Canceled
	}}
	h := NewHandler(src)

	req := httptest.NewRequest("GET", "/stream", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	_, params, err := mime.ParseMediaType(rec.Header().Get("Content-Type"))
	if err != nil {
		t.Fatalf("parse Content-Type: %v", err)
	}
	if params["boundary"] == "" {
		t.Error("expected a non-empty multipart boundary")
	}
}

func TestHandler_ServeHTTP_StopsWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	src := &fakeSource{next: func(ctx context.Context) ([]byte, error) {
		calls++
		if calls == 2 {
			cancel()
			return nil, ctx.Err()
		}
		return []byte("frame"), nil
	}}
	h := NewHandler(src)

	req := httptest.NewRequest("GET", "/stream", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if calls != 2 {
		t.Errorf("got %d calls to NextFrame, want 2", calls)
	}

	_, params, err := mime.ParseMediaType(rec.Header().Get("Content-Type"))
	if err != nil {
		t.Fatalf("parse Content-Type: %v", err)
	}
	parts := splitParts(t, rec.Body.Bytes(), params["boundary"])
	if len(parts) != 1 {
		t.Errorf("got %d parts, want 1", len(parts))
	}
}

func TestHandler_ServeHTTP_StopsOnSourceError(t *testing.T) {
	calls := 0
	src := &fakeSource{next: func(ctx context.Context) ([]byte, error) {
		calls++
		if calls == 1 {
			return []byte("frame"), nil
		}
		return nil, errors.New("boom")
	}}
	h := NewHandler(src)

	req := httptest.NewRequest("GET", "/stream", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	_, params, err := mime.ParseMediaType(rec.Header().Get("Content-Type"))
	if err != nil {
		t.Fatalf("parse Content-Type: %v", err)
	}
	parts := splitParts(t, rec.Body.Bytes(), params["boundary"])
	if len(parts) != 1 {
		t.Errorf("got %d parts, want 1", len(parts))
	}
}

func TestHandler_ServeHTTP_NoFramesOnImmediateError(t *testing.T) {
	mock, err := camera.NewMockSource(0)
	if err != nil {
		t.Fatalf("NewMockSource: %v", err)
	}
	src := &countingSource{src: mock, max: 0}
	h := NewHandler(src)

	req := httptest.NewRequest("GET", "/stream", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct == "" {
		t.Error("expected a Content-Type header even with zero frames")
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body, got %d bytes", rec.Body.Len())
	}
}

func TestHandler_ServeHTTP_NilSourceReturns500(t *testing.T) {
	h := &Handler{}

	req := httptest.NewRequest("GET", "/stream", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 500 {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

// nonFlushingRecorder implements http.ResponseWriter by delegating to an
// httptest.ResponseRecorder, but — critically — does not embed it, so it
// doesn't inherit ResponseRecorder's Flush method and therefore doesn't
// satisfy http.Flusher. This simulates a ResponseWriter that can't support
// streaming.
type nonFlushingRecorder struct {
	rec *httptest.ResponseRecorder
}

func (n *nonFlushingRecorder) Header() http.Header         { return n.rec.Header() }
func (n *nonFlushingRecorder) Write(b []byte) (int, error) { return n.rec.Write(b) }
func (n *nonFlushingRecorder) WriteHeader(code int)        { n.rec.WriteHeader(code) }

var _ http.ResponseWriter = (*nonFlushingRecorder)(nil)

func TestHandler_ServeHTTP_NonFlusherReturns500(t *testing.T) {
	mock, err := camera.NewMockSource(0)
	if err != nil {
		t.Fatalf("NewMockSource: %v", err)
	}
	h := NewHandler(mock)

	req := httptest.NewRequest("GET", "/stream", nil)
	rec := &nonFlushingRecorder{rec: httptest.NewRecorder()}
	h.ServeHTTP(rec, req)

	if rec.rec.Code != 500 {
		t.Errorf("status = %d, want 500", rec.rec.Code)
	}
}
