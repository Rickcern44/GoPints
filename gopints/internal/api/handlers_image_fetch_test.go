package api

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func smallPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	img.Set(0, 0, color.RGBA{R: 200, G: 130, B: 26, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

// allowAllHosts bypasses the loopback rejection so httptest.Server (which
// always binds to 127.0.0.1) can stand in for a real remote host.
func allowAllHosts(t *testing.T) {
	t.Helper()
	prev := checkHostIsPublic
	checkHostIsPublic = func(string) error { return nil }
	t.Cleanup(func() { checkHostIsPublic = prev })
}

func TestFetchRemoteImage_RejectsPrivateAndLoopbackHosts(t *testing.T) {
	cases := []string{
		"http://127.0.0.1:1/image.png",
		"http://localhost:1/image.png",
		"http://10.0.0.5:1/image.png",
		"http://172.16.0.5:1/image.png",
		"http://192.168.1.5:1/image.png",
		"http://169.254.169.254/latest/meta-data/",
		"http://[::1]:1/image.png",
	}
	for _, u := range cases {
		t.Run(u, func(t *testing.T) {
			_, _, err := fetchRemoteImage(context.Background(), u)
			if err == nil {
				t.Fatalf("expected rejection for %q, got nil error", u)
			}
			var fe *fetchError
			if !errors.As(err, &fe) || fe.status != http.StatusBadRequest {
				t.Fatalf("expected 400 fetchError, got %v", err)
			}
		})
	}
}

func TestFetchRemoteImage_RejectsNonHTTPScheme(t *testing.T) {
	_, _, err := fetchRemoteImage(context.Background(), "file:///etc/passwd")
	if err == nil {
		t.Fatal("expected rejection for file:// scheme")
	}
}

func TestFetchRemoteImage_RejectsRedirect(t *testing.T) {
	allowAllHosts(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/other", http.StatusFound)
	}))
	defer srv.Close()

	_, _, err := fetchRemoteImage(context.Background(), srv.URL+"/image.png")
	if err == nil {
		t.Fatal("expected rejection for redirecting URL")
	}
	var fe *fetchError
	if !errors.As(err, &fe) || fe.status != http.StatusBadRequest {
		t.Fatalf("expected 400 fetchError, got %v", err)
	}
}

func TestFetchRemoteImage_RejectsOversizedBody(t *testing.T) {
	allowAllHosts(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		oversized := make([]byte, maxImageBytes+1)
		_, _ = w.Write(oversized)
	}))
	defer srv.Close()

	_, _, err := fetchRemoteImage(context.Background(), srv.URL+"/image.png")
	if err == nil {
		t.Fatal("expected rejection for oversized body")
	}
	var fe *fetchError
	if !errors.As(err, &fe) || fe.status != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413 fetchError, got %v", err)
	}
}

func TestFetchRemoteImage_RejectsNonImageContent(t *testing.T) {
	allowAllHosts(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png") // lying header, real bytes are text
		_, _ = w.Write([]byte("<html>not an image</html>"))
	}))
	defer srv.Close()

	_, _, err := fetchRemoteImage(context.Background(), srv.URL+"/image.png")
	if err == nil {
		t.Fatal("expected rejection for non-image content")
	}
	var fe *fetchError
	if !errors.As(err, &fe) || fe.status != http.StatusBadRequest {
		t.Fatalf("expected 400 fetchError, got %v", err)
	}
}

func TestFetchRemoteImage_SucceedsOnValidImage(t *testing.T) {
	allowAllHosts(t)
	pngBytes := smallPNG(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream") // deliberately wrong; sniffing must win
		_, _ = w.Write(pngBytes)
	}))
	defer srv.Close()

	data, mimeType, err := fetchRemoteImage(context.Background(), srv.URL+"/image.png")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if mimeType != "image/png" {
		t.Fatalf("expected sniffed mime type image/png, got %q", mimeType)
	}
	if !bytes.Equal(data, pngBytes) {
		t.Fatal("returned data does not match served image bytes")
	}
}
