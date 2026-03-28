package db

import (
	"database/sql"

	"github.com/realglebivanov/hstd/xrayconnectord/internal/link"
	_ "modernc.org/sqlite"
)

type LinkInfo struct {
	Index   int
	Comment string
	Enabled bool
	IPs     string
}

func (d *DB) IsEnabled(l *link.Link) (bool, error) {
	row := d.db.QueryRow(
		`SELECT links.enabled FROM links WHERE links.idx = ?`,
		l.Index,
	)

	var isEnabled bool
	err := row.Scan(&isEnabled)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	return isEnabled, nil
}

func (d *DB) List(count int) ([]LinkInfo, error) {
	rows, err := d.db.Query(`
		WITH RECURSIVE g(value) AS (
			SELECT 0
			UNION ALL
			SELECT value + 1 FROM g WHERE value < ?-1
		)
		SELECT g.value, COALESCE(l.comment, ''), COALESCE(l.enabled, 1), COALESCE(GROUP_CONCAT(i.ip), '')
		FROM g
		LEFT JOIN links l ON l.idx = g.value
		LEFT JOIN ips i ON i.link_idx = g.value
		GROUP BY g.value
		ORDER BY COUNT(i.ip) DESC, g.value ASC
	`, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []LinkInfo
	for rows.Next() {
		var li LinkInfo
		if err := rows.Scan(&li.Index, &li.Comment, &li.Enabled, &li.IPs); err != nil {
			return nil, err
		}
		result = append(result, li)
	}
	return result, rows.Err()
}

func (d *DB) TrackIP(l *link.Link, ip string) error {
	_, err := d.db.Exec(
		`INSERT OR IGNORE INTO links (idx) VALUES (?)`,
		l.Index,
	)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(
		`INSERT OR IGNORE INTO ips (link_idx, ip) VALUES (?, ?)`,
		l.Index, ip,
	)
	return err
}

func (d *DB) SetComment(index int, comment string) error {
	_, err := d.db.Exec(
		`INSERT INTO links (idx, comment) VALUES (?, ?)
		 ON CONFLICT(idx) DO UPDATE SET comment = excluded.comment`,
		index, comment,
	)
	return err
}

func (d *DB) SetEnabled(index int, enabled bool) error {
	_, err := d.db.Exec(
		`INSERT INTO links (idx, enabled) VALUES (?, ?)
		 ON CONFLICT(idx) DO UPDATE SET enabled = excluded.enabled`,
		index, enabled,
	)
	return err
}
