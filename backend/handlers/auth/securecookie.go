package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
)

var cookies = map[string]*securecookie.SecureCookie{
	"previous": securecookie.New(
		securecookie.GenerateRandomKey(64), // hash key
		securecookie.GenerateRandomKey(32), // block key (16/24/32)
	),
	"current": securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	),
}

const cookieName = "token"

func SetTokenCookie(w http.ResponseWriter, token string) error {
	value := map[string]string{"token": token}

	encoded, err := cookies["current"].Encode(cookieName, value)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production", // false on http://localhost
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(72 * time.Hour),
	})

	return nil
}

func ReadTokenCookie(r *http.Request) (string, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	// try current first
	if token, err := decodeToken(c.Value, cookies["current"]); err == nil {
		return token, nil
	}

	// fallback to previous
	if token, err := decodeToken(c.Value, cookies["previous"]); err == nil {
		return token, nil
	}

	return "", errors.New("invalid or expired cookie")
}

func decodeToken(raw string, sc *securecookie.SecureCookie) (string, error) {
	value := map[string]string{}
	if err := sc.Decode(cookieName, raw, &value); err != nil {
		return "", err
	}
	return value["token"], nil
}
