package service

import (
	"context"
	"gin-quickstart/repository"
)

type PropertyService struct {
	propertyRepo *repository.PropertyRepo
}

func NewPropertyService(propertyRepo *repository.PropertyRepo) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
	}
}

func (s *PropertyService) AddProperty(ctx context.Context, userID int, name string, photoID int) (repository.PropertyData, error) {
	data, err := s.propertyRepo.CreateProperty(ctx, userID, name, photoID)
	if err != nil {
		return repository.PropertyData{}, err
	}
	return data, nil
}
