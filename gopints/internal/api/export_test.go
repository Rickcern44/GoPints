package api

// SetCheckHostIsPublicForTest overrides the SSRF host-allow check used by
// fetchRemoteImage, so external tests can exercise the from-url handlers
// against an httptest.Server, which always binds to loopback — a real
// target for this exact check in production. Only compiled under `go test`.
func SetCheckHostIsPublicForTest(fn func(host string) error) (restore func()) {
	prev := checkHostIsPublic
	checkHostIsPublic = fn
	return func() { checkHostIsPublic = prev }
}
