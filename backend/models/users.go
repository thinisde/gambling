package models

import "time"

type User struct {
	ID          int64
	Username    string
	DisplayName string
	Password    string
	Balance     float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
