package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"backend/database"
	"backend/database/repo"
	"backend/handlers/auth"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	db, err := database.PostgresDriver.Open(database.CreateDnsPostgres())
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer database.PostgresDriver.Close()

	r := mux.NewRouter()

	// Run Migrations
	if err := database.PostgresDriver.Migrate(); err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	repo.Users = *repo.NewUserRepo(db)

	// Auth routes
	authR := r.PathPrefix("/auth").Subrouter()
	authR.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	authR.HandleFunc("/register", auth.RegisterHandler).Methods("POST")

	http.ListenAndServe(":3000", r)
}
