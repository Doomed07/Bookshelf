package statistics_service

import (
	"context"
	"fmt"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (s *StatsService) GetSystemStats(ctx context.Context, from, to *time.Time, top *int) (core_domain.Statistics, error) {
	if err := core_domain.ValidateFromToQuery(from, to); err != nil {
		return core_domain.Statistics{}, fmt.Errorf("validate 'from/to' query param: %w", err)
	}

	t, err := core_domain.NormalizeTopQueryParam(top)
	if err != nil {
		return core_domain.Statistics{}, fmt.Errorf("validate 'top' query param: %w", err)
	}

	stats, err := s.statsRepository.GetSystemCounters(ctx, from, to)
	if err != nil {
		return core_domain.Statistics{}, fmt.Errorf("get system stats from repo: %w", err)
	}

	topScoreBooks, err := s.statsRepository.GetTopBooksByScore(ctx, s.minRatings, t)
	if err != nil {
		return core_domain.Statistics{}, fmt.Errorf("get top score books from repo: %w", err)
	}

	topReadBooks, err := s.statsRepository.GetTopBooksByReads(ctx, t)
	if err != nil {
		return core_domain.Statistics{}, fmt.Errorf("get top reads books from repo: %w", err)
	}

	stats.TopBooksOnScore = topScoreBooks
	stats.TopBooksOnCount = topReadBooks

	return stats, nil
}
