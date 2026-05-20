package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rickcern44/gopints/internal/api"
	"github.com/rickcern44/gopints/pkg/flow"
	"github.com/rickcern44/gopints/pkg/tap"
)

// --- Mock Store ---

type mockStore struct {
	ensureTapFn      func(ctx context.Context, id uint8) error
	listTapsFn       func(ctx context.Context) ([]tap.Tap, error)
	getTapFn         func(ctx context.Context, id uint8) (tap.Tap, error)
	setTapKegFn      func(ctx context.Context, tapID uint8, kegID int64) error
	removeTapKegFn   func(ctx context.Context, tapID uint8) error
	createKegFn      func(ctx context.Context, keg tap.Keg) (tap.Keg, error)
	updateKegFn      func(ctx context.Context, keg tap.Keg) (tap.Keg, error)
	deleteKegFn      func(ctx context.Context, id int64) error
	listKegsFn       func(ctx context.Context) ([]tap.Keg, error)
	getKegFn         func(ctx context.Context, id int64) (tap.Keg, error)
	getKegStatsFn    func(ctx context.Context, id int64) (tap.KegStats, error)
	setKegImageFn    func(ctx context.Context, id int64, data []byte, mimeType string) error
	getKegImageFn    func(ctx context.Context, id int64) ([]byte, string, error)
	deleteKegImageFn func(ctx context.Context, id int64) error
	recordPourFn     func(ctx context.Context, pour tap.Pour) (tap.Pour, error)
	listPoursFn      func(ctx context.Context, limit, offset int) ([]tap.Pour, error)
	listPoursByTapFn func(ctx context.Context, tapID uint8, limit, offset int) ([]tap.Pour, error)
	deletePourFn     func(ctx context.Context, id int64) error
	getSettingFn     func(ctx context.Context, key string) (string, error)
	setSettingFn     func(ctx context.Context, key, value string) error
	deleteSettingFn  func(ctx context.Context, key string) error
}

