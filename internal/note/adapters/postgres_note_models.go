package adapters

import (
	"time"

	domain "github.com/MarinDmitrii/notes-service/internal/note/domain"
)

type NoteModel struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	CreateDt    time.Time `db:"created_at"`
	Description string    `db:"description"`
}

func NewNoteModel(note domain.Note) (NoteModel, error) {
	return NoteModel{
		ID:          note.ID,
		UserID:      note.UserID,
		CreateDt:    note.CreateDt,
		Description: note.Description,
	}, nil
}

func (model *NoteModel) mapToDomain() domain.Note {
	return domain.Note{
		ID:          model.ID,
		UserID:      model.UserID,
		CreateDt:    model.CreateDt,
		Description: model.Description,
	}
}
