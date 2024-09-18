package usecase

import (
	"context"

	"github.com/MarinDmitrii/notes-service/internal/note/domain"
)

type GetNotesUseCase struct {
	noteRepository domain.Repository
}

func NewGetNotesUseCase(noteRepository domain.Repository) *GetNotesUseCase {
	return &GetNotesUseCase{noteRepository: noteRepository}
}

func (uc *GetNotesUseCase) Execute(ctx context.Context, userID int) ([]domain.Note, error) {
	return uc.noteRepository.GetNotes(ctx, userID)
}
