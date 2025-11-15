package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func New() *Store {
	return &Store{}
}

func (s *Store) Open(config string) error {
	db, err := sql.Open("postgres", config)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Store) GetConn() *sql.DB {
	return s.db
}

func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		return
	}
}
func (s *Store) SetConn(db *sql.DB) {
	s.db = db
}
