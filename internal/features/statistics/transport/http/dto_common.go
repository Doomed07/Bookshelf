package statistics_transport_http

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
	UsersCount         int               `json:"users_count" example:"128"`
	BooksCount         int               `json:"books_count" example:"1022"`
	BooksOnShelves     int               `json:"books_on_shelves" example:"3450"`
	ReadBooksOnShelves int               `json:"reads_count" example:"2130"`
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
	BooksOnShelf            int      `json:"books_on_shelf" example:"24"`
	BooksRead               int      `json:"books_read" example:"18"`
	BooksReadRate           *float64 `json:"read_percent" example:"75"`
	BooksAverageReadTime    *float64 `json:"avg_read_time_hours" example:"36.5"`
	BooksReadFavoriteGenre  *string  `json:"favorite_genre" example:"Фэнтези"`
	BooksReadFavoriteAuthor *string  `json:"favorite_author" example:"Дж. Р. Р. Толкин"`
	UserReviews             int      `json:"reviews_count" example:"9"`
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
	Rating int `json:"rating" example:"5"`
	Count  int `json:"count" example:"37"`
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
	UsersAddedOnShelf  int              `json:"users_on_shelf" example:"87"`
	UsersRead          int              `json:"users_read" example:"52"`
	BookAverageRate    *float64         `json:"average_rating" example:"4.6"`
	RatingDistribution []RatingCountDTO `json:"rating_distribution"`
	ReviewCount        int              `json:"reviews_count" example:"23"`
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
