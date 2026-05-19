package tap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

const schemaVersion = 3

// SQLiteStore implements Store using a local SQLite database.
// Use ":memory:" as path for in-process testing.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (or creates) a SQLite database at path and applies
// any pending schema migrations.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("tap: open db %s: %w", path, err)
	}
	db.SetMaxOpenConns(1) // SQLite writer serialisation
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

// Close releases the database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// --- Tap ---

func (s *SQLiteStore) EnsureTap(ctx context.Context, id uint8) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO taps (id, keg_id) VALUES (?, NULL)`, id)
	if err != nil {
		return fmt.Errorf("tap: ensure tap %d: %w", id, err)
	}
	return nil
}

func (s *SQLiteStore) ListTaps(ctx context.Context) ([]Tap, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, keg_id FROM taps ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("tap: list taps: %w", err)
	}
	defer rows.Close()

	var taps []Tap
	for rows.Next() {
		var t Tap
		if err := rows.Scan(&t.ID, &t.KegID); err != nil {
			return nil, err
		}
		if t.KegID != nil {
			keg, err := s.GetKeg(ctx, *t.KegID)
			if err == nil {
				t.Keg = &keg
			}
		}
		taps = append(taps, t)
	}
	return taps, rows.Err()
}

func (s *SQLiteStore) GetTap(ctx context.Context, id uint8) (Tap, error) {
	var t Tap
	err := s.db.QueryRowContext(ctx, `SELECT id, keg_id FROM taps WHERE id = ?`, id).
		Scan(&t.ID, &t.KegID)
	if errors.Is(err, sql.ErrNoRows) {
		return Tap{}, fmt.Errorf("tap: tap %d not found", id)
	}
	if err != nil {
		return Tap{}, fmt.Errorf("tap: get tap %d: %w", id, err)
	}
	if t.KegID != nil {
		keg, err := s.GetKeg(ctx, *t.KegID)
		if err == nil {
			t.Keg = &keg
		}
	}
	return t, nil
}

func (s *SQLiteStore) SetTapKeg(ctx context.Context, tapID uint8, kegID int64) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO taps (id, keg_id) VALUES (?, ?)
		 ON CONFLICT(id) DO UPDATE SET keg_id = excluded.keg_id`,
		tapID, kegID)
	if err != nil {
		return fmt.Errorf("tap: set tap %d keg: %w", tapID, err)
	}
	return nil
}

func (s *SQLiteStore) RemoveTapKeg(ctx context.Context, tapID uint8) error {
	_, err := s.db.ExecContext(ctx, `UPDATE taps SET keg_id = NULL WHERE id = ?`, tapID)
	if err != nil {
		return fmt.Errorf("tap: remove tap %d keg: %w", tapID, err)
	}
	return nil
}

// --- Keg ---

func (s *SQLiteStore) CreateKeg(ctx context.Context, keg Keg) (Keg, error) {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO kegs (beer_name, style, abv, brewery, notes, capacity_ml, added_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		keg.BeerName, keg.Style, keg.ABV, keg.Brewery, keg.Notes, keg.CapacityMl, now.UnixMilli())
	if err != nil {
		return Keg{}, fmt.Errorf("tap: create keg: %w", err)
	}
	id, _ := res.LastInsertId()
	keg.ID = id
	keg.AddedAt = now
	return keg, nil
}

