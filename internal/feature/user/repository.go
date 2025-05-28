package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/users-info/internal/gender"
	api "github.com/Sanchir01/users-info/pkg/lib/api/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

func (r *Repository) CreateUserRepository(
	name, surname, patronymic, nationality string,
	age int, gender gender.Gender,
	tx pgx.Tx, ctx context.Context,
) error {
	query, args, err := sq.Insert("users").
		Columns("name", "surname", "patronymic", "nationality", "age", "gender").
		Values(name, surname, patronymic, nationality, age, gender).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return api.ErrQueryString
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
