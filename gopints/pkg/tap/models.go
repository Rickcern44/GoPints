package tap

import "time"

// Keg represents a keg of beer installed in the kegerator.
type Keg struct {
	ID            int64     `json:"id"`
	BeerName      string    `json:"beer_name"`
	Style         string    `json:"style"`
	ABV           float64   `json:"abv"`
	Brewery       string    `json:"brewery"`
	Notes         string    `json:"notes"`
	CapacityMl    float64   `json:"capacity_ml"`
	AddedAt       time.Time `json:"added_at"`
	ImageMimeType        string    `json:"image_mime_type"`
	ImageStyle           string    `json:"image_style"`
	BreweryImageMimeType string    `json:"brewery_image_mime_type"`
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
	ID    uint8  `json:"id"`
	KegID *int64 `json:"keg_id"`
	Keg   *Keg   `json:"keg"`
}

// Pour is a recorded pour event with final volume.
type Pour struct {
	ID        int64     `json:"id"`
	TapID     uint8     `json:"tap_id"`
	VolumeMl  float64   `json:"volume_ml"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}
