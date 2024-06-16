package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token error")

type Auth struct {
	key           []byte
	signingMethod jwt.SigningMethod
}

func NewAuth(key string) Auth {
	return Auth{
		key:           []byte(key),
		signingMethod: jwt.SigningMethodHS256,
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (a Auth) NewToken(uID string) (string, error) {
	token := jwt.NewWithClaims(a.signingMethod, Claims{
		jwt.RegisteredClaims{
			Subject:   uID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return token.SignedString(a.key)
}

// ParseToken parses and verifies a token.
func (a Auth) ParseToken(tokenString string) (uID string, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrInvalidToken
	}
	return token.Claims.GetSubject()
}
