package domain

import "errors"

var (
	ErrLocationNotFound        = errors.New("location not found")
	ErrDistributorNotFound     = errors.New("distributor not found")
	ErrInvalidParent           = errors.New("invalid parent distributor")
	ErrLocationAlreadyIncluded = errors.New("location already included")
	ErrLocationAlreadyExcluded = errors.New("location already excluded")
)