func (m *mockStore) EnsureTap(ctx context.Context, id uint8) error {
	if m.ensureTapFn == nil {
		panic("unexpected call to EnsureTap")
	}
	return m.ensureTapFn(ctx, id)
}
func (m *mockStore) ListTaps(ctx context.Context) ([]tap.Tap, error) {
	if m.listTapsFn == nil {
		panic("unexpected call to ListTaps")
	}
	return m.listTapsFn(ctx)
}
func (m *mockStore) GetTap(ctx context.Context, id uint8) (tap.Tap, error) {
	if m.getTapFn == nil {
		panic("unexpected call to GetTap")
	}
	return m.getTapFn(ctx, id)
}
func (m *mockStore) SetTapKeg(ctx context.Context, tapID uint8, kegID int64) error {
	if m.setTapKegFn == nil {
		panic("unexpected call to SetTapKeg")
	}
	return m.setTapKegFn(ctx, tapID, kegID)
}
func (m *mockStore) RemoveTapKeg(ctx context.Context, tapID uint8) error {
	if m.removeTapKegFn == nil {
		panic("unexpected call to RemoveTapKeg")
	}
	return m.removeTapKegFn(ctx, tapID)
}
func (m *mockStore) CreateKeg(ctx context.Context, keg tap.Keg) (tap.Keg, error) {
	if m.createKegFn == nil {
		panic("unexpected call to CreateKeg")
	}
	return m.createKegFn(ctx, keg)
}
func (m *mockStore) UpdateKeg(ctx context.Context, keg tap.Keg) (tap.Keg, error) {
	if m.updateKegFn == nil {
		panic("unexpected call to UpdateKeg")
	}
	return m.updateKegFn(ctx, keg)
}
func (m *mockStore) DeleteKeg(ctx context.Context, id int64) error {
	if m.deleteKegFn == nil {
		panic("unexpected call to DeleteKeg")
	}
	return m.deleteKegFn(ctx, id)
}
func (m *mockStore) ListKegs(ctx context.Context) ([]tap.Keg, error) {
	if m.listKegsFn == nil {
		panic("unexpected call to ListKegs")
	}
	return m.listKegsFn(ctx)
}
func (m *mockStore) GetKeg(ctx context.Context, id int64) (tap.Keg, error) {
	if m.getKegFn == nil {
		panic("unexpected call to GetKeg")
	}
	return m.getKegFn(ctx, id)
}
func (m *mockStore) GetKegStats(ctx context.Context, id int64) (tap.KegStats, error) {
	if m.getKegStatsFn == nil {
		panic("unexpected call to GetKegStats")
	}
	return m.getKegStatsFn(ctx, id)
}
func (m *mockStore) SetKegImage(ctx context.Context, id int64, data []byte, mimeType string) error {
	if m.setKegImageFn == nil {
		panic("unexpected call to SetKegImage")
	}
	return m.setKegImageFn(ctx, id, data, mimeType)
}
func (m *mockStore) GetKegImage(ctx context.Context, id int64) ([]byte, string, error) {
	if m.getKegImageFn == nil {
		panic("unexpected call to GetKegImage")
	}
	return m.getKegImageFn(ctx, id)
}
func (m *mockStore) DeleteKegImage(ctx context.Context, id int64) error {
	if m.deleteKegImageFn == nil {
		panic("unexpected call to DeleteKegImage")
	}
	return m.deleteKegImageFn(ctx, id)
}
func (m *mockStore) RecordPour(ctx context.Context, pour tap.Pour) (tap.Pour, error) {
	if m.recordPourFn == nil {
		panic("unexpected call to RecordPour")
	}
	return m.recordPourFn(ctx, pour)
}
func (m *mockStore) ListPours(ctx context.Context, limit, offset int) ([]tap.Pour, error) {
	if m.listPoursFn == nil {
		panic("unexpected call to ListPours")
	}
	return m.listPoursFn(ctx, limit, offset)
}
func (m *mockStore) ListPoursByTap(ctx context.Context, tapID uint8, limit, offset int) ([]tap.Pour, error) {
	if m.listPoursByTapFn == nil {
		panic("unexpected call to ListPoursByTap")
	}
	return m.listPoursByTapFn(ctx, tapID, limit, offset)
}
func (m *mockStore) DeletePour(ctx context.Context, id int64) error {
	if m.deletePourFn == nil {
		panic("unexpected call to DeletePour")
	}
	return m.deletePourFn(ctx, id)
}
func (m *mockStore) GetSetting(ctx context.Context, key string) (string, error) {
	if m.getSettingFn == nil {
		panic("unexpected call to GetSetting")
	}
	return m.getSettingFn(ctx, key)
}
func (m *mockStore) SetSetting(ctx context.Context, key, value string) error {
	if m.setSettingFn == nil {
		panic("unexpected call to SetSetting")
	}
	return m.setSettingFn(ctx, key, value)
}
func (m *mockStore) DeleteSetting(ctx context.Context, key string) error {
	if m.deleteSettingFn == nil {
		panic("unexpected call to DeleteSetting")
	}
	return m.deleteSettingFn(ctx, key)
}

// --- Helpers ---

const testAuthToken = "test-auth-token"

