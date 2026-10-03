package bookshelf_transport_http

import (
	"time"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	books_transport_http "github.com/Doomed07/Bookshelf/internal/featurs/books/transport/http"
)

type ShelfBookDTOResponse struct {
	UserID  int        `json:"user_id" example:"1"`
	BookID  int        `json:"book_id" example:"1"`
	Version int        `json:"version" example:"3"`
	Read    bool       `json:"read" example:"true"`
	Rating  *int       `json:"rating" example:"5"`
	Review  *string    `json:"review" example:"Одна из лучших книг, что я читал."`
	AddedAt time.Time  `json:"added_at" example:"2026-01-10T09:00:00Z"`
	ReadAt  *time.Time `json:"read_at" example:"2026-01-15T12:00:00Z"`
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
	Name   string    `json:"name" example:"finished"`
	BookID int       `json:"book_id" example:"1"`
	Title  string    `json:"title" example:"Хоббит, или Туда и обратно"`
	Author string    `json:"author" example:"Дж. Р. Р. Толкин"`
	Rating *int      `json:"rating" example:"5"`
	Review *string   `json:"review" example:"Одна из лучших книг, что я читал."`
	At     time.Time `json:"at" example:"2026-01-15T12:00:00Z"`
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
