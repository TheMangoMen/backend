package store

import (
	"database/sql"
	"errors"

	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetUser(uID string) (user model.User, err error) {
	err = s.db.Get(&user, "SELECT * FROM Users WHERE UID = $1;", uID)
	return
}

func (s *Store) CreateUser(uID string) error {
	_, err := s.db.Exec("INSERT INTO Users (UID) VALUES ($1) ON CONFLICT DO NOTHING;", uID)
	return err
}

func (s *Store) GetIsAdmin(uID string) (bool, error) {
	var a int
	err := s.db.Get(&a, "SELECT 1 FROM Admins WHERE UID = $1", uID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}
