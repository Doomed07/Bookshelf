package statistics_service

import (
	"context"
	"fmt"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func (s *StatsService) GetUserStats(ctx context.Context, userID int, from, to *time.Time) (core_domain.UserStats, error) {
	if err := core_domain.ValidateFromToQuery(from, to); err != nil {
		return core_domain.UserStats{}, fmt.Errorf("validate 'from/to' query param: %w", err)
	}

	exists, err := s.statsRepository.UserExists(ctx, userID)
	if err != nil {
		return core_domain.UserStats{}, fmt.Errorf("check user exists: %w", err)
	}

	if !exists {
		return core_domain.UserStats{}, fmt.Errorf("user %d: %w", userID, core_errors.ErrNotFound)
	}

	userBooks, err := s.statsRepository.GetUserShelf(ctx, userID)
	if err != nil {
		return core_domain.UserStats{}, fmt.Errorf("get user statistics from repo: %w", err)
	}

	userStatsDomain := calcUserStatistics(userBooks, from, to)
	return userStatsDomain, nil
}

func calcUserStatistics(books []core_domain.ShelfBookWithBook, from, to *time.Time) core_domain.UserStats {
	booksOnShelf := 0
	booksRead := 0
	addedAndRead := 0
	userReviews := 0
	var totalReadDuration time.Duration
	genres := make(map[string]int)
	authors := make(map[string]int)

	for _, v := range books {
		if inPeriod(v.ShelfBook.AddedAt, from, to) {
			booksOnShelf++
			if v.ShelfBook.Read && v.ShelfBook.ReadAt != nil && inPeriod(*v.ShelfBook.ReadAt, from, to) {
				addedAndRead++
			}
		}

		if !v.ShelfBook.Read || v.ShelfBook.ReadAt == nil || !inPeriod(*v.ShelfBook.ReadAt, from, to) {
			continue
		}

		booksRead++

		authors[v.Book.Author]++

		for _, g := range v.Book.Genres {
			genres[g]++
		}

		readDuration := v.ShelfBook.ReadDuration()
		if readDuration != nil {
			totalReadDuration += *readDuration
		}

		if v.ShelfBook.Review != nil {
			userReviews++
		}
	}

	var booksReadRate *float64
	if booksOnShelf > 0 {
		rate := float64(addedAndRead) / float64(booksOnShelf) * 100
		booksReadRate = &rate
	}

	var booksAverageReadTime *time.Duration
	if booksRead > 0 {
		avg := totalReadDuration / time.Duration(booksRead)
		booksAverageReadTime = &avg
	}

	return core_domain.UserStats{
		BooksOnShelf:            booksOnShelf,
		BooksRead:               booksRead,
		BooksReadRate:           booksReadRate,
		BooksAverageReadTime:    booksAverageReadTime,
		BooksReadFavoriteGenre:  favorite(genres),
		BooksReadFavoriteAuthor: favorite(authors),
		UserReviews:             userReviews,
	}
}

func favorite(counts map[string]int) *string {
	var best string
	bestCount := 0

	for name, count := range counts {
		if count > bestCount || (count == bestCount && name < best) {
			best, bestCount = name, count
		}
	}

	if bestCount == 0 {
		return nil
	}
	return &best
}
