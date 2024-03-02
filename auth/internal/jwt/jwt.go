package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/OliviaDilan/async_arc/auth/internal/user"

	"github.com/golang-jwt/jwt"
)

type Service interface {
	CreateToken(user *user.User) (string, error)
	DecodeToken(token string) (*Claims, error)
}

type Claims struct {
	Username string
	Role     string
	jwt.StandardClaims
}

type jwtService struct{
	secret []byte
}

func NewService(secret string) Service {
	return &jwtService{
		secret: []byte(secret),
	}
}

func (r *jwtService) CreateToken(user *user.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(r.secret)
}

var ErrInvalidToken = fmt.Errorf("invalid token")
func (r *jwtService) DecodeToken(token string) (*Claims, error) {
	claims := &Claims{}

	token = strings.TrimPrefix(token, "Bearer ")

	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return r.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}