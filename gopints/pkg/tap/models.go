package tap

import "time"

// Keg represents a keg of beer installed in the kegerator.
type Keg struct {
	ID            int64
	BeerName      string
	Style         string  // e.g. "IPA", "Stout", "Lager"
	ABV           float64 // alcohol by volume, e.g. 6.5
	Brewery       string
	Notes         string // tasting notes or description
	CapacityMl    float64
	AddedAt       time.Time
	ImageMimeType string // empty if no image stored
}

// KegStats summarises pour activity for a keg.
type KegStats struct {
	KegID        int64   `json:"keg_id"`
	PourCount    int     `json:"pour_count"`
	PouredMl     float64 `json:"poured_ml"`
	RemainingMl  float64 `json:"remaining_ml"`
	PctRemaining float64 `json:"pct_remaining"` // 0–100
}

// Tap is a physical tap line. It may or may not have a keg assigned.
type Tap struct {
	ID    uint8
	KegID *int64
	Keg   *Keg // populated by store reads; nil if no keg assigned
}

// Pour is a recorded pour event with final volume.
type Pour struct {
	ID        int64
	TapID     uint8
	VolumeMl  float64
	StartedAt time.Time
	EndedAt   time.Time
}
