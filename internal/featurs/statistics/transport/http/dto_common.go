package statistics_transport_http

import (
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
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

func bookDTOFromDomain(book core_domain.Book) BookDTOResponse {
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
		booksDTO[i] = bookDTOFromDomain(v)
	}

	return booksDTO
}

type StatisticsDTO struct {
	UsersCount         int               `json:"users_count"`
	BooksCount         int               `json:"books_count"`
	BooksOnShelves     int               `json:"books_on_shelves"`
	ReadBooksOnShelves int               `json:"reads_count"`
	TopBooksOnScore    []BookDTOResponse `json:"top_by_score"`
	TopBooksOnCount    []BookDTOResponse `json:"top_by_reads"`
}

func statisticsDTOFromDomain(stats core_domain.Statistics) StatisticsDTO {
	return StatisticsDTO{
		UsersCount:         stats.UsersCount,
		BooksCount:         stats.BooksCount,
		BooksOnShelves:     stats.BooksOnShelves,
		ReadBooksOnShelves: stats.ReadBooksOnShelves,
		TopBooksOnScore:    booksDTOFromDomains(stats.TopBooksOnScore),
		TopBooksOnCount:    booksDTOFromDomains(stats.TopBooksOnCount),
	}
}

type UserStatsDTO struct {
	BooksOnShelf            int      `json:"books_on_shelf"`
	BooksRead               int      `json:"books_read"`
	BooksReadRate           *float64 `json:"read_percent"`
	BooksAverageReadTime    *float64 `json:"avg_read_time_hours"`
	BooksReadFavoriteGenre  *string  `json:"favorite_genre"`
	BooksReadFavoriteAuthor *string  `json:"favorite_author"`
	UserReviews             int      `json:"reviews_count"`
}

func userStatsDTOFromDomain(u core_domain.UserStats) UserStatsDTO {
	return UserStatsDTO{
		BooksOnShelf:            u.BooksOnShelf,
		BooksRead:               u.BooksRead,
		BooksReadRate:           u.BooksReadRate,
		BooksAverageReadTime:    durationHours(u.BooksAverageReadTime),
		BooksReadFavoriteGenre:  u.BooksReadFavoriteGenre,
		BooksReadFavoriteAuthor: u.BooksReadFavoriteAuthor,
		UserReviews:             u.UserReviews,
	}
}

func durationHours(d *time.Duration) *float64 {
	if d == nil {
		return nil
	}
	h := d.Hours()
	return &h
}

type RatingCountDTO struct {
	Rating int `json:"rating"`
	Count  int `json:"count"`
}

func ratingCountDTOFromDomain(r core_domain.RatingCount) RatingCountDTO {
	return RatingCountDTO{
		Rating: r.Rating,
		Count:  r.Count,
	}
}

func ratingCountDTOsFromDomains(rates []core_domain.RatingCount) []RatingCountDTO {
	rateDTO := make([]RatingCountDTO, len(rates))
	for i, v := range rates {
		rateDTO[i] = ratingCountDTOFromDomain(v)
	}
	return rateDTO
}

type BookStatsDTO struct {
	UsersAddedOnShelf  int              `json:"users_on_shelf"`
	UsersRead          int              `json:"users_read"`
	BookAverageRate    *float64         `json:"average_rating"`
	RatingDistribution []RatingCountDTO `json:"rating_distribution"`
	ReviewCount        int              `json:"reviews_count"`
}

func bookStatsDTOFromDomain(b core_domain.BookStats) BookStatsDTO {
	return BookStatsDTO{
		UsersAddedOnShelf:  b.UsersAddedOnShelf,
		UsersRead:          b.UsersRead,
		BookAverageRate:    b.BookAverageRate,
		RatingDistribution: ratingCountDTOsFromDomains(b.RatingDistribution),
		ReviewCount:        b.ReviewCount,
	}
}
