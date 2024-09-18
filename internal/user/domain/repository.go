package domain

import (
	"context"
)

type Repository interface {
	SaveUser(ctx context.Context, user User) (int, error)
	GetUserById(ctx context.Context, userId int) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
