package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const remoteImageFetchTimeout = 8 * time.Second

var errRedirectNotAllowed = errors.New("redirects are not allowed")

// fetchError carries an HTTP status code alongside a message so handlers can
// surface an admin-facing reason without leaking internal error details.
type fetchError struct {
	status int
	msg    string
}

func (e *fetchError) Error() string { return e.msg }

func newFetchError(status int, msg string) error {
	return &fetchError{status: status, msg: msg}
}

// writeFetchError maps a fetchRemoteImage error to an HTTP response, falling
// back to 502 for anything that wasn't explicitly classified.
func writeFetchError(w http.ResponseWriter, err error) {
	var fe *fetchError
	if errors.As(err, &fe) {
		writeError(w, fe.status, fe.msg)
		return
	}
	writeError(w, http.StatusBadGateway, "failed to fetch image")
}

// fetchRemoteImage downloads an admin-supplied image URL server-side. It
// rejects non-http(s) schemes, loopback/private/link-local/unspecified
// addresses, redirects, oversized bodies, and content that doesn't sniff as
// an image — see design.md for the SSRF threat model this defends against.
func fetchRemoteImage(ctx context.Context, rawURL string) (data []byte, mimeType string, err error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return nil, "", newFetchError(http.StatusBadRequest, "invalid URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, "", newFetchError(http.StatusBadRequest, "URL must use http or https")
	}
	if err := checkHostIsPublic(parsed.Hostname()); err != nil {
		return nil, "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", newFetchError(http.StatusBadRequest, "invalid URL")
	}

	client := &http.Client{
		Timeout: remoteImageFetchTimeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return errRedirectNotAllowed
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, errRedirectNotAllowed) {
			return nil, "", newFetchError(http.StatusBadRequest, "URL redirected; please provide a direct image link")
		}
		return nil, "", newFetchError(http.StatusBadGateway, "could not reach URL")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", newFetchError(http.StatusBadGateway, fmt.Sprintf("remote server returned status %d", resp.StatusCode))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxImageBytes+1))
	if err != nil {
		return nil, "", newFetchError(http.StatusBadGateway, "failed to read remote image")
	}
	if len(body) > maxImageBytes {
		return nil, "", newFetchError(http.StatusRequestEntityTooLarge, "image exceeds 10 MB limit")
	}
	if len(body) == 0 {
		return nil, "", newFetchError(http.StatusBadRequest, "empty image body")
	}

	sniffed := http.DetectContentType(body)
	if !strings.HasPrefix(sniffed, "image/") {
		return nil, "", newFetchError(http.StatusBadRequest, "URL does not point to a valid image")
	}

	return body, sniffed, nil
}

// checkHostIsPublic is a package-level var so tests can stub it out when
// exercising fetchRemoteImage against an httptest.Server, which always binds
// to loopback — a real target for this exact check in production.
var checkHostIsPublic = defaultCheckHostIsPublic

// defaultCheckHostIsPublic rejects hosts that resolve to loopback, private,
// link-local, or unspecified addresses (including the common cloud metadata
// address 169.254.169.254, covered by IsLinkLocalUnicast).
func defaultCheckHostIsPublic(host string) error {
	if ip := net.ParseIP(host); ip != nil {
		if isDisallowedHostIP(ip) {
			return newFetchError(http.StatusBadRequest, "URL resolves to a disallowed address")
		}
		return nil
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return newFetchError(http.StatusBadRequest, "could not resolve host")
	}
	for _, ip := range ips {
		if isDisallowedHostIP(ip) {
			return newFetchError(http.StatusBadRequest, "URL resolves to a disallowed address")
		}
	}
	return nil
}

func isDisallowedHostIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified()
}
