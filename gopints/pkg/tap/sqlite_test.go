package tap_test

import (
	"context"
	"testing"
	"time"

	"github.com/rickcern44/gopints/pkg/tap"
)

func newTestStore(t *testing.T) *tap.SQLiteStore {
	t.Helper()
	s, err := tap.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func ctx() context.Context { return context.Background() }

// --- Migration ---

func TestSQLiteStore_Migration(t *testing.T) {
	// Opening a fresh :memory: store should succeed without error,
	// meaning both migrations ran cleanly.
	s := newTestStore(t)
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

// --- Tap ---

func TestSQLiteStore_SetTapKeg_OneKegPerTap(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "IPA", CapacityMl: 19000})
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}
	if err := s.EnsureTap(ctx(), 2); err != nil {
		t.Fatal(err)
	}

	// Assign keg to tap 1.
	if err := s.SetTapKeg(ctx(), 1, keg.ID); err != nil {
		t.Fatalf("SetTapKeg tap 1: %v", err)
	}

	// Assign same keg to tap 2 — should silently clear it from tap 1.
	if err := s.SetTapKeg(ctx(), 2, keg.ID); err != nil {
		t.Fatalf("SetTapKeg tap 2: %v", err)
	}

	tap1, _ := s.GetTap(ctx(), 1)
	if tap1.KegID != nil {
		t.Errorf("tap 1 should have no keg after reassignment, got keg_id=%d", *tap1.KegID)
	}
	tap2, _ := s.GetTap(ctx(), 2)
	if tap2.KegID == nil || *tap2.KegID != keg.ID {
		t.Errorf("tap 2 should have keg %d, got %v", keg.ID, tap2.KegID)
	}
}

func TestSQLiteStore_EnsureTap_Idempotent(t *testing.T) {
	s := newTestStore(t)
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatalf("first EnsureTap: %v", err)
	}
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatalf("second EnsureTap (idempotent): %v", err)
	}
}

func TestSQLiteStore_GetTap_NotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetTap(ctx(), 99); err == nil {
		t.Fatal("expected error for missing tap, got nil")
	}
}

func TestSQLiteStore_SetTapKeg_And_Get(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "Test IPA", CapacityMl: 19000})
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}

	if err := s.SetTapKeg(ctx(), 1, keg.ID); err != nil {
		t.Fatalf("SetTapKeg: %v", err)
	}
	got, err := s.GetTap(ctx(), 1)
	if err != nil {
		t.Fatalf("GetTap: %v", err)
	}
	if got.KegID == nil || *got.KegID != keg.ID {
		t.Errorf("want KegID %d, got %v", keg.ID, got.KegID)
	}
	if got.Keg == nil || got.Keg.BeerName != "Test IPA" {
		t.Error("Keg field not populated on GetTap")
	}
}

func TestSQLiteStore_RemoveTapKeg(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "Stout", CapacityMl: 19000})
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}
	if err := s.SetTapKeg(ctx(), 1, keg.ID); err != nil {
		t.Fatal(err)
	}

	if err := s.RemoveTapKeg(ctx(), 1); err != nil {
		t.Fatalf("RemoveTapKeg: %v", err)
	}
	got, _ := s.GetTap(ctx(), 1)
	if got.KegID != nil {
		t.Errorf("want nil KegID after remove, got %v", got.KegID)
	}
}

func TestSQLiteStore_ListTaps(t *testing.T) {
	s := newTestStore(t)
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}
	if err := s.EnsureTap(ctx(), 2); err != nil {
		t.Fatal(err)
	}

	taps, err := s.ListTaps(ctx())
	if err != nil {
		t.Fatalf("ListTaps: %v", err)
	}
	if len(taps) != 2 {
		t.Errorf("want 2 taps, got %d", len(taps))
	}
}

// --- Keg ---

func TestSQLiteStore_CreateKeg_RoundTrip(t *testing.T) {
	s := newTestStore(t)
	in := tap.Keg{
		BeerName:   "Hazy IPA",
		Style:      "IPA",
		ABV:        6.8,
		Brewery:    "Test Brewery",
		Notes:      "Tropical notes",
		CapacityMl: 19000,
	}
	got, err := s.CreateKeg(ctx(), in)
	if err != nil {
		t.Fatalf("CreateKeg: %v", err)
	}
	if got.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if got.BeerName != in.BeerName {
		t.Errorf("BeerName: want %q, got %q", in.BeerName, got.BeerName)
	}
	if got.ABV != in.ABV {
		t.Errorf("ABV: want %f, got %f", in.ABV, got.ABV)
	}
	if got.AddedAt.IsZero() {
		t.Error("AddedAt should be set")
	}
}

