package models

import (
	"time"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string
}

type UserActivity struct {
	ID     int
	UserID int
	Route  string
	Time   time.Time
}
