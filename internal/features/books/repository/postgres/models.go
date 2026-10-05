package books_repository_postgres

import (
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
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

func bookDomainFromModel(bm BookModel) core_domain.Book {
	return core_domain.NewBook(
		bm.ID,
		bm.Title,
		bm.Author,
		bm.Year,
		bm.Pages,
		bm.Genres,
		bm.Description,
		bm.Score,
		bm.ReadsCount,
	)
}

func booksDomainsFromModels(books []BookModel) []core_domain.Book {
	booksDomain := make([]core_domain.Book, len(books))
	for i, book := range books {
		booksDomain[i] = bookDomainFromModel(book)
	}

	return booksDomain
}

type ReviewModel struct {
	UserID     int
	Username   string
	Rating     *int
	Review     *string
	ReadAt     time.Time
	ReviewedAt *time.Time
}

func reviewDomainFromModel(rm ReviewModel) core_domain.Review {
	return core_domain.NewReview(
		rm.UserID,
		rm.Username,
		rm.Rating,
		rm.Review,
		rm.ReadAt,
		rm.ReviewedAt,
	)
}

func reviewsDomainsFromModels(reviews []ReviewModel) []core_domain.Review {
	reviewsDomain := make([]core_domain.Review, len(reviews))
	for i, r := range reviews {
		reviewsDomain[i] = reviewDomainFromModel(r)
	}

	return reviewsDomain
}

type ReviewWithBookModel struct {
	Review ReviewModel
	Book   BookModel
}

func reviewWithBookDomainFromModel(model ReviewWithBookModel) core_domain.ReviewWithBook {
	return core_domain.NewReviewWithBook(
		reviewDomainFromModel(model.Review),
		bookDomainFromModel(model.Book),
	)
}

func reviewsWithBookDomainsFromModels(models []ReviewWithBookModel) []core_domain.ReviewWithBook {
	reviewsDomain := make([]core_domain.ReviewWithBook, len(models))
	for i, m := range models {
		reviewsDomain[i] = reviewWithBookDomainFromModel(m)
	}

	return reviewsDomain
}
