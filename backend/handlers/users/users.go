package users

import (
	"backend/database/repo"
	"backend/handlers/auth"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateUsersHandler(r *mux.Router) {
	users := r.PathPrefix("/users").Subrouter()
	users.Use(auth.AuthMiddleware)

	users.HandleFunc("/self", GetSelf).Methods("GET")
}

// CRUD operations for users
func GetSelf(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int64)

	user, err := repo.Users.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		return
	}
}
