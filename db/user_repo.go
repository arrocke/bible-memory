package db

import (
	"context"
	"main/domain_model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
    Get(id int) (*domain_model.User, error)
    GetByEmail(email string) (*domain_model.User, error)
    Create(*domain_model.User) error
    Update(*domain_model.User) error
}

type PgUserRepo struct {
    pool *pgxpool.Pool
}

func CreatePgUserRepo(pool *pgxpool.Pool) PgUserRepo {
    return PgUserRepo{pool}
}

type userModel struct {
    Id int
    FirstName string
    LastName string
    Email string
    Password string
}

func (dbModel *userModel) toDomain() domain_model.User {
    user := domain_model.User{
        Id: dbModel.Id,
        FirstName: dbModel.FirstName,
        LastName: dbModel.LastName,
        EmailAddress: dbModel.Email,
        Password: dbModel.Password,
    }

    return user
}

func userToDb(user *domain_model.User) userModel {
    dbModel := userModel{
        Id: user.Id,
        FirstName: user.FirstName,
        LastName: user.LastName,
        Email: user.EmailAddress,
        Password: user.Password,
    }

    return dbModel
}

func (repo PgUserRepo) Get(id int) (*domain_model.User, error) {
	query := `
        SELECT id, first_name, last_name, email, password
        FROM user
        WHERE id = $1
    `
	rows, _ := repo.pool.Query(context.Background(), query, id)
	dbModel, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[userModel])
	if err != nil {
		return nil, err
	}

	user := dbModel.toDomain()
	return &user, nil
}

func (repo PgUserRepo) GetByEmail(email string) (*domain_model.User, error) {
	query := `
        SELECT id, first_name, last_name, email, password
        FROM user
        WHERE email = $1
    `
	rows, _ := repo.pool.Query(context.Background(), query, email)
	dbModel, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[userModel])
	if err != nil {
		return nil, err
	}

	user := dbModel.toDomain()
	return &user, nil
}

func (repo PgUserRepo) Create(user *domain_model.User) error {
	query := `
        INSERT INTO passage (id, first_name, last_name, email, password)
        VALUES ($1, $2, $3, $4, $5)
    `
	dbModel := userToDb(user)
	_, err := repo.pool.Exec(
		context.Background(),
        query,
		dbModel.Id,
		dbModel.FirstName,
		dbModel.LastName,
		dbModel.Email,
		dbModel.Password,
	)
	return err
}

func (repo PgUserRepo) Update(user *domain_model.User) error {
	query := `
        UPDATE passage
            SET first_name = $2,
            last_name = $3,
            email = $4,
            password = $5
        WHERE id = $1
    `
	dbModel := userToDb(user)
	_, err := repo.pool.Exec(
		context.Background(),
        query,
		dbModel.Id,
		dbModel.FirstName,
		dbModel.LastName,
		dbModel.Email,
		dbModel.Password,
	)
	return err
}
