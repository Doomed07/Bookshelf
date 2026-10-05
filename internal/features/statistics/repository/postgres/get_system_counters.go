package statistics_repository_postgres

import (
	"context"
	"fmt"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *StatsRepository) GetSystemCounters(ctx context.Context, from, to *time.Time) (core_domain.Statistics, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
    (SELECT count(*) FROM bookshelfapp.users),
    (SELECT count(*) FROM bookshelfapp.books),
    count(*) FILTER (
        WHERE ($1::timestamptz IS NULL OR s.added_at >= $1)
          AND ($2::timestamptz IS NULL OR s.added_at <  $2 + interval '1 day')
    ),
    count(*) FILTER (
        WHERE s.read
          AND ($1::timestamptz IS NULL OR s.read_at >= $1)
          AND ($2::timestamptz IS NULL OR s.read_at <  $2 + interval '1 day')
    )
FROM bookshelfapp.bookshelf s;
	`

	row := r.pool.QueryRow(ctx, query, from, to)

	var s StatisticsModel
	err := row.Scan(
		&s.UsersCount,
		&s.BooksCount,
		&s.BooksOnShelves,
		&s.ReadBooksOnShelves,
	)
	if err != nil {
		return core_domain.Statistics{}, fmt.Errorf("scan row: %w", err)
	}

	statsDomain := statisticsDomainFromModel(s)
	return statsDomain, nil
}
