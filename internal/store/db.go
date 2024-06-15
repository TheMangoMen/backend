package store

import "github.com/jmoiron/sqlx"

type Store struct {
	db *sqlx.DB
}
