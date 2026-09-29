package core_domain

import "time"

type Statistics struct {
	UsersCount         int
	BooksCount         int
	BooksOnShelves     int
	ReadBooksOnShelves int
	TopBooksOnScore    []Book
	TopBooksOnCount    []Book
}

func NewStatistics(
	usersCount, booksCount, booksOnShelves, readBookOnShelves int,
	topScore, topCount []Book,
) Statistics {
	return Statistics{
		UsersCount:         usersCount,
		BooksCount:         booksCount,
		BooksOnShelves:     booksOnShelves,
		ReadBooksOnShelves: readBookOnShelves,
		TopBooksOnScore:    topScore,
		TopBooksOnCount:    topCount,
	}
}

type UserStats struct {
	BooksOnShelf            int
	BooksRead               int
	BooksReadRate           *float64
	BooksAverageReadTime    *time.Duration
	BooksReadFavoriteGenre  *string
	BooksReadFavoriteAuthor *string
	UserReviews             int
}

func NewUserStats(
	booksOnShelf, booksRead, userReview int,
	booksReadRate *float64,
	booksART *time.Duration,
	booksFG, booksFA *string,
) UserStats {
	return UserStats{
		BooksOnShelf:            booksOnShelf,
		BooksRead:               booksRead,
		BooksReadRate:           booksReadRate,
		BooksAverageReadTime:    booksART,
		BooksReadFavoriteGenre:  booksFG,
		BooksReadFavoriteAuthor: booksFA,
		UserReviews:             userReview,
	}
}

type RatingCount struct {
	Rating int
	Count  int
}

func NewRatingCount(r, c int) RatingCount {
	return RatingCount{
		Rating: r,
		Count:  c,
	}
}

type BookStats struct {
	UsersAddedOnShelf  int
	UsersRead          int
	BookAverageRate    *float64
	RatingDistribution []RatingCount
	ReviewCount        int
}

func NewBookStats(
	usersAOS, usersRead, reviewCount int,
	bookAR *float64,
	ratingDistribution []RatingCount,
) BookStats {
	return BookStats{
		UsersAddedOnShelf:  usersAOS,
		UsersRead:          usersRead,
		BookAverageRate:    bookAR,
		RatingDistribution: ratingDistribution,
		ReviewCount:        reviewCount,
	}
}
