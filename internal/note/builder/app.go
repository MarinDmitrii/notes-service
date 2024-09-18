package builder

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarinDmitrii/notes-service/internal/note/adapters"
	"github.com/MarinDmitrii/notes-service/internal/note/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Application struct {
	SaveNote *usecase.SaveNoteUseCase
	GetNotes *usecase.GetNotesUseCase
}

type PostgresConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func NewPostgresConfig() *PostgresConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	return &PostgresConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

func NewApplication(ctx context.Context) (*Application, func()) {
	PostgresConfig := NewPostgresConfig()
	postgresConnect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		PostgresConfig.host,
		PostgresConfig.port,
		PostgresConfig.user,
		PostgresConfig.password,
		PostgresConfig.dbname,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", postgresConnect)
	if err != nil {
		panic(err)
	}

	noteRepository := adapters.NewPostgresNoteRepository(db)

	return &Application{
		SaveNote: usecase.NewSaveNoteUseCase(noteRepository),
		GetNotes: usecase.NewGetNotesUseCase(noteRepository),
	}, func() { db.Close() }
}
