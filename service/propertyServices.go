package service

import (
	"context"
	db "gin-quickstart/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type PropertyService struct {
	queries *db.Queries
}

func NewPropertyService(queries *db.Queries) *PropertyService {
	return &PropertyService{
		queries: queries,
	}
}

func (s *PropertyService) AddProperty(ctx context.Context, userID int, name string, photoID int) (db.Property, error) {
	var photoIDParam pgtype.Int8
	if photoID != 0 {
		photoIDParam = pgtype.Int8{Int64: int64(photoID), Valid: true}
	}

	data, err := s.queries.CreateProperty(ctx, db.CreatePropertyParams{
		UserID:  int64(userID),
		Name:    name,
		PhotoID: photoIDParam,
	})
	if err != nil {
		return db.Property{}, err
	}
	return data, nil
}

func (s *PropertyService) DeleteProperty(ctx context.Context, userID int) error {
	return s.queries.DeleteProperty(ctx, int64(userID))
}
