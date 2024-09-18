package usecase

import (
	"context"

	"github.com/MarinDmitrii/notes-service/internal/user/domain"
)

type GetUserByIdUseCase struct {
	userRepository domain.Repository
}

func NewGetUserByIdUseCase(userRepository domain.Repository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{userRepository: userRepository}
}

func (uc *GetUserByIdUseCase) Execute(ctx context.Context, userId int) (domain.User, error) {
	return uc.userRepository.GetUserById(ctx, userId)
}
