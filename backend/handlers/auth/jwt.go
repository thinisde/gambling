package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signJWT(userid int64) (string, error) {
	key := []byte(os.Getenv("JWT_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "Gambling Server",
		"sub": userid,
		"ver": os.Getenv("JWT_VER"),
		"exp": jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
	})
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func verifyJWT(signedToken string) (*jwt.Token, error) {
	key := []byte(os.Getenv("JWT_KEY"))
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return key, nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}
	return token, nil
}
