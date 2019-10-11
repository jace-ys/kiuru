package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTConfig struct {
	SecretKey string
	Issuer    string
	TTL       time.Duration
}

type JWTClaims struct {
	*jwt.StandardClaims
	UserId   string
	Username string
}
