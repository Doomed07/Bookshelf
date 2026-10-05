package books_transport_http

import (
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

type BookDTOResponse struct {
	ID          int      `json:"id" example:"1"`
	Title       string   `json:"title" example:"Хоббит, или Туда и обратно"`
	Author      string   `json:"author" example:"Дж. Р. Р. Толкин"`
	Year        int      `json:"year" example:"1937"`
	Pages       int      `json:"pages" example:"310"`
	Genres      []string `json:"genres" example:"Фэнтези,Приключения"`
	Description string   `json:"description" example:"Повесть о путешествии хоббита Бильбо Бэггинса, втянутого в опасное приключение с гномами и драконом."`
	Score       *int     `json:"score" example:"5"`
	ReadsCount  int      `json:"reads_count" example:"42"`
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
	UserID     int        `json:"user_id" example:"1"`
	Username   string     `json:"username" example:"book_worm07"`
	Rating     *int       `json:"rating" example:"5"`
	Review     *string    `json:"review" example:"Одна из лучших книг, что я читал."`
	ReadAt     time.Time  `json:"read_at" example:"2026-01-15T12:00:00Z"`
	ReviewedAt *time.Time `json:"reviewed_at" example:"2026-01-20T18:30:00Z"`
}

func reviewDTOFromDomain(review core_domain.Review) ReviewDTOResponse {
	return ReviewDTOResponse{
		UserID:     review.UserID,
		Username:   review.Username,
		Rating:     review.Rating,
		Review:     review.Review,
		ReadAt:     review.ReadAt,
		ReviewedAt: review.ReviewedAt,
	}
}

func reviewsDTOFromDomains(reviews []core_domain.Review) []ReviewDTOResponse {
	reviewsDTO := make([]ReviewDTOResponse, len(reviews))

	for i, r := range reviews {
		reviewsDTO[i] = reviewDTOFromDomain(r)
	}

	return reviewsDTO
}

type RecentReviewDTOResponse struct {
	Review ReviewDTOResponse `json:"review"`
	Book   BookDTOResponse   `json:"book"`
}

func recentReviewsDTOFromDomains(reviews []core_domain.ReviewWithBook) []RecentReviewDTOResponse {
	reviewsDTO := make([]RecentReviewDTOResponse, len(reviews))
	for i, r := range reviews {
		reviewsDTO[i] = RecentReviewDTOResponse{
			Review: reviewDTOFromDomain(r.Review),
			Book:   BookDTOFromDomain(r.Book),
		}
	}

	return reviewsDTO
}
