package tap

import "context"

// Store is the persistence interface for taps, kegs, and pours.
// Implementations must be safe for concurrent use.
type Store interface {
	// Tap operations
	EnsureTap(ctx context.Context, id uint8) error
	CreateTap(ctx context.Context, id uint8) error
	DeleteTap(ctx context.Context, id uint8) error
	ListTaps(ctx context.Context) ([]Tap, error)
	GetTap(ctx context.Context, id uint8) (Tap, error)
	SetTapKeg(ctx context.Context, tapID uint8, kegID int64) error
	RemoveTapKeg(ctx context.Context, tapID uint8) error

	// Keg operations
	CreateKeg(ctx context.Context, keg Keg) (Keg, error)
	UpdateKeg(ctx context.Context, keg Keg) (Keg, error)
	DeleteKeg(ctx context.Context, id int64) error
	ListKegs(ctx context.Context) ([]Keg, error)
	GetKeg(ctx context.Context, id int64) (Keg, error)
	GetKegStats(ctx context.Context, id int64) (KegStats, error)

	// Keg image operations (stored as a blob)
	SetKegImage(ctx context.Context, id int64, data []byte, mimeType string) error
	GetKegImage(ctx context.Context, id int64) (data []byte, mimeType string, err error)
	DeleteKegImage(ctx context.Context, id int64) error

	// Brewery image operations
	SetBreweryImage(ctx context.Context, id int64, data []byte, mimeType string) error
	GetBreweryImage(ctx context.Context, id int64) (data []byte, mimeType string, err error)
	DeleteBreweryImage(ctx context.Context, id int64) error

	// Pour operations
	RecordPour(ctx context.Context, pour Pour) (Pour, error)
	ListPours(ctx context.Context, limit, offset int) ([]Pour, error)
	ListPoursByTap(ctx context.Context, tapID uint8, limit, offset int) ([]Pour, error)
	DeletePour(ctx context.Context, id int64) error

	// Settings operations (key-value store for server-side configuration)
	GetSetting(ctx context.Context, key string) (string, error)
	SetSetting(ctx context.Context, key, value string) error
	DeleteSetting(ctx context.Context, key string) error
}
