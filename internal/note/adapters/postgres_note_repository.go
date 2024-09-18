package adapters

import (
	"context"
	"time"

	domain "github.com/MarinDmitrii/notes-service/internal/note/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresNoteRepository struct {
	db *sqlx.DB
}

func NewPostgresNoteRepository(db *sqlx.DB) *PostgresNoteRepository {
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS "notes" (
		"id" serial NOT NULL,
		"user_id" integer NOT NULL,
		"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"description" varchar(10000) NOT NULL,
		CONSTRAINT "notes_pk" PRIMARY KEY ("id"),
		CONSTRAINT "notes_fk0" FOREIGN KEY ("user_id") REFERENCES "users"("id")
	);
	`)

	return &PostgresNoteRepository{db: db}
}

func (r *PostgresNoteRepository) SaveNote(ctx context.Context, domainNote domain.Note) (int, error) {
	if domainNote.ID == 0 {
		err := r.db.QueryRowContext(ctx, "SELECT nextval('notes_id_seq')").Scan(&domainNote.ID)
		if err != nil {
			return 0, err
		}

		domainNote.CreateDt = time.Now()
	}

	query := `
		INSERT INTO notes (id, user_id, created_at, description)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET user_id = EXCLUDED.user_id, created_at = EXCLUDED.created_at, description = EXCLUDED.description
		RETURNING id
	`

	var noteID int
	err := r.db.QueryRowContext(ctx, query, domainNote.ID, domainNote.UserID, domainNote.CreateDt, domainNote.Description).Scan(&noteID)
	if err != nil {
		return 0, err
	}

	return noteID, nil
}

func (r *PostgresNoteRepository) GetNotes(ctx context.Context, userID int) ([]domain.Note, error) {
	query := `SELECT * FROM notes WHERE user_id = $1`

	var noteModels []NoteModel
	err := r.db.SelectContext(ctx, &noteModels, query, userID)
	if err != nil {
		return nil, err
	}

	notes := make([]domain.Note, len(noteModels))
	for i, model := range noteModels {
		notes[i] = model.mapToDomain()
	}

	return notes, nil
}

func (r *PostgresNoteRepository) GetNoteById(ctx context.Context, noteID int) (domain.Note, error) {
	query := `SELECT * FROM notes WHERE id = $1`

	var noteModel NoteModel
	err := r.db.GetContext(ctx, &noteModel, query, noteID)
	if err != nil {
		return domain.Note{}, err
	}

	note := noteModel.mapToDomain()

	return note, nil
}
