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
	Admin *bool `json:"admin,omitempty"`
}

type User struct {
	UID   string
	Admin bool
}

func (a Auth) NewToken(user User) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	if user.Admin {
		claims.Admin = &user.Admin
	}

	token := jwt.NewWithClaims(a.signingMethod, claims)
	return token.SignedString(a.key)
}

// ParseToken parses and verifies a token.
func (a Auth) ParseToken(tokenString string) (User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err != nil {
		return User{}, err
	}
	if !token.Valid {
		return User{}, ErrInvalidToken
	}

	claims := token.Claims.(Claims)
	return User{
		UID:   claims.Subject,
		Admin: claims.Admin != nil && *claims.Admin,
	}, nil
}
