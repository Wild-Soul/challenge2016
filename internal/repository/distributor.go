package repository

import "distributor/internal/domain"

type DistributorRepository interface {
	Store(distributor *domain.Distributor) error
	Find(name string) (*domain.Distributor, error)
	FindAll() (map[string]*domain.Distributor, error)
}
