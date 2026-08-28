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

var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrInvalidCredentials       = errors.New("invalid email or password")
	ErrUserNotFound             = errors.New("user not found")
	ErrTokenNotFound            = errors.New("refresh token not found")
	ErrJWTSecretNotConfigured   = errors.New("JWT secret not configured")
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
		return ErrUserAlreadyExists
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
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	verifyPass, err := argon2id.ComparePasswordAndHash(password, data.Password)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if !verifyPass {
		return "", "", ErrInvalidCredentials
	}

	secret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", "", ErrJWTSecretNotConfigured
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
		return "", "", errors.New("failed to generate access token")
	}

	generateToken := make([]byte, 32)
	if _, err := rand.Read(generateToken); err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	stringToken := base64.RawURLEncoding.EncodeToString(generateToken)
	hash := sha256.Sum256([]byte(stringToken))
	tokenHash := hex.EncodeToString(hash[:])

	addingToken, err := s.userRepo.AddRefreshToken(ctx, tokenHash, email)
	if err != nil {
		return "", "", err
	}
	if !addingToken {
		return "", "", errors.New("failed to save refresh token")
	}

	return signedToken, stringToken, nil
}

func (s *UserService) FetchUser(ctx context.Context, id int) (repository.User, error) {
	userData, err := s.userRepo.FetchData(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.User{}, ErrUserNotFound
		}
		return repository.User{}, err
	}
	return userData, nil
}

func (s *UserService) VerfiyToken(ctx context.Context, token string) (repository.User, error) {
	userData, err := s.userRepo.VerfiyToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.User{}, ErrTokenNotFound
		}
		return repository.User{}, err
	}

	return userData, nil
}

func (s *UserService) LogoutRemoveToken(ctx context.Context, hashtoken string) (bool, error) {
	deleteToken, err := s.userRepo.RemoveToken(ctx, hashtoken)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return false, ErrTokenNotFound
		}
		return false, err
	}

	return deleteToken, nil
}
