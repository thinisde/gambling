package auth

import (
	"context"
	"fmt"
	"net/http"
)

type contenxtKey string

const UserIDKey contenxtKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := ReadTokenCookie(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		claims, err := verifyJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		userId, err := claims.UserID()
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
