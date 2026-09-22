package service

import (
	"context"
	db "gin-quickstart/db/sqlc"
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
	data, err := s.queries.CreateProperty(ctx, db.CreatePropertyParams{
		UserID:  int64(userID),
		Name:    name,
		PhotoID: int32(photoID),
	})
	if err != nil {
		return db.Property{}, err
	}
	return data, nil
}
