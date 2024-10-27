package memory

import (
	"fmt"
	"sync"

	"distributor/internal/domain"
	"distributor/internal/repository"
)

type locationRepository struct {
	locations map[string]*domain.Location
	mu        sync.RWMutex
}

func NewLocationRepository() repository.LocationRepository {
	return &locationRepository{
		locations: make(map[string]*domain.Location),
	}
}

func (r *locationRepository) Store(location *domain.Location, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.locations[key] = location
	return nil
}

func (r *locationRepository) Find(key string) (*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// for key := range r.locations {
	// 	fmt.Println("KEY::", key)
	// }
	location, exists := r.locations[key]
	if !exists {
		fmt.Println("LOCATION KEY::", key)
		return nil, domain.ErrLocationNotFound
	}
	return location, nil
}

func (r *locationRepository) FindAll() (map[string]*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent concurrent map access
	locations := make(map[string]*domain.Location, len(r.locations))
	for k, v := range r.locations {
		locations[k] = v
	}
	return locations, nil
}
