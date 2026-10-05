package statistics_service

import (
	"context"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

type StatsService struct {
	statsRepository StatsRepository
	minRatings      int
}

type StatsRepository interface {
	GetSystemCounters(ctx context.Context, from, to *time.Time) (core_domain.Statistics, error)
	GetTopBooksByScore(ctx context.Context, minRatings, top int) ([]core_domain.Book, error)
	GetTopBooksByReads(ctx context.Context, top int) ([]core_domain.Book, error)
	GetUserShelf(ctx context.Context, userID int) ([]core_domain.ShelfBookWithBook, error)
	UserExists(ctx context.Context, userID int) (bool, error)
	GetBookShelves(ctx context.Context, bookID int) ([]core_domain.ShelfBook, error)
	BookExists(ctx context.Context, bookID int) (bool, error)
}

func NewStatsService(statsRepository StatsRepository, config Config) *StatsService {
	return &StatsService{
		statsRepository: statsRepository,
		minRatings:      config.MinRatings,
	}
}
