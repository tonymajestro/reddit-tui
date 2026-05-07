package images

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeRenderer struct {
	output string
	err    error
	opts   RenderOptions
}

func (f *fakeRenderer) Render(img image.Image, opts RenderOptions) (string, error) {
	f.opts = opts
	if f.err != nil {
		return "", f.err
	}
	return f.output, nil
}

func TestLooksLikeImagePath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/img/test.png", true},
		{"/img/test.jpeg", true},
		{"/img/test.jpg", true},
		{"/img/test.gif", true},
		{"/img/test.txt", false},
	}

	for _, tt := range tests {
		if got := looksLikeImagePath(tt.path); got != tt.want {
			t.Fatalf("looksLikeImagePath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestRenderFromURLRejectsNonImageContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html></html>"))
	}))
	defer server.Close()

	client := NewPreviewClient(5, "auto")
	_, err := client.RenderFromURL(server.URL+"/post", 30, 10)
	if !errors.Is(err, ErrUnsupportedImageURL) {
		t.Fatalf("expected unsupported image url error, got %v", err)
	}
}

func TestRenderFromURLUsesRendererOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(makePNGBytes(t))
	}))
	defer server.Close()

	renderer := &fakeRenderer{output: "rendered-preview"}
	client := NewPreviewClient(5, "kitty")
	client.renderer = renderer

	got, err := client.RenderFromURL(server.URL+"/image", 50, 25)
	if err != nil {
		t.Fatalf("unexpected error rendering image preview: %v", err)
	}

	if got != "rendered-preview" {
		t.Fatalf("expected rendered output, got %q", got)
	}

	if renderer.opts.Width != 50 || renderer.opts.Height != 25 {
		t.Fatalf("unexpected render size opts: %+v", renderer.opts)
	}

	if renderer.opts.Protocol != ProtocolKitty {
		t.Fatalf("expected kitty protocol option, got %q", renderer.opts.Protocol)
	}
}

func TestRenderFromURLReturnsDecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid image bytes"))
	}))
	defer server.Close()

	client := NewPreviewClient(5, "auto")
	_, err := client.RenderFromURL(server.URL+"/image", 30, 10)
	if !errors.Is(err, ErrCannotDecodeImage) {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func makePNGBytes(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("could not encode test png: %v", err)
	}
	return buf.Bytes()
}
