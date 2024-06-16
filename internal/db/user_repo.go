package db

import (
	"context"
	"errors"
	"fmt"
	"main/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
    Pool *pgxpool.Pool
}

var NotFoundError = errors.New("not found")

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
    fmt.Printf("%v %v %v %v %v %v", user.Id, user.FirstName, user.LastName, user.Email, user.Password)

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
