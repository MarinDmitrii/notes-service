package adapters

import (
	"time"

	domain "github.com/MarinDmitrii/notes-service/internal/user/domain"
)

type UserModel struct {
	ID       int       `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	CreateDt time.Time `db:"created_at"`
	UpdateDt time.Time `db:"updated_at"`
}

func NewUserModel(user domain.User) (UserModel, error) {
	return UserModel{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		CreateDt: user.CreateDt,
		UpdateDt: user.UpdateDt,
	}, nil
}

func (model *UserModel) mapToDomain() domain.User {
	return domain.User{
		ID:       model.ID,
		Email:    model.Email,
		Password: model.Password,
		CreateDt: model.CreateDt,
		UpdateDt: model.UpdateDt,
	}
}
