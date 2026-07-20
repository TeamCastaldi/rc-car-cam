package camera

import (
	"bytes"
	"context"
	"image/jpeg"
	"testing"
	"time"
)

func TestMockSource_NextFrame_ReturnsValidJPEG(t *testing.T) {
	src, err := NewMockSource(0)
	if err != nil {
		t.Fatalf("NewMockSource: %v", err)
	}

	frame, err := src.NextFrame(context.Background())
	if err != nil {
		t.Fatalf("NextFrame: %v", err)
	}

	if _, err := jpeg.Decode(bytes.NewReader(frame)); err != nil {
		t.Errorf("NextFrame returned invalid JPEG: %v", err)
	}
}

func TestMockSource_NextFrame_LoopsSameFrame(t *testing.T) {
	src, err := NewMockSource(0)
	if err != nil {
		t.Fatalf("NewMockSource: %v", err)
	}

	first, err := src.NextFrame(context.Background())
	if err != nil {
		t.Fatalf("NextFrame (1st call): %v", err)
	}
	second, err := src.NextFrame(context.Background())
	if err != nil {
		t.Fatalf("NextFrame (2nd call): %v", err)
	}

	if !bytes.Equal(first, second) {
		t.Error("expected NextFrame to return the same looped frame on repeated calls")
	}
}

func TestMockSource_NextFrame_RespectsCanceledContext(t *testing.T) {
	src, err := NewMockSource(0)
	if err != nil {
		t.Fatalf("NewMockSource: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := src.NextFrame(ctx); err == nil {
		t.Error("expected NextFrame to return an error for an already-canceled context")
	}
}

func TestMockSource_NextFrame_CancelDuringWaitReturnsPromptly(t *testing.T) {
	src, err := NewMockSource(time.Hour) // would hang without ctx support
	if err != nil {
		t.Fatalf("NewMockSource: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err = src.NextFrame(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("expected NextFrame to return an error when context times out")
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("NextFrame took %v to honor context cancellation, expected well under the 1h interval", elapsed)
	}
}
