package usecase

import (
	"context"

	domain "github.com/MarinDmitrii/notes-service/internal/note/domain"
)

type SaveNote struct {
	UserID      int
	Description string
}

type SaveNoteUseCase struct {
	noteRepository domain.Repository
}

func NewSaveNoteUseCase(
	noteRepository domain.Repository,
) *SaveNoteUseCase {
	return &SaveNoteUseCase{
		noteRepository: noteRepository,
	}
}

func (uc *SaveNoteUseCase) Execute(ctx context.Context, createdNote SaveNote) (domain.Note, error) {
	noteID, err := uc.noteRepository.SaveNote(ctx, domain.Note{
		UserID:      createdNote.UserID,
		Description: createdNote.Description,
	})
	if err != nil {
		return domain.Note{}, err
	}

	note, err := uc.noteRepository.GetNoteById(ctx, noteID)
	if err != nil {
		return domain.Note{}, err
	}

	return note, nil
}
