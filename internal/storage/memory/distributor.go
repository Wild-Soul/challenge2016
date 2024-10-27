package memory

import (
	"sync"

	"distributor/internal/domain"
	"distributor/internal/repository"
)

type distributorRepository struct {
	distributors map[string]*domain.Distributor
	mu           sync.RWMutex
}

func NewDistributorRepository() repository.DistributorRepository {
	return &distributorRepository{
		distributors: make(map[string]*domain.Distributor),
	}
}

func (r *distributorRepository) Store(distributor *domain.Distributor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.distributors[distributor.Name] = distributor
	return nil
}

func (r *distributorRepository) Find(name string) (*domain.Distributor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	distributor, exists := r.distributors[name]
	if !exists {
		return nil, domain.ErrDistributorNotFound
	}
	return distributor, nil
}

func (r *distributorRepository) FindAll() (map[string]*domain.Distributor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	distributors := make(map[string]*domain.Distributor, len(r.distributors))
	for k, v := range r.distributors {
		distributors[k] = v
	}
	return distributors, nil
}
