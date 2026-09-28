package bookshelf_repository_postgres

import (
	"context"
	"fmt"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
)

func (r *BookshelfRepository) GetUserActivity(ctx context.Context, userID int, limit, offset int) ([]core_domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
		CASE WHEN s.read THEN 'finished' ELSE 'added' END AS name,
		s.book_id, b.title, b.author, s.rating, s.review,
		CASE WHEN s.read THEN s.read_at ELSE s.added_at END AS at
	FROM bookshelfapp.bookshelf s
	JOIN bookshelfapp.books b ON b.id = s.book_id
	WHERE s.user_id = $1
	ORDER BY at DESC, s.book_id
	LIMIT $2 OFFSET $3;
`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var events []EventModel
	for rows.Next() {
		var e EventModel
		err := rows.Scan(
			&e.Name,
			&e.BookID,
			&e.Title,
			&e.Author,
			&e.Rating,
			&e.Review,
			&e.At,
		)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows next: %w", err)
	}

	eventsDomain := eventsDomainFromModel(events)
	return eventsDomain, nil
}
