package logic

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"

	m "github.com/lesienchik/vk__test/internal/models"
	"github.com/lesienchik/vk__test/pkg/hashes"
	"github.com/lesienchik/vk__test/pkg/validator"
)

// Проводит валидацию полей пользователя при регистрации и проверяет на существование.
func (l *Logic) UserRegister(userReq *m.UserRegReq) *m.Err {
	// Валидация полей: username,email,password.
	if !validator.IsValidUsername(userReq.Username) {
		return &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Псевдоним не удовлетворяет требованиям",
		}
	}
	if !validator.IsValidEmail(userReq.Email) {
		return &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Неверный email",
		}
	}
	if !validator.IsValidPassword(userReq.Password) {
		return &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Пароль не удовлетворяет требованиям",
		}
	}

	// Проверяем пользователя на существование (по почте).
	_, exists, err := l.storage.User.GetByEmail(userReq.Email)
	if err != nil {
		return &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	if exists {
		return &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Пользователь с такой почтой уже существует",
			Error:     errors.New("email already exists"),
		}
	}

	// Проверяем пользователя на существование (по псевдониму).
	_, exists, err = l.storage.User.GetByUsername(userReq.Username)
	if err != nil {
		return &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	if exists {
		return &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Пользователь с таким псевдонимом уже существует",
			Error:     errors.New("username already exists"),
		}
	}

	verifyCode, err := hashes.HmacGenHash(userReq, hashes.ExpiresDefault, l.secret)
	if err != nil {
		return &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}

	// Не дожидаемся ответа об отправке сообщения на почту (дабы не заставлять пользователя ждать).
	go func() {
		if err := l.email.SendConfirmCode(userReq.Email, verifyCode); err != nil {
			l.logger.Error(fmt.Errorf("logic.UserRegister: %w", err))
			return
		}
	}()
	return nil
}

func (l *Logic) UserConfirm(verifyCode string) (int, *m.Err) {
	userReq := new(m.UserRegReq)
	hashingStatus, err := hashes.HmacParseAndValidateHash(verifyCode, userReq, l.secret)
	if err != nil {
		if hashingStatus == hashes.HashExpires {
			return -1, &m.Err{
				Code:      fasthttp.StatusBadRequest,
				ClientMsg: "Срок действия кода подтверждения истек",
				Error:     err,
			}
		}

		return -1, &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Неверный код подтверждения",
			Error:     err,
		}
	}

	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" {
		return -1, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     errors.New("empty fields for userReq"),
		}
	}

	// Проверяем пользователя на существование (по почте).
	_, exists, err := l.storage.User.GetByEmail(userReq.Email)
	if err != nil {
		return -1, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	if exists {
		return -1, &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Пользователь с такой почтой уже существует",
			Error:     errors.New("email already exists"),
		}
	}

	// Проверяем пользователя на существование (по псевдониму).
	_, exists, err = l.storage.User.GetByUsername(userReq.Username)
	if err != nil {
		return -1, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	if exists {
		return -1, &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Пользователь с таким псевдонимом уже существует",
			Error:     errors.New("username already exists"),
		}
	}

	userId, errs := l.userCreate(userReq)
	if errs != nil {
		return -1, errs
	}
	return userId, nil
}

// Сохраняет нового пользователя.
func (l *Logic) userCreate(userReq *m.UserRegReq) (int, *m.Err) {
	user := &m.User{
		Username: userReq.Username,
		Email:    userReq.Email,
	}

	hashPassword, err := hashes.HashPassword(userReq.Password)
	if err != nil {
		return -1, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	user.Password = string(hashPassword)

	userId, err := l.storage.User.Create(user)
	if err != nil {
		return -1, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	return userId, nil
}

// Возвращает пару access,refresh токенов.
func (l *Logic) UserSetJwtTokens(userId int) ([]string, *m.Err) {
	accessClaims := m.UserAuthClaims{
		Id: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiresAccessTime)),
		},
	}
	refreshClaims := m.UserAuthClaims{
		Id: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiresRefreshTime)),
		},
	}

	// Генерируем токены.
	access, err := hashes.JwtGenToken(accessClaims, l.secret)
	if err != nil {
		return nil, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	refresh, err := hashes.JwtGenToken(refreshClaims, l.secret)
	if err != nil {
		return nil, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	return []string{access, refresh}, nil
}

func (l *Logic) UserAuth(userReq *m.UserAuthReq) (int, *m.Err) {
	// Проверяем пользователя на существование (по почте).
	userDb, exists, err := l.storage.User.GetByEmail(userReq.Email)
	if err != nil {
		return -1, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	if !exists {
		return -1, &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Неверная почта или пароль",
			Error:     errors.New("invalid email or password"),
		}
	}

	if err := hashes.CompareHashAndPassword(userDb.Password, userReq.Password); err != nil {
		return -1, &m.Err{
			Code:      fasthttp.StatusBadRequest,
			ClientMsg: "Неверная почта или пароль",
			Error:     err,
		}
	}
	return userDb.Id, nil
}

func (l *Logic) UserVerify(token string) (int, *m.Err) {
	claims, _, err := hashes.JwtParseAndValidateToken(token, &m.UserAuthClaims{}, l.secret)
	if err != nil {
		return -1, &m.Err{
			Code:  fasthttp.StatusBadRequest,
			Error: errors.New("invalid access token"),
		}
	}

	// Приведение claims к типу UserAuthClaims.
	userAuthClaims, ok := claims.(*m.UserAuthClaims)
	if !ok {
		return -1, &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     errors.New("failed to convert claims to the UserAuthClaims type"),
		}
	}
	return userAuthClaims.Id, nil
}

func (l *Logic) UserRefresh(access, refresh string) (string, *m.Err) {
	accessClaims, tokenStatus, err := hashes.JwtParseAndValidateToken(access, &m.UserAuthClaims{}, l.secret)
	if err != nil {
		return "", &m.Err{
			Code:  fasthttp.StatusBadRequest,
			Error: errors.New("invalid access token"),
		}
	}
	if tokenStatus == hashes.JwtTokenValid {
		return "", &m.Err{
			Code:  fasthttp.StatusBadRequest,
			Error: errors.New("access token has not expired"),
		}
	}

	if _, _, err = hashes.JwtParseAndValidateToken(refresh, &m.UserAuthClaims{}, l.secret); err != nil {
		return "", &m.Err{
			Code:  fasthttp.StatusBadRequest,
			Error: errors.New("invalid refresh token"),
		}
	}

	// Приведение claims к типу UserAuthClaims.
	userAuthClaims, ok := accessClaims.(*m.UserAuthClaims)
	if !ok {
		return "", &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     errors.New("failed to convert access claims to the UserAuthClaims type"),
		}
	}

	claims := m.UserAuthClaims{
		Id: userAuthClaims.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiresAccessTime)),
		},
	}

	token, err := hashes.JwtGenToken(claims, l.secret)
	if err != nil {
		return "", &m.Err{
			Code:      fasthttp.StatusInternalServerError,
			ClientMsg: msgInternalServerError,
			Error:     err,
		}
	}
	return token, nil
}
