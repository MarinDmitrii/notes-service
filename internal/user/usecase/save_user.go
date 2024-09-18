package usecase

import (
	"context"

	domain "github.com/MarinDmitrii/notes-service/internal/user/domain"
)

type SaveUser struct {
	Email    string
	Password string
}

type SaveUserUseCase struct {
	userRepository domain.Repository
}

func NewSaveUserUseCase(
	userRepository domain.Repository,
) *SaveUserUseCase {
	return &SaveUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *SaveUserUseCase) Execute(ctx context.Context, createdUser SaveUser) (domain.User, error) {
	userID, err := uc.userRepository.SaveUser(ctx, domain.User{
		Email:    createdUser.Email,
		Password: createdUser.Password,
	})
	if err != nil {
		return domain.User{}, err
	}

	user, err := uc.userRepository.GetUserById(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