func (s *SQLiteStore) UpdateKeg(ctx context.Context, keg Keg) (Keg, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE kegs SET beer_name=?, style=?, abv=?, brewery=?, notes=?, capacity_ml=?
		 WHERE id=?`,
		keg.BeerName, keg.Style, keg.ABV, keg.Brewery, keg.Notes, keg.CapacityMl, keg.ID)
	if err != nil {
		return Keg{}, fmt.Errorf("tap: update keg %d: %w", keg.ID, err)
	}
	return s.GetKeg(ctx, keg.ID)
}

func (s *SQLiteStore) DeleteKeg(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM kegs WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("tap: delete keg %d: %w", id, err)
	}
	return nil
}

func (s *SQLiteStore) ListKegs(ctx context.Context) ([]Keg, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, beer_name, style, abv, brewery, notes, capacity_ml, added_at, image_mime_type
		 FROM kegs ORDER BY added_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("tap: list kegs: %w", err)
	}
	defer rows.Close()

	var kegs []Keg
	for rows.Next() {
		k, err := scanKeg(rows)
		if err != nil {
			return nil, err
		}
		kegs = append(kegs, k)
	}
	return kegs, rows.Err()
}

func (s *SQLiteStore) GetKeg(ctx context.Context, id int64) (Keg, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, beer_name, style, abv, brewery, notes, capacity_ml, added_at, image_mime_type
		 FROM kegs WHERE id = ?`, id)
	k, err := scanKeg(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Keg{}, fmt.Errorf("tap: keg %d not found", id)
	}
	if err != nil {
		return Keg{}, fmt.Errorf("tap: get keg %d: %w", id, err)
	}
	return k, nil
}

func (s *SQLiteStore) GetKegStats(ctx context.Context, id int64) (KegStats, error) {
	keg, err := s.GetKeg(ctx, id)
	if err != nil {
		return KegStats{}, err
	}

	var count int
	var totalMl sql.NullFloat64
	err = s.db.QueryRowContext(ctx,
		`SELECT COUNT(*), SUM(volume_ml) FROM pours
		 WHERE tap_id IN (SELECT id FROM taps WHERE keg_id = ?)`, id).
		Scan(&count, &totalMl)
	if err != nil {
		return KegStats{}, fmt.Errorf("tap: keg stats %d: %w", id, err)
	}

	poured := totalMl.Float64
	remaining := keg.CapacityMl - poured
	if remaining < 0 {
		remaining = 0
	}
	pct := 0.0
	if keg.CapacityMl > 0 {
		pct = (remaining / keg.CapacityMl) * 100
	}
	return KegStats{
		KegID:        id,
		PourCount:    count,
		PouredMl:     poured,
		RemainingMl:  remaining,
		PctRemaining: pct,
	}, nil
}

// --- Keg image ---

func (s *SQLiteStore) SetKegImage(ctx context.Context, id int64, data []byte, mimeType string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE kegs SET image = ?, image_mime_type = ? WHERE id = ?`,
		data, mimeType, id)
	if err != nil {
		return fmt.Errorf("tap: set keg %d image: %w", id, err)
	}
	return nil
}

func (s *SQLiteStore) GetKegImage(ctx context.Context, id int64) ([]byte, string, error) {
	var data []byte
	var mimeType string
	err := s.db.QueryRowContext(ctx,
		`SELECT image, image_mime_type FROM kegs WHERE id = ?`, id).
		Scan(&data, &mimeType)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", fmt.Errorf("tap: keg %d not found", id)
	}
	if err != nil {
		return nil, "", fmt.Errorf("tap: get keg %d image: %w", id, err)
	}
	if len(data) == 0 {
		return nil, "", fmt.Errorf("tap: keg %d has no image", id)
	}
	return data, mimeType, nil
}

func (s *SQLiteStore) DeleteKegImage(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE kegs SET image = NULL, image_mime_type = '' WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("tap: delete keg %d image: %w", id, err)
	}
	return nil
}

// --- Pour ---

func (s *SQLiteStore) RecordPour(ctx context.Context, pour Pour) (Pour, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO pours (tap_id, volume_ml, started_at, ended_at) VALUES (?, ?, ?, ?)`,
		pour.TapID, pour.VolumeMl, pour.StartedAt.UnixMilli(), pour.EndedAt.UnixMilli())
	if err != nil {
		return Pour{}, fmt.Errorf("tap: record pour: %w", err)
	}
	id, _ := res.LastInsertId()
	pour.ID = id
	return pour, nil
}

func (s *SQLiteStore) ListPours(ctx context.Context, limit, offset int) ([]Pour, error) {
	return s.queryPours(ctx,
		`SELECT id, tap_id, volume_ml, started_at, ended_at FROM pours ORDER BY started_at DESC LIMIT ? OFFSET ?`,
		limit, offset)
}

