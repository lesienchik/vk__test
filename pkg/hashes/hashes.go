package hashes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	ExpiresDefault    = 10 * time.Minute
	ExpiresOneMinute  = 1 * time.Minute
	ExpiresFiveMinute = 5 * time.Minute
	ExpiresTenMinute  = 10 * time.Minute
)

const (
	JwtTokenValid byte = iota
	JwtTokenError
	JwtTokenExpires
	HashValid
	HashError
	HashExpires
)

func HashPassword(password string) ([]byte, error) {
	if len(password) == 0 {
		return nil, fmt.Errorf("hashes.HashPassword(1): password is empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashes.HashPassword(2): %w", err)
	}
	return hash, nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	if len(hashedPassword) == 0 {
		return fmt.Errorf("hashes.CompareHashAndPassword(1): hashedPassword is empty")
	}
	if len(password) == 0 {
		return fmt.Errorf("hashes.CompareHashAndPassword(2): password is empty")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("hashes.CompareHashAndPassword(3): %w", err)
	}
	return nil
}

// Генерирует хэш из любых входных структур (с включенной сигнатурой json).
// Принимает время жизни генерируемого хэша.
func HmacGenHash(in any, ttl time.Duration, secret string) (string, error) {
	if ttl == 0 {
		ttl = ExpiresDefault
	}
	expiration := time.Now().Add(ttl).Unix()

	data, err := json.Marshal(in)
	if err != nil {
		return "", fmt.Errorf("hashes.HmacGenHash (1): %w", err)
	}

	message := fmt.Sprintf("%s|%d", data, expiration)

	hmacHash := hmac.New(sha256.New, []byte(secret))
	hmacHash.Write([]byte(message))
	signature := hmacHash.Sum(nil)

	// Используем RawStdEncoding для кодирования без паддинга
	hash := fmt.Sprintf("%s.%s",
		base64.RawStdEncoding.EncodeToString([]byte(message)),
		base64.RawStdEncoding.EncodeToString(signature),
	)
	return hash, nil
}

func HmacParseAndValidateHash(hash string, out any, secret string) (byte, error) {
	parts := strings.Split(hash, ".")
	if len(parts) != 2 {
		return HashError, errors.New("hashes.HmacParseAndValidateHash(1): invalid format hash")
	}

	message, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return HashError, fmt.Errorf("hashes.HmacParseAndValidateHash(2): %w", err)
	}

	signature, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return HashError, fmt.Errorf("hashes.HmacParseAndValidateHash(3): %w", err)
	}

	hmacHash := hmac.New(sha256.New, []byte(secret))
	hmacHash.Write(message)
	expectedSignature := hmacHash.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return HashError, errors.New("hashes.HmacParseAndValidateHash(4): invalid hash signature")
	}

	partsMessage := strings.SplitN(string(message), "|", 2)
	if len(partsMessage) != 2 {
		return HashError, errors.New("hashes.HmacParseAndValidateHash(5): invalid message format")
	}

	expiration, err := strconv.ParseInt(partsMessage[1], 10, 64)
	if err != nil || time.Now().Unix() > expiration {
		return HashExpires, errors.New("hashes.HmacParseAndValidateHash(6): hash has expired")
	}

	if err := json.Unmarshal([]byte(partsMessage[0]), &out); err != nil {
		return HashError, fmt.Errorf("hashes.HmacParseAndValidateHash(7): %w", err)
	}

	return HashValid, nil
}

// Генерирует jwt-токен, в зависимости от переданного claims.
func JwtGenToken(claims jwt.Claims, secret string) (token string, err error) {
	jwtT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtT.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("hashes.JwtGenToken(1): %w", err)
	}
	return token, nil
}

// Расшифровывает и проверяет JWT токен с любыми claims.
func JwtParseAndValidateToken(token string, claims jwt.Claims, secret string) (jwt.Claims, byte, error) {
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("hashes.JwtParseAndValidateToken(1): unexpected signing method [%v]", jwtToken.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		// Проверяем, не истек ли срок действия токена.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, JwtTokenExpires, fmt.Errorf("hashes.JwtParseAndValidateToken(2): token has expired")
		}
		return nil, JwtTokenError, fmt.Errorf("hashes.JwtParseAndValidateToken(3): %w", err)
	}

	if !jwtToken.Valid {
		return nil, JwtTokenError, fmt.Errorf("hashes.JwtParseAndValidateToken(4): invalid token")
	}

	return claims, JwtTokenValid, nil
}
