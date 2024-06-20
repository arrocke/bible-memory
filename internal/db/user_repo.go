package db

import (
	"context"
	"errors"
	"main/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
    Pool *pgxpool.Pool
}

var NotFoundError = errors.New("not found")

func userNamedArgs(u model.User) pgx.NamedArgs {
    return pgx.NamedArgs{
        "id": u.Id,
        "email": u.Email,
        "first_name": u.FirstName,
        "last_name": u.LastName,
    }
}

func (r UserRepo) GetUserByEmail(c context.Context, email string) (model.User, error) {
    query := `
        SELECT id, first_name, last_name, email, password
        FROM "user"
        WHERE email = $1
    `

    rows, _ := r.Pool.Query(context.Background(), query, email)
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return model.User{}, NotFoundError
        } else {
            return model.User{}, err
        }
    }

    return user, nil
}

func (r UserRepo) GetUserById(c context.Context, id int) (model.User, error) {
    query := `
        SELECT id, first_name, last_name, email, password
        FROM "user"
        WHERE id = $1
    `

    rows, _ := r.Pool.Query(context.Background(), query, id)
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return model.User{}, NotFoundError
        } else {
            return model.User{}, err
        }
    }

    return user, nil
}

func (r UserRepo) Update(c context.Context, user model.User) error {
    query := `
        UPDATE "user" SET 
            email = @email,
            first_name = @first_name,
            last_name = @last_name
        WHERE id = @id
    `
    _, err := r.Pool.Exec(c, query, userNamedArgs(user))
	return err
}
