package modal

import "testing"

func TestImagePreviewContentSizeAppliesModalWidthCap(t *testing.T) {
	w, h := ImagePreviewContentSize(200, 80)
	if w != 124 {
		t.Fatalf("unexpected width: got %d, want 124", w)
	}
	if h != 72 {
		t.Fatalf("unexpected height: got %d, want 72", h)
	}
}

func TestImagePreviewContentSizeMinimums(t *testing.T) {
	w, h := ImagePreviewContentSize(1, 1)
	if w != 1 || h != 1 {
		t.Fatalf("expected minimum size 1x1, got %dx%d", w, h)
	}
}
