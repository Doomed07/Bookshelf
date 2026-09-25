package domain

import (
	"fmt"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

func NormalizePagination(limit, offset *int) (int, int, error) {
	lim, off := DefaultLimit, 0

	if limit != nil {
		if *limit < 1 || *limit > MaxLimit {
			return 0, 0, fmt.Errorf("invalid 'limit' query param: %w", core_errors.ErrInvalidArgument)
		}
		lim = *limit
	}

	if offset != nil {
		if *offset < 0 {
			return 0, 0, fmt.Errorf("invalid 'offset' query param: %w", core_errors.ErrInvalidArgument)
		}
		off = *offset
	}

	return lim, off, nil
}
