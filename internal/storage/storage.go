package storage

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type Storage struct {
	User User
}

func New(logger *logrus.Logger, db *sql.DB) *Storage {
	return &Storage{
		User: NewUser(logger, db),
	}
}
