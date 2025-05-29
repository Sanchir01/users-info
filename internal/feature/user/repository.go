package user

import (
	"context"

	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/users-info/internal/gender"
	api "github.com/Sanchir01/users-info/pkg/lib/api/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateUserRequestDB struct {
	Name        *string // pointer, чтобы можно было отличить "не передано" от пустой строки
	Surname     *string
	Patronymic  *string
	Nationality *string
	Age         *int
	Gender      *gender.Gender
}
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

func (r *Repository) GetAllUsers(ctx context.Context) ([]*UserDB, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	query, args, err := sq.Select("id,name, surname,patronymic,created_at,updated_at,age,gender,nationality,version").
		From("public.users").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, api.ErrQueryString
	}
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*UserDB

	for rows.Next() {
		var oneuserdb UserDB
		if err := rows.Scan(
			&oneuserdb.ID,
			&oneuserdb.Name,
			&oneuserdb.Surname,
			&oneuserdb.Patronymic,
			&oneuserdb.CreatedAt,
			&oneuserdb.UpdatedAt,
			&oneuserdb.Age,
			&oneuserdb.Gender,
			&oneuserdb.Nationality,
			&oneuserdb.Version,
		); err != nil {
			return nil, err
		}
		users = append(users, &oneuserdb)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	req UpdateUserRequestDB,
	tx pgx.Tx,
) error {
	updateBuilder := sq.Update("users").Where(sq.Eq{"id": id})

	if req.Name != nil {
		updateBuilder = updateBuilder.Set("name", *req.Name)
	}
	if req.Surname != nil {
		updateBuilder = updateBuilder.Set("surname", *req.Surname)
	}
	if req.Patronymic != nil {
		updateBuilder = updateBuilder.Set("patronymic", *req.Patronymic)
	}
	if req.Nationality != nil {
		updateBuilder = updateBuilder.Set("nationality", *req.Nationality)
	}
	if req.Age != nil {
		updateBuilder = updateBuilder.Set("age", *req.Age)
	}
	if req.Gender != nil {
		updateBuilder = updateBuilder.Set("gender", *req.Gender)
	}

	updateBuilder = updateBuilder.Set("updated_at", sq.Expr("NOW()"))

	query, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return api.ErrQueryString
	}

	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return api.ErrNotFoundById
	}
	return nil
}

func (r *Repository) DeleteUserById(ctx context.Context, id uuid.UUID, tx pgx.Tx) error {

	query, args, err := sq.Delete("users").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return api.ErrQueryString
	}
	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return api.ErrNotFoundById
	}
	return nil
}
