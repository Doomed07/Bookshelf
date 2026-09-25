package books_repository_postgres

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *BooksRepository) GetReviews(ctx context.Context, id, limit, offset int) ([]domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT b.user_id, u.username, b.rating, b.review, b.read_at
	FROM bookshelfapp.bookshelf b
	JOIN bookshelfapp.users u ON u.id = b.user_id
	WHERE b.book_id = $1
	AND (b.rating IS NOT NULL OR b.review IS NOT NULL)
	ORDER By b.read_at DESC
	LIMIT $2
	OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var reviews []ReviewModel
	for rows.Next() {
		var rev ReviewModel
		err := rows.Scan(
			&rev.UserID,
			&rev.Username,
			&rev.Rating,
			&rev.Review,
			&rev.ReadAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		reviews = append(reviews, rev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	reviewsDomain := reviewsDomainsFromModels(reviews)

	return reviewsDomain, nil
}
