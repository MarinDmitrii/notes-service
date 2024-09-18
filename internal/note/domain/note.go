package domain

import "time"

type Note struct {
	ID          int
	UserID      int
	CreateDt    time.Time
	Description string
}
