package storage

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	m "github.com/lesienchik/vk__test/internal/models"
)

type User interface {
	// Create info
	Create(user *m.User) (int, error)

	// Update info
	UpdatePasswordById(id int, newPassword string) error

	// Get info
	GetById(userId int) (*m.User, bool, error)
	GetByEmail(email string) (*m.User, bool, error)
	GetByUsername(username string) (*m.User, bool, error)
}

type user struct {
	logger *logrus.Logger
	db     *sql.DB
}

func NewUser(logger *logrus.Logger, db *sql.DB) *user {
	return &user{
		logger: logger,
		db:     db,
	}
}

func (u *user) Create(user *m.User) (int, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("storage.User.Create(1): %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	if err := tx.QueryRow(query, user.Username, user.Email, user.Password).Scan(&id); err != nil {
		return -1, fmt.Errorf("storage.User.Create(2): %w", err)
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("storage.User.Create(3): %w", err)
	}

	return id, nil
}

func (u *user) GetById(userId int) (*m.User, bool, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password
		FROM users WHERE id = $1
	`

	var user m.User
	if err := u.db.QueryRow(query, userId).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		if err != sql.ErrNoRows {
			return nil, false, fmt.Errorf("storage.User.GetById(1): %w", err)
		}
		return nil, false, nil
	}
	return &user, true, nil
}

func (u *user) GetByEmail(email string) (*m.User, bool, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password
		FROM users WHERE email = $1
	`

	var user m.User
	if err := u.db.QueryRow(query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		if err != sql.ErrNoRows {
			return nil, false, fmt.Errorf("storage.User.GetByEmail(1): %w", err)
		}
		return nil, false, nil
	}
	return &user, true, nil
}

func (u *user) GetByUsername(username string) (*m.User, bool, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password
		FROM users WHERE username = $1
	`

	var user m.User
	if err := u.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		if err != sql.ErrNoRows {
			return nil, false, fmt.Errorf("storage.User.GetByUsername(1): %w", err)
		}
		return nil, false, nil
	}
	return &user, true, nil
}

func (u *user) UpdatePasswordById(id int, newPassword string) error {
	query := `
		UPDATE users
		SET password = $2
		WHERE id = $1
	`

	tx, err := u.db.Begin()
	if err != nil {
		return fmt.Errorf("storage.UpdatePassword(1): %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(query, id, newPassword); err != nil {
		return fmt.Errorf("storage.UpdatePassword(2): %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("storage.UpdatePassword(3): %w", err)
	}
	return nil
}
