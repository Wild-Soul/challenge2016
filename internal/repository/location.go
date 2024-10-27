package repository

import "distributor/internal/domain"

type LocationRepository interface {
	Store(location *domain.Location, key string) error
	Find(key string) (*domain.Location, error)
	FindAll() (map[string]*domain.Location, error)
}
