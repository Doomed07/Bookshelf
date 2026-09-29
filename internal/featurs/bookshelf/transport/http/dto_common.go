package bookshelf_transport_http

import (
	"time"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	books_transport_http "github.com/Doomed07/Bookshelf/internal/featurs/books/transport/http"
)

type ShelfBookDTOResponse struct {
	UserID  int        `json:"user_id"`
	BookID  int        `json:"book_id"`
	Version int        `json:"version"`
	Read    bool       `json:"read"`
	Rating  *int       `json:"rating"`
	Review  *string    `json:"review"`
	AddedAt time.Time  `json:"added_at"`
	ReadAt  *time.Time `json:"read_at"`
}

func shelfBookDTOFromDomain(shelf core_domain.ShelfBook) ShelfBookDTOResponse {
	return ShelfBookDTOResponse{
		UserID:  shelf.UserID,
		BookID:  shelf.BookID,
		Version: shelf.Version,
		Read:    shelf.Read,
		Rating:  shelf.Rating,
		Review:  shelf.Review,
		AddedAt: shelf.AddedAt,
		ReadAt:  shelf.ReadAt,
	}
}

type ShelfBookWithBookDTOResponse struct {
	ShelfBook ShelfBookDTOResponse                 `json:"shelf"`
	Book      books_transport_http.BookDTOResponse `json:"book"`
}

func shelfBookWithBookDTOFromDomain(s core_domain.ShelfBookWithBook) ShelfBookWithBookDTOResponse {
	return ShelfBookWithBookDTOResponse{
		ShelfBook: shelfBookDTOFromDomain(s.ShelfBook),
		Book:      books_transport_http.BookDTOFromDomain(s.Book),
	}
}

func booksFBSDTOFromDomain(b []core_domain.ShelfBookWithBook) []ShelfBookWithBookDTOResponse {
	booksDTO := make([]ShelfBookWithBookDTOResponse, len(b))
	for i, v := range b {
		booksDTO[i] = shelfBookWithBookDTOFromDomain(v)
	}
	return booksDTO
}

type EventDTO struct {
	Name   string    `json:"name"`
	BookID int       `json:"book_id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Rating *int      `json:"rating"`
	Review *string   `json:"review"`
	At     time.Time `json:"at"`
}

func eventDTOFromDomain(d core_domain.Event) EventDTO {
	return EventDTO{
		Name:   d.Name,
		BookID: d.BookID,
		Title:  d.Title,
		Author: d.Author,
		Rating: d.Rating,
		Review: d.Review,
		At:     d.At,
	}
}

func eventsDTOFromDomain(d []core_domain.Event) []EventDTO {
	events := make([]EventDTO, len(d))
	for i, v := range d {
		events[i] = eventDTOFromDomain(v)
	}
	return events
}