func TestSQLiteStore_UpdateKeg(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "Original", CapacityMl: 19000})

	keg.BeerName = "Updated"
	keg.Style = "Stout"
	keg.ABV = 7.5
	updated, err := s.UpdateKeg(ctx(), keg)
	if err != nil {
		t.Fatalf("UpdateKeg: %v", err)
	}
	if updated.BeerName != "Updated" {
		t.Errorf("want BeerName Updated, got %q", updated.BeerName)
	}
	if updated.Style != "Stout" {
		t.Errorf("want Style Stout, got %q", updated.Style)
	}
}

func TestSQLiteStore_DeleteKeg(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "Gone", CapacityMl: 19000})

	if err := s.DeleteKeg(ctx(), keg.ID); err != nil {
		t.Fatalf("DeleteKeg: %v", err)
	}
	if _, err := s.GetKeg(ctx(), keg.ID); err == nil {
		t.Fatal("expected error after DeleteKeg, got nil")
	}
}

func TestSQLiteStore_GetKeg_NotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetKeg(ctx(), 9999); err == nil {
		t.Fatal("expected error for missing keg, got nil")
	}
}

func TestSQLiteStore_ListKegs_OrderedByDate(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.CreateKeg(ctx(), tap.Keg{BeerName: "First", CapacityMl: 19000}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond) // ensure different timestamps
	if _, err := s.CreateKeg(ctx(), tap.Keg{BeerName: "Second", CapacityMl: 19000}); err != nil {
		t.Fatal(err)
	}

	kegs, err := s.ListKegs(ctx())
	if err != nil {
		t.Fatalf("ListKegs: %v", err)
	}
	if len(kegs) != 2 {
		t.Fatalf("want 2 kegs, got %d", len(kegs))
	}
	// Should be newest first.
	if kegs[0].BeerName != "Second" {
		t.Errorf("want Second first, got %q", kegs[0].BeerName)
	}
}

// --- KegStats ---

func TestSQLiteStore_GetKegStats_WithPours(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "Lager", CapacityMl: 19000})
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}
	if err := s.SetTapKeg(ctx(), 1, keg.ID); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	if _, err := s.RecordPour(ctx(), tap.Pour{TapID: 1, VolumeMl: 500, StartedAt: now, EndedAt: now.Add(10 * time.Second)}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.RecordPour(ctx(), tap.Pour{TapID: 1, VolumeMl: 300, StartedAt: now.Add(time.Minute), EndedAt: now.Add(time.Minute + 5*time.Second)}); err != nil {
		t.Fatal(err)
	}

	stats, err := s.GetKegStats(ctx(), keg.ID)
	if err != nil {
		t.Fatalf("GetKegStats: %v", err)
	}
	if stats.PourCount != 2 {
		t.Errorf("want PourCount 2, got %d", stats.PourCount)
	}
	if stats.PouredMl != 800 {
		t.Errorf("want PouredMl 800, got %.1f", stats.PouredMl)
	}
	if stats.RemainingMl != 18200 {
		t.Errorf("want RemainingMl 18200, got %.1f", stats.RemainingMl)
	}
}

func TestSQLiteStore_GetKegStats_NoPours(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "Fresh", CapacityMl: 19000})

	stats, err := s.GetKegStats(ctx(), keg.ID)
	if err != nil {
		t.Fatalf("GetKegStats: %v", err)
	}
	if stats.PctRemaining != 100 {
		t.Errorf("want 100%% remaining for new keg, got %.1f", stats.PctRemaining)
	}
	if stats.RemainingMl != 19000 {
		t.Errorf("want RemainingMl 19000, got %.1f", stats.RemainingMl)
	}
}

// --- Image ---

func TestSQLiteStore_Image_SetGetDelete(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "Visual", CapacityMl: 19000})

	imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes
	mimeType := "image/jpeg"

	if err := s.SetKegImage(ctx(), keg.ID, imageData, mimeType); err != nil {
		t.Fatalf("SetKegImage: %v", err)
	}

	gotData, gotMime, err := s.GetKegImage(ctx(), keg.ID)
	if err != nil {
		t.Fatalf("GetKegImage: %v", err)
	}
	if string(gotData) != string(imageData) {
		t.Error("image data mismatch")
	}
	if gotMime != mimeType {
		t.Errorf("want mime %q, got %q", mimeType, gotMime)
	}
	// ImageMimeType should be reflected in keg record.
	updated, _ := s.GetKeg(ctx(), keg.ID)
	if updated.ImageMimeType != mimeType {
		t.Errorf("keg ImageMimeType: want %q, got %q", mimeType, updated.ImageMimeType)
	}

	// Delete.
	if err := s.DeleteKegImage(ctx(), keg.ID); err != nil {
		t.Fatalf("DeleteKegImage: %v", err)
	}
	if _, _, err := s.GetKegImage(ctx(), keg.ID); err == nil {
		t.Fatal("expected error after DeleteKegImage, got nil")
	}
}

