package service

import (
	"context"
	"errors"
	"gin-quickstart/repository"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
)

type User struct {
	name  string
	email string
	pass  string
}

type ValidationError struct {
	Field string
	Msg   string
}

type UserService struct {
	userRepo *repository.UserRepo
}

func (e ValidationError) Error() string {
	return e.Msg
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) ValidateUser(ctx context.Context, email string) error {

	_, err := s.userRepo.FindUserByEmail(email)

	if err == nil {
		return ValidationError{
			Field: "email",
			Msg:   "user already exists",
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	return err
}

func (s *UserService) AddUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (repository.User, error) {

	user, err := s.userRepo.CreateUser(name, email, password)

	if err != nil {
		return repository.User{}, err
	}

	return user, nil
}

func (s *UserService) CheckPassword(ctx context.Context, email string, password string) (bool, error) {
	data, err := s.userRepo.VerifyPassword(ctx, email)
	if err != nil {
		return false, ValidationError{
			Field: "credentials",
			Msg:   "invalid email or password",
		}
	}

	verifyPass, err := argon2id.ComparePasswordAndHash(password, data)
	if err != nil {
		return false, err
	}

	return verifyPass, nil
}
