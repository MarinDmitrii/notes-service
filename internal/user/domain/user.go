package domain

import "time"

type User struct {
	ID       int
	Email    string
	Password string
	CreateDt time.Time
	UpdateDt time.Time
}
