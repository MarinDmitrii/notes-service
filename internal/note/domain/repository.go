package domain

import (
	"context"
)

type Repository interface {
	SaveNote(ctx context.Context, note Note) (int, error)
	GetNotes(ctx context.Context, userID int) ([]Note, error)
	GetNoteById(ctx context.Context, noteId int) (Note, error)
}
