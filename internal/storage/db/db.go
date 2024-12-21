package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Ord1nI/netStats/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	DB *sql.DB
}

func NewDb() (*DB, error) {
	db, err := sql.Open("sqlite3","stats.db")
	db.SetMaxOpenConns(1)

	if err != nil {
		return nil, err
	}
	return &DB{db}, err
}

func (d *DB) Add(stats []storage.Stat, name string) error {
	_, err := d.DB.Exec("INSERT INTO snapshots (name) VALUES ('Snapshot 1');")
	if err != nil {
		return err
	}

	sid := d.DB.QueryRow("SELECT id FROM snapshots ORDER BY created_at DESC LIMIT 1;")

	var id int64

	err = sid.Scan(&id)
	if err != nil {
		return err
	}

	childCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := d.DB.BeginTx(childCtx, nil)
	if err != nil {
		return err
	}

	for _, v := range stats {

		js, err := json.Marshal(v)

		_, err = tx.Exec("INSERT INTO snapshot_data (snapshot_id, data) VALUES ($1, $2);", id, js)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *DB) Close() error {
	return db.DB.Close()
}
