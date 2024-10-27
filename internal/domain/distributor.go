package domain

import "sync"

type Distributor struct {
	Name     string
	Parent   string
	Includes map[string]bool
	Excludes map[string]bool
	mu       sync.RWMutex
}

func NewDistributor(name, parent string) *Distributor {
	return &Distributor{
		Name:     name,
		Parent:   parent,
		Includes: make(map[string]bool),
		Excludes: make(map[string]bool),
	}
}

func (d *Distributor) UpdatePermission(location string, isInclude bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if isInclude {
		if d.Excludes[location] {
			return ErrLocationAlreadyExcluded
		}
		delete(d.Excludes, location)
		d.Includes[location] = true
	} else {
		if d.Includes[location] {
			return ErrLocationAlreadyIncluded
		}
		delete(d.Includes, location)
		d.Excludes[location] = true
	}

	return nil
}
