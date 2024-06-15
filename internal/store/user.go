package store

import "github.com/TheMangoMen/backend/internal/model"

func (s *Store) GetUser(uID string) (user model.User, err error) {
	err = s.db.Select(&user, "SELECT * FROM Users WHERE UID = $1", uID)
	return
}

func (s *Store) CreateUser(uID string) error {
	_, err := s.db.Exec("INSERT INTO Users (UID) VALUES ($1)", uID)
	return err
}
