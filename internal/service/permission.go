package service

import (
	"strings"

	"distributor/internal/domain"
	"distributor/internal/repository"
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
	// Check parent permissions first
	if dist.Parent != "" {
		parent, err := s.distributorRepo.Find(dist.Parent)
		if err != nil {
			return false, err
		}
		hasParentPerm, err := s.hasPermission(parent, location)
		if err != nil || !hasParentPerm {
			return false, err
		}
	}

	locationParts := strings.Split(location, "-")

	// Check excludes first
	for excluded := range dist.Excludes {
		excludedParts := strings.Split(excluded, "-")
		if isLocationMatch(locationParts, excludedParts) {
			return false, nil
		}
	}

	// Check includes
	for included := range dist.Includes {
		includedParts := strings.Split(included, "-")
		if isLocationMatch(locationParts, includedParts) {
			return true, nil
		}
	}

	return false, nil
}

func isLocationMatch(location, pattern []string) bool {
	if len(pattern) > len(location) {
		return false
	}

	for i := range pattern {
		if location[len(location)-1-i] != pattern[len(pattern)-1-i] {
			return false
		}
	}

	return true
}
