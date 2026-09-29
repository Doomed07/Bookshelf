package statistics_repository_postgres

import (
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

type StatisticsModel struct {
	UsersCount         int
	BooksCount         int
	BooksOnShelves     int
	ReadBooksOnShelves int
}

func statisticsDomainFromModel(stat StatisticsModel) core_domain.Statistics {
	return core_domain.NewStatistics(
		stat.UsersCount,
		stat.BooksCount,
		stat.BooksOnShelves,
		stat.ReadBooksOnShelves,
		booksDomainsFromModels([]BookModel{}),
		booksDomainsFromModels([]BookModel{}),
	)
}

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

func bookDomainFromModel(book BookModel) core_domain.Book {
	return core_domain.NewBook(
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

func booksDomainsFromModels(books []BookModel) []core_domain.Book {
	booksDomain := make([]core_domain.Book, len(books))
	for i, book := range books {
		booksDomain[i] = bookDomainFromModel(book)
	}

	return booksDomain
}

type ShelfBookModel struct {
	UserID  int
	BookID  int
	Version int
	Read    bool
	Rating  *int
	Review  *string
	AddedAt time.Time
	ReadAt  *time.Time
}

func shelfDomainFromModel(s ShelfBookModel) core_domain.ShelfBook {
	return core_domain.NewShelfBook(
		s.UserID, s.BookID, s.Version, s.Read,
		s.Rating, s.Review, s.AddedAt, s.ReadAt,
	)
}

func shelvesDomainsFromModels(s []ShelfBookModel) []core_domain.ShelfBook {
	shelves := make([]core_domain.ShelfBook, len(s))
	for i, v := range s {
		shelves[i] = shelfDomainFromModel(v)
	}
	return shelves
}

type ShelfBookWithBookModel struct {
	ShelfBook ShelfBookModel
	Book      BookModel
}

func domainSBWBFromModel(m ShelfBookWithBookModel) core_domain.ShelfBookWithBook {
	return core_domain.NewShelfBookWithBook(
		m.ShelfBook.UserID,
		m.ShelfBook.BookID,
		m.ShelfBook.Version,
		m.ShelfBook.Read,
		m.ShelfBook.Rating,
		m.ShelfBook.Review,
		m.ShelfBook.AddedAt,
		m.ShelfBook.ReadAt,
		m.Book.ID,
		m.Book.Title,
		m.Book.Author,
		m.Book.Year,
		m.Book.Pages,
		m.Book.Genres,
		m.Book.Description,
		m.Book.Score,
		m.Book.ReadsCount,
	)
}

func booksFBSdomainFromModel(m []ShelfBookWithBookModel) []core_domain.ShelfBookWithBook {
	booksDomain := make([]core_domain.ShelfBookWithBook, len(m))
	for i, v := range m {
		booksDomain[i] = domainSBWBFromModel(v)
	}
	return booksDomain
}
