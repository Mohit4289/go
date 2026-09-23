package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	db "gin-quickstart/db/sqlc"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrUserNotFound           = errors.New("user not found")
	ErrTokenNotFound          = errors.New("refresh token not found")
	ErrJWTSecretNotConfigured = errors.New("JWT secret not configured")
)

type ValidationError struct {
	Field string
	Msg   string
}

func (e ValidationError) Error() string {
	return e.Msg
}

type UserService struct {
	queries *db.Queries
	redis   *redis.Client
}

func NewUserService(queries *db.Queries, redis *redis.Client) *UserService {
	return &UserService{
		queries: queries,
		redis:   redis,
	}
}

func (s *UserService) ValidateUser(ctx context.Context, email string) error {
	_, err := s.queries.FindUserByEmail(ctx, email)
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
) (db.CreateUserRow, error) {
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return db.CreateUserRow{}, err
	}

	userData, err := json.Marshal(user)
	if err == nil && s.redis != nil {
		key := "user:" + strconv.FormatInt(user.ID, 10)
		_ = s.redis.Set(ctx, key, userData, 10*time.Minute).Err()
	}

	return user, nil
}

func (s *UserService) CheckPassword(ctx context.Context, email string, password string) (string, string, error) {
	data, err := s.queries.VerifyPassword(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	verifyPass, err := argon2id.ComparePasswordAndHash(password, data.Password)
	if err != nil || !verifyPass {
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

	rowsAffected, err := s.queries.AddRefreshToken(ctx, db.AddRefreshTokenParams{
		RefreshToken: pgtype.Text{String: tokenHash, Valid: true},
		Email:        email,
	})
	if err != nil {
		return "", "", err
	}
	if rowsAffected == 0 {
		return "", "", errors.New("failed to save refresh token")
	}

	return signedToken, stringToken, nil
}

func (s *UserService) FetchUser(ctx context.Context, id int) (db.FetchUserByIDRow, error) {
	key := "user:" + strconv.Itoa(id)

	if s.redis != nil {
		userData, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var user db.FetchUserByIDRow
			if err := json.Unmarshal([]byte(userData), &user); err == nil {
				return user, nil
			}
		}
	}

	user, err := s.queries.FetchUserByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.FetchUserByIDRow{}, ErrUserNotFound
		}
		return db.FetchUserByIDRow{}, err
	}

	if s.redis != nil {
		userDataBytes, err := json.Marshal(user)
		if err == nil {
			_ = s.redis.Set(ctx, key, userDataBytes, 10*time.Minute).Err()
		}
	}

	return user, nil
}

func (s *UserService) VerfiyToken(ctx context.Context, token string) (db.VerifyRefreshTokenRow, error) {
	user, err := s.queries.VerifyRefreshToken(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.VerifyRefreshTokenRow{}, ErrTokenNotFound
		}
		return db.VerifyRefreshTokenRow{}, err
	}

	return user, nil
}

func (s *UserService) LogoutRemoveToken(ctx context.Context, hashtoken string) (bool, error) {
	rowsAffected, err := s.queries.RemoveRefreshToken(ctx, pgtype.Text{String: hashtoken, Valid: true})
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, ErrTokenNotFound
	}

	return true, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]db.ListUsersRow, error) {
	return s.queries.ListUsers(ctx)
}
