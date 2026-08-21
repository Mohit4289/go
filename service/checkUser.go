package service

import (
	"errors"
	"gin-quickstart/repository"

	"github.com/jackc/pgx/v5"
)

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

func (s *UserService) ValidateUser(email string) error {

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