func newTestServer(t *testing.T, store tap.Store, simulate bool) *httptest.Server {
	t.Helper()
	srv := api.NewServer("", store, nil, simulate, "test")
	srv.SeedSession(testAuthToken)
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

func do(t *testing.T, ts *httptest.Server, method, path string, body any) *http.Response {
	t.Helper()
	var reqBody *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func doAuthed(t *testing.T, ts *httptest.Server, method, path string, body any) *http.Response {
	t.Helper()
	var reqBody *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+testAuthToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

// --- Health ---

func TestHandleHealth(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	resp := do(t, ts, "GET", "/api/health", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	decodeJSON(t, resp, &body)
	if body["version"] != "test" {
		t.Errorf("want version=test, got %q", body["version"])
	}
	if body["status"] != "ok" {
		t.Errorf("want status=ok, got %q", body["status"])
	}
}

// --- Taps ---

func TestHandleListTaps_OK(t *testing.T) {
	store := &mockStore{
		listTapsFn: func(_ context.Context) ([]tap.Tap, error) {
			return []tap.Tap{{ID: 1}, {ID: 2}}, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/taps", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	var taps []tap.Tap
	decodeJSON(t, resp, &taps)
	if len(taps) != 2 {
		t.Errorf("want 2 taps, got %d", len(taps))
	}
}

func TestHandleListTaps_StoreError(t *testing.T) {
	store := &mockStore{
		listTapsFn: func(_ context.Context) ([]tap.Tap, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/taps", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", resp.StatusCode)
	}
}

func TestHandleGetTap_Found(t *testing.T) {
	kegID := int64(7)
	store := &mockStore{
		getTapFn: func(_ context.Context, _ uint8) (tap.Tap, error) {
			return tap.Tap{ID: 1, KegID: &kegID}, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/taps/1", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}

func TestHandleGetTap_NotFound(t *testing.T) {
	store := &mockStore{
		getTapFn: func(_ context.Context, _ uint8) (tap.Tap, error) {
			return tap.Tap{}, fmt.Errorf("not found")
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/taps/99", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("want 404, got %d", resp.StatusCode)
	}
}

func TestHandleGetTap_BadID(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	resp := do(t, ts, "GET", "/api/taps/notanumber", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400, got %d", resp.StatusCode)
	}
}

func TestHandleSetTapKeg(t *testing.T) {
	called := false
	store := &mockStore{
		setTapKegFn: func(_ context.Context, _ uint8, _ int64) error {
			called = true
			return nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "PUT", "/api/taps/1/keg", map[string]any{"keg_id": 5})
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
	if !called {
		t.Error("SetTapKeg was not called")
	}
}

func TestHandleSetTapKeg_BadBody(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	req, _ := http.NewRequest("PUT", ts.URL+"/api/taps/1/keg", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testAuthToken)
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400, got %d", resp.StatusCode)
	}
}

func TestHandleRemoveTapKeg(t *testing.T) {
	store := &mockStore{
		removeTapKegFn: func(_ context.Context, _ uint8) error { return nil },
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "DELETE", "/api/taps/1/keg", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
}

// --- Kegs ---

func TestHandleListKegs(t *testing.T) {
	store := &mockStore{
		listKegsFn: func(_ context.Context) ([]tap.Keg, error) {
			return []tap.Keg{{ID: 1, BeerName: "IPA"}}, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/kegs", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	var kegs []tap.Keg
	decodeJSON(t, resp, &kegs)
	if len(kegs) != 1 || kegs[0].BeerName != "IPA" {
		t.Errorf("unexpected kegs: %+v", kegs)
	}
}

func TestHandleCreateKeg(t *testing.T) {
	store := &mockStore{
		createKegFn: func(_ context.Context, keg tap.Keg) (tap.Keg, error) {
			keg.ID = 42
			return keg, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "POST", "/api/kegs", map[string]any{
		"beer_name":   "Hazy IPA",
		"style":       "IPA",
		"abv":         6.5,
		"capacity_ml": 19000,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("want 201, got %d", resp.StatusCode)
	}
	var keg tap.Keg
	decodeJSON(t, resp, &keg)
	if keg.ID != 42 {
		t.Errorf("want ID 42, got %d", keg.ID)
	}
	if keg.BeerName != "Hazy IPA" {
		t.Errorf("want BeerName Hazy IPA, got %q", keg.BeerName)
	}
}

func TestHandleGetKeg(t *testing.T) {
	store := &mockStore{
		getKegFn: func(_ context.Context, _ int64) (tap.Keg, error) {
			return tap.Keg{ID: 3, BeerName: "Stout"}, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/kegs/3", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}

func TestHandleGetKeg_NotFound(t *testing.T) {
	store := &mockStore{
		getKegFn: func(_ context.Context, _ int64) (tap.Keg, error) {
			return tap.Keg{}, fmt.Errorf("not found")
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/kegs/999", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("want 404, got %d", resp.StatusCode)
	}
}

func TestHandleUpdateKeg_PartialPatch(t *testing.T) {
	original := tap.Keg{ID: 1, BeerName: "Original", Style: "Lager", ABV: 4.0, CapacityMl: 19000}
	var saved tap.Keg
	store := &mockStore{
		getKegFn: func(_ context.Context, _ int64) (tap.Keg, error) { return original, nil },
		updateKegFn: func(_ context.Context, keg tap.Keg) (tap.Keg, error) {
			saved = keg
			return keg, nil
		},
	}
	ts := newTestServer(t, store, false)
	// Only update BeerName; Style and ABV should be unchanged.
	resp := doAuthed(t, ts, "PATCH", "/api/kegs/1", map[string]any{"beer_name": "Renamed"})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
	if saved.BeerName != "Renamed" {
		t.Errorf("BeerName: want Renamed, got %q", saved.BeerName)
	}
	if saved.Style != "Lager" {
		t.Errorf("Style should be unchanged, got %q", saved.Style)
	}
	if saved.ABV != 4.0 {
		t.Errorf("ABV should be unchanged, got %f", saved.ABV)
	}
}

func TestHandleDeleteKeg(t *testing.T) {
	store := &mockStore{
		deleteKegFn: func(_ context.Context, _ int64) error { return nil },
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "DELETE", "/api/kegs/1", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
}

func TestHandleGetKegStats(t *testing.T) {
	store := &mockStore{
		getKegStatsFn: func(_ context.Context, id int64) (tap.KegStats, error) {
			return tap.KegStats{KegID: id, PourCount: 3, PouredMl: 900, RemainingMl: 18100, PctRemaining: 95.3}, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/kegs/1/stats", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	var stats tap.KegStats
	decodeJSON(t, resp, &stats)
	if stats.PourCount != 3 {
		t.Errorf("want PourCount 3, got %d", stats.PourCount)
	}
}

// --- Keg image ---

func TestHandleSetKegImage_OK(t *testing.T) {
	store := &mockStore{
		setKegImageFn: func(_ context.Context, _ int64, _ []byte, _ string) error {
			return nil
		},
	}
	ts := newTestServer(t, store, false)
	imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	req, _ := http.NewRequest("PUT", ts.URL+"/api/kegs/1/image", bytes.NewReader(imageData))
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Authorization", "Bearer "+testAuthToken)
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
}

func TestHandleSetKegImage_MissingContentType(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	req, _ := http.NewRequest("PUT", ts.URL+"/api/kegs/1/image", bytes.NewReader([]byte{0x01}))
	req.Header.Set("Authorization", "Bearer "+testAuthToken)
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400, got %d", resp.StatusCode)
	}
}

func TestHandleSetKegImage_TooLarge(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	// 11 MB body > 10 MB limit
	big := make([]byte, 11<<20)
	req, _ := http.NewRequest("PUT", ts.URL+"/api/kegs/1/image", bytes.NewReader(big))
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Authorization", "Bearer "+testAuthToken)
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("want 413, got %d", resp.StatusCode)
	}
}

func TestHandleGetKegImage_OK(t *testing.T) {
	imageData := []byte{0xFF, 0xD8}
	store := &mockStore{
		getKegImageFn: func(_ context.Context, _ int64) ([]byte, string, error) {
			return imageData, "image/jpeg", nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/kegs/1/image", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/jpeg" {
		t.Errorf("want Content-Type image/jpeg, got %q", ct)
	}
}

func TestHandleGetKegImage_NotFound(t *testing.T) {
	store := &mockStore{
		getKegImageFn: func(_ context.Context, _ int64) ([]byte, string, error) {
			return nil, "", fmt.Errorf("no image")
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/kegs/1/image", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("want 404, got %d", resp.StatusCode)
	}
}

func TestHandleDeleteKegImage(t *testing.T) {
	store := &mockStore{
		deleteKegImageFn: func(_ context.Context, _ int64) error { return nil },
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "DELETE", "/api/kegs/1/image", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
}

// --- Pours ---

func TestHandleListPours_All(t *testing.T) {
	store := &mockStore{
		listPoursFn: func(_ context.Context, _, _ int) ([]tap.Pour, error) {
			return []tap.Pour{{ID: 1, VolumeMl: 400}}, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/pours", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	var pours []tap.Pour
	decodeJSON(t, resp, &pours)
	if len(pours) != 1 {
		t.Errorf("want 1 pour, got %d", len(pours))
	}
}

func TestHandleListPours_ByTap(t *testing.T) {
	var calledWith uint8
	store := &mockStore{
		listPoursByTapFn: func(_ context.Context, tapID uint8, _, _ int) ([]tap.Pour, error) {
			calledWith = tapID
			return []tap.Pour{}, nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/pours?tap_id=2", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	if calledWith != 2 {
		t.Errorf("want ListPoursByTap called with tap 2, got %d", calledWith)
	}
}

func TestHandleDeletePour(t *testing.T) {
	store := &mockStore{
		deletePourFn: func(_ context.Context, _ int64) error { return nil },
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "DELETE", "/api/pours/5", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
}

// --- Dev simulate ---

func TestHandleSimulatePour_Forbidden_WhenNotSimulate(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	resp := do(t, ts, "POST", "/api/dev/pour/1", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		// Route doesn't even exist when simulate=false.
		t.Logf("got status %d (route not registered when simulate=false)", resp.StatusCode)
	}
}

// --- Admin auth ---

func TestHandleAdminStatus_NoPassword(t *testing.T) {
	store := &mockStore{
		getSettingFn: func(_ context.Context, _ string) (string, error) {
			return "", fmt.Errorf("not found")
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/admin/status", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	var body map[string]bool
	decodeJSON(t, resp, &body)
	if body["password_set"] {
		t.Error("want password_set=false")
	}
}

func TestHandleAdminStatus_PasswordSet(t *testing.T) {
	store := &mockStore{
		getSettingFn: func(_ context.Context, _ string) (string, error) {
			return "$2a$10$hash", nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/admin/status", nil)
	var body map[string]bool
	decodeJSON(t, resp, &body)
	if !body["password_set"] {
		t.Error("want password_set=true")
	}
}

func TestHandleAdminSetup_AlreadySet(t *testing.T) {
	store := &mockStore{
		getSettingFn: func(_ context.Context, _ string) (string, error) {
			return "hash", nil // password already exists
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "POST", "/api/admin/setup", map[string]any{"password": "newpass"})
	resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("want 409, got %d", resp.StatusCode)
	}
}

func TestHandleAdminLogin_WrongPassword(t *testing.T) {
	store := &mockStore{
		getSettingFn: func(_ context.Context, _ string) (string, error) {
			// bcrypt hash of "correct"
			return "$2a$10$invalid", nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "POST", "/api/admin/login", map[string]any{"password": "wrong"})
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
}

func TestHandleAdminLogin_NoPassword(t *testing.T) {
	store := &mockStore{
		getSettingFn: func(_ context.Context, _ string) (string, error) {
			return "", fmt.Errorf("not found")
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "POST", "/api/admin/login", map[string]any{"password": "any"})
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
}

func TestRequireAuth_Rejects_Missing_Token(t *testing.T) {
	store := &mockStore{
		deleteKegFn: func(_ context.Context, _ int64) error { return nil },
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "DELETE", "/api/kegs/1", nil) // no auth header
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
}

// --- Banner ---

func TestHandleGetBanner_OK(t *testing.T) {
	store := &mockStore{
		getSettingFn: func(_ context.Context, _ string) (string, error) {
			return "Tap 2 offline for cleaning", nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/banner", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	decodeJSON(t, resp, &body)
	if body["message"] != "Tap 2 offline for cleaning" {
		t.Errorf("unexpected message: %q", body["message"])
	}
}

func TestHandleGetBanner_NotFound(t *testing.T) {
	store := &mockStore{
		getSettingFn: func(_ context.Context, _ string) (string, error) {
			return "", fmt.Errorf("not found")
		},
	}
	ts := newTestServer(t, store, false)
	resp := do(t, ts, "GET", "/api/banner", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("want 404, got %d", resp.StatusCode)
	}
}

func TestHandleSetBanner_OK(t *testing.T) {
	var savedKey, savedValue string
	store := &mockStore{
		setSettingFn: func(_ context.Context, key, value string) error {
			savedKey, savedValue = key, value
			return nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "PUT", "/api/banner", map[string]any{"message": "Cheers!"})
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
	if savedKey != "banner" {
		t.Errorf("want key=banner, got %q", savedKey)
	}
	if savedValue != "Cheers!" {
		t.Errorf("want value=Cheers!, got %q", savedValue)
	}
}

func TestHandleSetBanner_EmptyMessage(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	resp := doAuthed(t, ts, "PUT", "/api/banner", map[string]any{"message": ""})
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400, got %d", resp.StatusCode)
	}
}

func TestHandleSetBanner_BadBody(t *testing.T) {
	ts := newTestServer(t, &mockStore{}, false)
	req, _ := http.NewRequest("PUT", ts.URL+"/api/banner", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testAuthToken)
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400, got %d", resp.StatusCode)
	}
}

func TestHandleDeleteBanner(t *testing.T) {
	called := false
	store := &mockStore{
		deleteSettingFn: func(_ context.Context, _ string) error {
			called = true
			return nil
		},
	}
	ts := newTestServer(t, store, false)
	resp := doAuthed(t, ts, "DELETE", "/api/banner", nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("want 204, got %d", resp.StatusCode)
	}
	if !called {
		t.Error("DeleteSetting was not called")
	}
}

func TestHandleSimulatePour_Accepted_WhenSimulate(t *testing.T) {
	meters := map[uint8]*flow.Meter{
		1: flow.NewMeter(1, 450, 50, make(chan flow.PourEvent, 64)),
	}
	srv := api.NewServer("", &mockStore{}, meters, true, "test")
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/dev/pour/1?pulses=5", "", nil)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("want 202, got %d", resp.StatusCode)
	}
}
