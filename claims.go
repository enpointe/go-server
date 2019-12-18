package server

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CustomClaims the JWS Claims structure used to authenticate a user
type CustomClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	IssuedAt int64  `json:"iat"`
	jwt.StandardClaims
}

// GenerateToken generate the JWT Authentication token
func GenerateToken(username string, isAdmin bool, expires int, signingKey []byte) (string, error) {
	iat := time.Now()
	expirationTime := iat.Add(time.Duration(expires) * time.Second)
	// Create the JWT claims, which includes the username and expiry time
	claims := &CustomClaims{
		Username: username,
		IsAdmin:  isAdmin,
		IssuedAt: iat.Unix(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string.
	return token.SignedString(signingKey)
}

// ParseClaims from the given JWT token retrieve the custom Claims structure
func ParseClaims(tokenStr string, signingKey []byte) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
