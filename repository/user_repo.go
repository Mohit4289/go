package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserPass struct {
	ID       int64
	Password string
}

type UserRepo struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func NewUserRepo(db *pgxpool.Pool, redis *redis.Client) *UserRepo {
	return &UserRepo{
		DB:    db,
		Redis: redis,
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

	userData, err := json.Marshal(user)
	if err != nil {
		return User{}, err
	}

	key := "user:" + strconv.Itoa(int(user.ID))

	err = r.Redis.Set(
		ctx,
		key,
		userData,
		10*time.Minute,
	).Err()

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

func (r *UserRepo) RemoveToken(ctx context.Context, refresh_token string) (bool, error) {
	res, err := r.DB.Exec(ctx, `UPDATE "user" SET refresh_token = NULL WHERE refresh_token = $1`, refresh_token)
	if err != nil {
		return false, err
	}

	if res.RowsAffected() == 0 {
		return false, ErrRefreshTokenNotFound
	}

	return true, nil
}
