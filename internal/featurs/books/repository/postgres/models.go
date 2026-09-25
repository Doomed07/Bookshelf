package books_repository_postgres

import (
	"time"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

type BookModel struct {
	ID          int
	Title       string
	Author      string
	Year        int
	Pages       int
	Genres      []string
	Description string
	Score       *int
	ReadsCount  int
}

func bookDomainFromModel(book BookModel) domain.Book {
	return domain.NewBook(
		book.ID,
		book.Title,
		book.Author,
		book.Year,
		book.Pages,
		book.Genres,
		book.Description,
		book.Score,
		book.ReadsCount,
	)
}

func booksDomainsFromModels(books []BookModel) []domain.Book {
	booksDomain := make([]domain.Book, len(books))
	for i, book := range books {
		booksDomain[i] = bookDomainFromModel(book)
	}

	return booksDomain
}

type ReviewModel struct {
	UserID   int
	Username string
	Rating   *int
	Review   *string
	ReadAt   time.Time
}

func reviewDomainFromModel(review ReviewModel) domain.Review {
	return domain.NewReview(
		review.UserID,
		review.Username,
		review.Rating,
		review.Review,
		review.ReadAt,
	)
}

func reviewsDomainsFromModels(reviews []ReviewModel) []domain.Review {
	reviewsDomain := make([]domain.Review, len(reviews))
	for i, r := range reviews {
		reviewsDomain[i] = reviewDomainFromModel(r)
	}

	return reviewsDomain
}
