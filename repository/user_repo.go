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
	ID       int64
	Password string
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
	row := r.DB.QueryRow(ctx, `SELECT id, password FROM public."user" WHERE email = $1`,
		email,
	)

	var userPass UserPass
	err := row.Scan(
		&userPass.ID,
		&userPass.Password,
	)
	if err != nil {
		return UserPass{}, err
	}

	return userPass, nil

}

func (r *UserRepo) FetchData(ctx context.Context, id int) (User, error) {

	row := r.DB.QueryRow(ctx, `SELECT id, name, email FROM public."user" WHERE id = $1`, id)

	var UserData User
	err := row.Scan(
		&UserData.ID,
		&UserData.Name,
		&UserData.Email,
	)
	if err != nil {
		return User{}, err
	}
	return UserData, nil
}

func (r *UserRepo) AddRefreshToken(ctx context.Context, refresh_token string, email string) (bool, error) {
	_, err := r.DB.Exec(
		ctx,
		`UPDATE "user"
         SET refresh_token = $1
         WHERE email = $2`,
		refresh_token, email,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepo) VerfiyToken(ctx context.Context, refresh_token string) (User, error) {
	row := r.DB.QueryRow(ctx, `SELECT id, name, email FROM public."user" WHERE refresh_token = $1`, refresh_token)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