func TestSQLiteStore_GetKegImage_NoImage(t *testing.T) {
	s := newTestStore(t)
	keg, _ := s.CreateKeg(ctx(), tap.Keg{BeerName: "NoImg", CapacityMl: 19000})

	if _, _, err := s.GetKegImage(ctx(), keg.ID); err == nil {
		t.Fatal("expected error for keg with no image, got nil")
	}
}

// --- Pour ---

func TestSQLiteStore_RecordAndListPours(t *testing.T) {
	s := newTestStore(t)
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}
	if err := s.EnsureTap(ctx(), 2); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	if _, err := s.RecordPour(ctx(), tap.Pour{TapID: 1, VolumeMl: 400, StartedAt: now, EndedAt: now.Add(5 * time.Second)}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.RecordPour(ctx(), tap.Pour{TapID: 2, VolumeMl: 600, StartedAt: now.Add(time.Minute), EndedAt: now.Add(time.Minute + 8*time.Second)}); err != nil {
		t.Fatal(err)
	}

	all, err := s.ListPours(ctx(), 10, 0)
	if err != nil {
		t.Fatalf("ListPours: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("want 2 pours, got %d", len(all))
	}

	byTap, err := s.ListPoursByTap(ctx(), 1, 10, 0)
	if err != nil {
		t.Fatalf("ListPoursByTap: %v", err)
	}
	if len(byTap) != 1 {
		t.Errorf("want 1 pour for tap 1, got %d", len(byTap))
	}
	if byTap[0].VolumeMl != 400 {
		t.Errorf("want 400ml, got %.1f", byTap[0].VolumeMl)
	}
}

func TestSQLiteStore_ListPours_Pagination(t *testing.T) {
	s := newTestStore(t)
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		if _, err := s.RecordPour(ctx(), tap.Pour{TapID: 1, VolumeMl: float64(i * 100), StartedAt: now, EndedAt: now.Add(5 * time.Second)}); err != nil {
			t.Fatal(err)
		}
	}

	page1, _ := s.ListPours(ctx(), 2, 0)
	page2, _ := s.ListPours(ctx(), 2, 2)

	if len(page1) != 2 {
		t.Errorf("page1: want 2, got %d", len(page1))
	}
	if len(page2) != 2 {
		t.Errorf("page2: want 2, got %d", len(page2))
	}
	// No overlap.
	if page1[0].ID == page2[0].ID {
		t.Error("pages overlap")
	}
}

func TestSQLiteStore_DeletePour(t *testing.T) {
	s := newTestStore(t)
	if err := s.EnsureTap(ctx(), 1); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	pour, _ := s.RecordPour(ctx(), tap.Pour{TapID: 1, VolumeMl: 300, StartedAt: now, EndedAt: now.Add(5 * time.Second)})

	if err := s.DeletePour(ctx(), pour.ID); err != nil {
		t.Fatalf("DeletePour: %v", err)
	}
	pours, _ := s.ListPours(ctx(), 10, 0)
	if len(pours) != 0 {
		t.Errorf("want 0 pours after delete, got %d", len(pours))
	}
}

// --- Settings ---

func TestSQLiteStore_Setting_SetAndGet(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetSetting(ctx(), "banner", "Cheers!"); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	got, err := s.GetSetting(ctx(), "banner")
	if err != nil {
		t.Fatalf("GetSetting: %v", err)
	}
	if got != "Cheers!" {
		t.Errorf("want Cheers!, got %q", got)
	}
}

func TestSQLiteStore_Setting_Upsert(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetSetting(ctx(), "banner", "First"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetSetting(ctx(), "banner", "Second"); err != nil {
		t.Fatal(err)
	}
	got, _ := s.GetSetting(ctx(), "banner")
	if got != "Second" {
		t.Errorf("want Second after upsert, got %q", got)
	}
}

func TestSQLiteStore_Setting_GetNotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetSetting(ctx(), "missing"); err == nil {
		t.Error("expected error for missing setting, got nil")
	}
}

func TestSQLiteStore_Setting_Delete(t *testing.T) {
	s := newTestStore(t)
	if err := s.SetSetting(ctx(), "banner", "To be removed"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteSetting(ctx(), "banner"); err != nil {
		t.Fatalf("DeleteSetting: %v", err)
	}
	if _, err := s.GetSetting(ctx(), "banner"); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestSQLiteStore_Setting_DeleteNonExistent(t *testing.T) {
	s := newTestStore(t)
	// Deleting a key that doesn't exist should not error.
	if err := s.DeleteSetting(ctx(), "ghost"); err != nil {
		t.Errorf("unexpected error deleting non-existent key: %v", err)
	}
}
