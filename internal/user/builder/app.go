package builder

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarinDmitrii/notes-service/internal/user/adapters"
	"github.com/MarinDmitrii/notes-service/internal/user/usecase"
	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Application struct {
	SaveUser       *usecase.SaveUserUseCase
	GetUserById    *usecase.GetUserByIdUseCase
	GetUserByEmail *usecase.GetUserByEmailUseCase
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

	userRepository := adapters.NewPostgresUserRepository(db)

	return &Application{
		SaveUser:       usecase.NewSaveUserUseCase(userRepository),
		GetUserById:    usecase.NewGetUserByIdUseCase(userRepository),
		GetUserByEmail: usecase.NewGetUserByEmailUseCase(userRepository),
	}, func() { db.Close() }
}
