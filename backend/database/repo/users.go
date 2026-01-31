package repo

import (
	"backend/lib"
	"backend/models"
	"database/sql"
)

var Users UserRepo

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(username, displayName, password string) (int64, error) {
	var userID int64

	hashedPassword, err := lib.HashPassword(password)
	if err != nil {
		return 0, err
	}

	err = r.db.QueryRow(
		"INSERT INTO users (username, display_name, password, balance, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id",
		username, displayName, hashedPassword, 100,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepo) GetUserByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, display_name, password, balance, created_at FROM users WHERE username = $1",
		username,
	).Scan(&u.ID, &u.DisplayName, &u.Password, &u.Balance, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetUserByID(userID int64) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, username, display_name, balance, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&u.ID, &u.Username, &u.DisplayName, &u.Balance, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdateUserBalance(userID int64, newBalance float64) error {
	_, err := r.db.Exec(
		"UPDATE users SET balance = $1, updated_at = NOW() WHERE id = $2",
		newBalance, userID,
	)
	return err
}
