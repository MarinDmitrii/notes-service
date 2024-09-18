package usecase

import (
	"context"

	"github.com/MarinDmitrii/notes-service/internal/user/domain"
)

type GetUserByEmailUseCase struct {
	userRepository domain.Repository
}

func NewGetUserByEmailUseCase(userRepository domain.Repository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{userRepository: userRepository}
}

func (uc *GetUserByEmailUseCase) Execute(ctx context.Context, email string) (domain.User, error) {
	return uc.userRepository.GetUserByEmail(ctx, email)
}
