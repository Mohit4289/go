package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"gin-quickstart/repository"
	"time"

	"os"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
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

func (s *UserService) ValidateUser(ctx context.Context, email string) error {

	_, err := s.userRepo.FindUserByEmail(ctx, email)

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

	user, err := s.userRepo.CreateUser(ctx, name, email, password)

	if err != nil {
		return repository.User{}, err
	}

	return user, nil
}

func (s *UserService) CheckPassword(ctx context.Context, email string, password string) (string, string, error) {
	data, err := s.userRepo.VerifyPassword(ctx, email)
	if err != nil {
		return "", "", ValidationError{
			Field: "credentials",
			Msg:   "invalid email or password",
		}
	}

	verifyPass, err := argon2id.ComparePasswordAndHash(password, data.Password)
	if err != nil {
		return "", "", err
	}

	if !verifyPass {
		return "", "", ValidationError{
			Field: "Credentials",
			Msg:   "Wrong password",
		}
	}

	secret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", "", errors.New("Token not configuered")
	}

	claims := jwt.MapClaims{
		"user_id": data.ID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", errors.New("Not able to generate Token")
	}

	generateToken := make([]byte, 32)

	if _, err := rand.Read(generateToken); err != nil {
		return "", "", ValidationError{
			Field: "Refresh token",
			Msg:   "not able to generate token for refresh",
		}
	}

	stringToken := base64.RawURLEncoding.EncodeToString(generateToken)
	hash := sha256.Sum256([]byte(stringToken))
	tokenHash := hex.EncodeToString(hash[:])

	addingToken, err := s.userRepo.AddRefreshToken(ctx, tokenHash, email)
	if addingToken == false {
		return "", "", ValidationError{
			Field: "DB error",
			Msg:   "Not added refresh token",
		}
	}

	return signedToken, stringToken, nil
}

func (s *UserService) FetchUser(ctx context.Context, id int) (repository.User, error) {
	UserData, err := s.userRepo.FetchData(ctx, id)
	if err != nil {
		return repository.User{}, ValidationError{
			Field: "data not fetched",
			Msg:   err.Error(),
		}
	}
	return UserData, nil
}

func (s *UserService) VerfiyToken(ctx context.Context, token string) (repository.User, error) {
	UserData, err := s.userRepo.VerfiyToken(ctx, token)
	if err != nil {
		return repository.User{}, err
	}

	return UserData, nil
}

func (s *UserService) LogoutRemoveToken(ctx context.Context, hashtoken string) (bool, error) {
	DeleteToken, err := s.userRepo.RemoveToken(ctx, hashtoken)
	if err != nil {
		return false, ValidationError{
			Field: "DB Error",
			Msg:   "Failed to Delete Token",
		}
	}

	return DeleteToken, nil
}
