package statistics_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

func (r *StatsRepository) GetUserShelf(ctx context.Context, userID int) ([]core_domain.ShelfBookWithBook, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT s.user_id, s.book_id, s.version, s.read, s.rating, s.review, s.added_at, s.read_at, 
	b.id, b.title, b.author, b.year, b.pages, b.genres, b.description, b.score, b.reads_count
	FROM bookshelfapp.bookshelf s
	JOIN bookshelfapp.books b ON b.id = s.book_id
	WHERE s.user_id = $1
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get query rows: %w", err)
	}
	defer rows.Close()

	var booksModel []ShelfBookWithBookModel
	for rows.Next() {
		var book ShelfBookWithBookModel
		err := rows.Scan(
			&book.ShelfBook.UserID,
			&book.ShelfBook.BookID,
			&book.ShelfBook.Version,
			&book.ShelfBook.Read,
			&book.ShelfBook.Rating,
			&book.ShelfBook.Review,
			&book.ShelfBook.AddedAt,
			&book.ShelfBook.ReadAt,
			&book.Book.ID,
			&book.Book.Title,
			&book.Book.Author,
			&book.Book.Year,
			&book.Book.Pages,
			&book.Book.Genres,
			&book.Book.Description,
			&book.Book.Score,
			&book.Book.ReadsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
		booksModel = append(booksModel, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows next: %w", err)
	}

	booksDomain := booksFBSdomainFromModel(booksModel)
	return booksDomain, nil
}
