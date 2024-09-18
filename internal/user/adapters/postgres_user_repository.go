package adapters

import (
	"context"
	"time"

	domain "github.com/MarinDmitrii/notes-service/internal/user/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS "users" (
		"id" serial NOT NULL,
		"email" varchar(255) NOT NULL UNIQUE,
		"password" varchar(255) NOT NULL,
		"created_at" TIMESTAMP NOT NULL,
		"updated_at" TIMESTAMP NOT NULL,
		CONSTRAINT "users_pk" PRIMARY KEY ("id")
	);
	`)

	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SaveUser(ctx context.Context, domainUser domain.User) (int, error) {
	if domainUser.ID == 0 {
		err := r.db.QueryRowContext(ctx, "SELECT nextval('users_id_seq')").Scan(&domainUser.ID)
		if err != nil {
			return 0, err
		}

		domainUser.CreateDt = time.Now()
	}

	domainUser.UpdateDt = time.Now()

	query := `
	INSERT INTO users (id, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id) DO UPDATE
	SET email = EXCLUDED.email, password = EXCLUDED.password, created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at
	RETURNING id
	`

	var userID int
	err := r.db.QueryRowContext(ctx, query, domainUser.ID, domainUser.Email, domainUser.Password, domainUser.CreateDt, domainUser.UpdateDt).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *PostgresUserRepository) GetUserById(ctx context.Context, userID int) (domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var userModel UserModel
	err := r.db.GetContext(ctx, &userModel, query, userID)
	if err != nil {
		return domain.User{}, err
	}

	user := userModel.mapToDomain()

	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var userModel UserModel
	err := r.db.GetContext(ctx, &userModel, query, email)
	if err != nil {
		return domain.User{}, err
	}

	user := userModel.mapToDomain()

	return user, nil
}