func (s *SQLiteStore) ListPoursByTap(ctx context.Context, tapID uint8, limit, offset int) ([]Pour, error) {
	return s.queryPours(ctx,
		`SELECT id, tap_id, volume_ml, started_at, ended_at FROM pours WHERE tap_id = ? ORDER BY started_at DESC LIMIT ? OFFSET ?`,
		tapID, limit, offset)
}

func (s *SQLiteStore) DeletePour(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM pours WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("tap: delete pour %d: %w", id, err)
	}
	return nil
}

func (s *SQLiteStore) queryPours(ctx context.Context, query string, args ...any) ([]Pour, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("tap: query pours: %w", err)
	}
	defer rows.Close()

	var pours []Pour
	for rows.Next() {
		var p Pour
		var startedAt, endedAt int64
		if err := rows.Scan(&p.ID, &p.TapID, &p.VolumeMl, &startedAt, &endedAt); err != nil {
			return nil, err
		}
		p.StartedAt = time.UnixMilli(startedAt).UTC()
		p.EndedAt = time.UnixMilli(endedAt).UTC()
		pours = append(pours, p)
	}
	return pours, rows.Err()
}

// --- Scan helpers ---

type scanner interface {
	Scan(dest ...any) error
}

func scanKeg(s scanner) (Keg, error) {
	var k Keg
	var addedAt int64
	err := s.Scan(&k.ID, &k.BeerName, &k.Style, &k.ABV, &k.Brewery, &k.Notes,
		&k.CapacityMl, &addedAt, &k.ImageMimeType)
	if err != nil {
		return Keg{}, err
	}
	k.AddedAt = time.UnixMilli(addedAt).UTC()
	return k, nil
}

// --- Schema migrations ---

func migrate(db *sql.DB) error {
	var v int
	if err := db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		return fmt.Errorf("tap: read schema version: %w", err)
	}

	if v < 1 {
		if err := applyMigration1(db); err != nil {
			return fmt.Errorf("tap: migration 1: %w", err)
		}
	}
	if v < 2 {
		if err := applyMigration2(db); err != nil {
			return fmt.Errorf("tap: migration 2: %w", err)
		}
	}
	if v < 3 {
		if err := applyMigration3(db); err != nil {
			return fmt.Errorf("tap: migration 3: %w", err)
		}
	}

	_, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %d", schemaVersion))
	return err
}

func applyMigration1(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA foreign_keys=ON;

		CREATE TABLE IF NOT EXISTS kegs (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			beer_name   TEXT    NOT NULL DEFAULT '',
			capacity_ml REAL    NOT NULL DEFAULT 0,
			added_at    INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS taps (
			id     INTEGER PRIMARY KEY,
			keg_id INTEGER REFERENCES kegs(id)
		);

		CREATE TABLE IF NOT EXISTS pours (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			tap_id     INTEGER NOT NULL,
			volume_ml  REAL    NOT NULL,
			started_at INTEGER NOT NULL,
			ended_at   INTEGER NOT NULL
		);

		CREATE INDEX IF NOT EXISTS pours_tap_id     ON pours(tap_id);
		CREATE INDEX IF NOT EXISTS pours_started_at ON pours(started_at DESC);
	`)
	return err
}

// applyMigration3 converts second-precision timestamps to milliseconds.
// The heuristic: any timestamp < 2e12 is in seconds (year ~2033 in ms).
func applyMigration3(db *sql.DB) error {
	stmts := []string{
		`UPDATE kegs  SET added_at   = added_at   * 1000 WHERE added_at   < 2000000000000`,
		`UPDATE pours SET started_at = started_at * 1000 WHERE started_at < 2000000000000`,
		`UPDATE pours SET ended_at   = ended_at   * 1000 WHERE ended_at   < 2000000000000`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func applyMigration2(db *sql.DB) error {
	stmts := []string{
		`ALTER TABLE kegs ADD COLUMN style            TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE kegs ADD COLUMN abv              REAL NOT NULL DEFAULT 0`,
		`ALTER TABLE kegs ADD COLUMN brewery          TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE kegs ADD COLUMN notes            TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE kegs ADD COLUMN image            BLOB`,
		`ALTER TABLE kegs ADD COLUMN image_mime_type  TEXT NOT NULL DEFAULT ''`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
