package logic

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/lesienchik/vk__test/internal/storage"
	"github.com/lesienchik/vk__test/pkg/email"
)

// Популярные клиентские сообщения для ошибок.
var (
	msgInternalServerError = "Упс! Что-то пошло не так..."
)

const (
	jwtExpiresAccessTime  = 10 * time.Minute
	jwtExpiresRefreshTime = 24 * time.Hour * 30
)

type Logic struct {
	secret  string
	logger  *logrus.Logger
	email   *email.Email
	storage *storage.Storage
}

func New(secret string, logger *logrus.Logger, email *email.Email, storage *storage.Storage) *Logic {
	return &Logic{
		secret:  secret,
		logger:  logger,
		email:   email,
		storage: storage,
	}
}
