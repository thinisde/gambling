package auth

import (
	"encoding/json"
	"net/http"

	"backend/database/repo"
	"backend/lib"
)

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerHandler struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	user, err := repo.Users.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if ok := lib.CheckPasswordHash(req.Password, user.Password); !ok {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := signJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerHandler

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	userID, err := repo.Users.CreateUser(req.Username, req.DisplayName, req.Password)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := signJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
