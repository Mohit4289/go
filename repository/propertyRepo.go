package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type PropertyData struct {
	UserID  int
	Name    string
	PhotoID int
}

type PropertyRepo struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func NewPropertyRepo(db *pgxpool.Pool, redis *redis.Client) *PropertyRepo {
	return &PropertyRepo{
		DB:    db,
		Redis: redis,
	}
}

func (r *PropertyRepo) CreateProperty(ctx context.Context, userID int, propName string, photoID int) (PropertyData, error) {
	var property PropertyData

	err := r.DB.QueryRow(ctx, `INSERT INTO property (userID, name, photoID) VALUES ($1, $2, $3) RETURNING userID, name, photoID`, userID, propName, photoID).Scan(&property.UserID, &property.Name, &property.PhotoID)
	if err != nil {
		return PropertyData{}, err
	}

	return property, nil
}
