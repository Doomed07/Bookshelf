package bookshelf_repository_postgres

import (
	"time"

	"github.com/Doomed07/Bookshelf/internal/core/core_domain"
)

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

type EventModel struct {
	Name   string
	BookID int
	Title  string
	Author string
	Rating *int
	Review *string
	At     time.Time
}

func eventDomainFromModel(m EventModel) core_domain.Event {
	return core_domain.NewEvent(
		m.Name, m.Title, m.Author, m.BookID, m.Rating, m.Review, m.At,
	)
}

func eventsDomainFromModel(m []EventModel) []core_domain.Event {
	events := make([]core_domain.Event, len(m))
	for i, v := range m {
		events[i] = eventDomainFromModel(v)
	}
	return events
}
