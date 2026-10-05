package statistics_service

import (
	"context"
	"fmt"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (s *StatsService) GetBookStats(ctx context.Context, bookID int, from, to *time.Time) (core_domain.BookStats, error) {
	if err := core_domain.ValidateFromToQuery(from, to); err != nil {
		return core_domain.BookStats{}, fmt.Errorf("validate 'from/to' query param: %w", err)
	}

	exists, err := s.statsRepository.BookExists(ctx, bookID)
	if err != nil {
		return core_domain.BookStats{}, fmt.Errorf("check book exists: %w", err)
	}

	if !exists {
		return core_domain.BookStats{}, fmt.Errorf("book with ID: %d - %w",
			bookID, core_errors.ErrNotFound)
	}

	shelvesBook, err := s.statsRepository.GetBookShelves(ctx, bookID)
	if err != nil {
		return core_domain.BookStats{}, fmt.Errorf("get book statistics from repo: %w", err)
	}

	bookStats := calcBookStatistics(shelvesBook, from, to)
	return bookStats, nil
}

func calcBookStatistics(s []core_domain.ShelfBook, from, to *time.Time) core_domain.BookStats {
	users := 0
	usersRead := 0
	ratingSum := 0
	ratingCount := 0
	counts := [101]int{}
	reviewCount := 0

	for _, v := range s {
		if inPeriod(v.AddedAt, from, to) {
			users++
		}

		if !v.Read || v.ReadAt == nil || !inPeriod(*v.ReadAt, from, to) {
			continue
		}

		usersRead++

		if v.Rating != nil {
			ratingSum += *v.Rating
			ratingCount++
			counts[*v.Rating]++
		}

		if v.Review != nil {
			reviewCount++
		}
	}

	var bookRate *float64
	if ratingCount > 0 {
		rate := float64(ratingSum) / float64(ratingCount)
		bookRate = &rate
	}

	ratingCounts := make([]core_domain.RatingCount, 0)
	for rating, count := range counts {
		if count > 0 {
			ratingCounts = append(ratingCounts, core_domain.NewRatingCount(rating, count))
		}
	}

	return core_domain.NewBookStats(
		users, usersRead, reviewCount,
		bookRate,
		ratingCounts,
	)
}
