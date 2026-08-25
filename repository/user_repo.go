package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID    int64
	Name  string
	Email string
}

type UserPass struct {
	Id       int64
	Password string
}

type validationError struct {
	Field string
	Msg   string
}

type UserRepo struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) CreateUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (User, error) {

	var user User

	err := r.DB.QueryRow(
		ctx,
		`INSERT INTO "user" (name, email, password)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, email`,
		name,
		email,
		password,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (int64, error) {
	row := r.DB.QueryRow(
		ctx,
		`SELECT id FROM public."user" WHERE email = $1`,
		email,
	)

	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepo) VerifyPassword(ctx context.Context, email string) (UserPass, error) {
	row := r.DB.QueryRow(ctx, `SELECT id password FROM public."user" WHERE email = $1`,
		email,
	)

	var userPass UserPass
	err := row.Scan(&userPass)
	if err != nil {
		return userPass, err
	}

	return userPass, nil

}
