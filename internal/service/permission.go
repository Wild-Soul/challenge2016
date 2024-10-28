package service

import (
	"distributor/internal/domain"
	"distributor/internal/repository"
	"strings"
)

type PermissionService struct {
	locationRepo    repository.LocationRepository
	distributorRepo repository.DistributorRepository
}

func NewPermissionService(
	locationRepo repository.LocationRepository,
	distributorRepo repository.DistributorRepository,
) *PermissionService {
	return &PermissionService{
		locationRepo:    locationRepo,
		distributorRepo: distributorRepo,
	}
}

func (s *PermissionService) CheckPermission(distributorName, location string) (bool, error) {
	// Validate location exists
	if _, err := s.locationRepo.Find(location); err != nil {
		return false, err
	}

	// Get distributor
	dist, err := s.distributorRepo.Find(distributorName)
	if err != nil {
		return false, err
	}

	return s.hasPermission(dist, location)
}

func (s *PermissionService) hasPermission(dist *domain.Distributor, location string) (bool, error) {
	locationParts := strings.Split(location, "-")

	// Check excludes first
	for excluded := range dist.Excludes {
		excludedParts := strings.Split(excluded, "-")
		if isLocationMatch(locationParts, excludedParts) {
			return false, nil // Denied if the location matches any exclusion
		}
	}

	// Check parent, as their exlude is also applied to children.
	if dist.Parent != "" {
		parent, err := s.distributorRepo.Find(dist.Parent)
		if err != nil {
			return false, err
		}
		return s.hasPermission(parent, location)
	}

	// Finally check includes of the distributor.
	for included := range dist.Includes {
		includedParts := strings.Split(included, "-")
		if isLocationMatch(locationParts, includedParts) {
			return true, nil // Granted if the location matches any inclusion
		}
	}

	// If no permissions match and no parent exists
	return false, nil
}

func isLocationMatch(location, pattern []string) bool {
	// Allow matching based on the hierarchical structure
	if len(pattern) > len(location) {
		return false
	}

	for i := range pattern {
		if location[len(location)-1-i] != pattern[len(pattern)-1-i] {
			return false
		}
	}

	// If pattern is shorter or equal, we can also allow broader matches
	if len(pattern) < len(location) {
		return true
	}

	return false
}
