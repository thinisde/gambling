package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Ver string `json:"ver"`
	jwt.RegisteredClaims
}

func jwtKey() ([]byte, error) {
	k := os.Getenv("JWT_KEY")
	if k == "" {
		return nil, errors.New("JWT_KEY is not set")
	}
	return []byte(k), nil
}

func signJWT(userID int64) (string, error) {
	key, err := jwtKey()
	if err != nil {
		return "", err
	}

	claims := Claims{
		Ver: os.Getenv("JWT_VER"),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Gambling Server",
			Subject:   strconv.FormatInt(userID, 10), // keep sub as string
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(key)
}

func verifyJWT(signed string) (*Claims, error) {
	key, err := jwtKey()
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		signed,
		claims,
		func(t *jwt.Token) (any, error) {
			// strictly enforce HS256 (don’t just accept any HMAC)
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return key, nil
		},
		// optional but recommended:
		jwt.WithIssuer("Gambling Server"),
		// jwt.WithLeeway(30*time.Second), // if you want small clock skew tolerance
	)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Optional: enforce your JWT_VER
	expectedVer := os.Getenv("JWT_VER")
	if expectedVer != "" && claims.Ver != expectedVer {
		return nil, errors.New("invalid token version")
	}

	return claims, nil
}

// Helper to get userID back as int64
func (c *Claims) UserID() (int64, error) {
	return strconv.ParseInt(c.Subject, 10, 64)
}
