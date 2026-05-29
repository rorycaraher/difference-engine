package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"
)

// Request holds the parameters needed to reproduce a mixdown.
type Request struct {
	ID        int64
	CreatedAt string
	Track     string
	Stems     []string
	Volumes   []float64
}

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	// SQLite supports one writer at a time; cap the pool to avoid locking errors.
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable WAL: %w", err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func migrate(db *sql.DB) error {
	var version int
	_ = db.QueryRow("PRAGMA user_version").Scan(&version)

	if version < 1 {
		if _, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS requests (
				id         INTEGER PRIMARY KEY AUTOINCREMENT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				track      TEXT NOT NULL DEFAULT '',
				stems      TEXT NOT NULL,
				volumes    TEXT NOT NULL
			)
		`); err != nil {
			return err
		}
		// No-op on a fresh DB (column exists); adds the column when upgrading
		// a pre-versioned DB that was created without the track column.
		_, _ = db.Exec(`ALTER TABLE requests ADD COLUMN track TEXT NOT NULL DEFAULT ''`)

		if _, err := db.Exec(`PRAGMA user_version = 1`); err != nil {
			return err
		}
	}

	return nil
}

// RecordRequest persists the track, stems, and (already-rounded) volumes and returns the new row ID.
func (s *Store) RecordRequest(track string, stems []string, volumes []float64) (int64, error) {
	stemsJSON, err := json.Marshal(stems)
	if err != nil {
		return 0, err
	}
	volsJSON, err := json.Marshal(volumes)
	if err != nil {
		return 0, err
	}
	result, err := s.db.Exec(
		`INSERT INTO requests (track, stems, volumes) VALUES (?, ?, ?)`,
		track, string(stemsJSON), string(volsJSON),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetRequest fetches a stored request by ID. Returns nil, nil if not found.
func (s *Store) GetRequest(id int64) (*Request, error) {
	row := s.db.QueryRow(
		`SELECT id, created_at, track, stems, volumes FROM requests WHERE id = ?`, id,
	)
	var req Request
	var stemsJSON, volsJSON string
	if err := row.Scan(&req.ID, &req.CreatedAt, &req.Track, &stemsJSON, &volsJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := json.Unmarshal([]byte(stemsJSON), &req.Stems); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(volsJSON), &req.Volumes); err != nil {
		return nil, err
	}
	return &req, nil
}
