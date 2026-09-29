package books_transport_http

import (
	"time"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
)

type BookDTOResponse struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Year        int      `json:"year"`
	Pages       int      `json:"pages"`
	Genres      []string `json:"genres"`
	Description string   `json:"description"`
	Score       *int     `json:"score"`
	ReadsCount  int      `json:"reads_count"`
}

func BookDTOFromDomain(book core_domain.Book) BookDTOResponse {
	return BookDTOResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Year:        book.Year,
		Pages:       book.Pages,
		Genres:      book.Genres,
		Description: book.Description,
		Score:       book.Score,
		ReadsCount:  book.ReadsCount,
	}
}

func booksDTOFromDomains(books []core_domain.Book) []BookDTOResponse {
	booksDTO := make([]BookDTOResponse, len(books))

	for i, v := range books {
		booksDTO[i] = BookDTOFromDomain(v)
	}

	return booksDTO
}

type ReviewDTOResponse struct {
	UserID   int       `json:"user_id"`
	Username string    `json:"username"`
	Rating   *int      `json:"rating"`
	Review   *string   `json:"review"`
	ReadAt   time.Time `json:"read_at"`
}

func reviewDTOFromDomain(review core_domain.Review) ReviewDTOResponse {
	return ReviewDTOResponse{
		UserID:   review.UserID,
		Username: review.Username,
		Rating:   review.Rating,
		Review:   review.Review,
		ReadAt:   review.ReadAt,
	}
}

func reviewsDTOFromDomains(reviews []core_domain.Review) []ReviewDTOResponse {
	reviewsDTO := make([]ReviewDTOResponse, len(reviews))

	for i, r := range reviews {
		reviewsDTO[i] = reviewDTOFromDomain(r)
	}

	return reviewsDTO
}
