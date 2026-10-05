package books_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

// GetRecentReviews — последние рецензии всех читателей. Отдельного времени
// написания рецензии в базе нет, поэтому сортируем по дате прочтения.
func (r *BooksRepository) GetRecentReviews(ctx context.Context, limit, offset int) ([]core_domain.ReviewWithBook, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT s.user_id, u.username, s.rating, s.review, s.read_at, s.reviewed_at,
	       b.id, b.title, b.author, b.year, b.pages, b.genres, b.description, b.score, b.reads_count
	FROM bookshelfapp.bookshelf s
	JOIN bookshelfapp.users u ON u.id = s.user_id
	JOIN bookshelfapp.books b ON b.id = s.book_id
	WHERE s.review IS NOT NULL
	ORDER BY s.reviewed_at DESC, s.user_id, s.book_id
	LIMIT $1 OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var reviews []ReviewWithBookModel
	for rows.Next() {
		var m ReviewWithBookModel
		err := rows.Scan(
			&m.Review.UserID,
			&m.Review.Username,
			&m.Review.Rating,
			&m.Review.Review,
			&m.Review.ReadAt,
			&m.Review.ReviewedAt,
			&m.Book.ID,
			&m.Book.Title,
			&m.Book.Author,
			&m.Book.Year,
			&m.Book.Pages,
			&m.Book.Genres,
			&m.Book.Description,
			&m.Book.Score,
			&m.Book.ReadsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		reviews = append(reviews, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return reviewsWithBookDomainsFromModels(reviews), nil
}
