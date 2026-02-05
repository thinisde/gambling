package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"backend/cache"
	"backend/database"
	"backend/database/repo"

	"backend/handlers/auth"
	"backend/handlers/users"
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

	repo.Users = *repo.NewUserRepo(db)
	cache.CacheClient = cache.NewRedisClient()

	auth.CreateAuthHandler(r)
	users.CreateUsersHandler(r)

	http.ListenAndServe(":3000", r)
}
