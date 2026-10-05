package core_domain

import "time"

type Review struct {
	UserID     int
	Username   string
	Rating     *int
	Review     *string
	ReadAt     time.Time
	ReviewedAt *time.Time
}

func NewReview(
	id int,
	username string,
	rating *int,
	review *string,
	readAt time.Time,
	reviewedAt *time.Time,
) Review {
	return Review{
		UserID:     id,
		Username:   username,
		Rating:     rating,
		Review:     review,
		ReadAt:     readAt,
		ReviewedAt: reviewedAt,
	}
}

type ReviewWithBook struct {
	Review Review
	Book   Book
}

func NewReviewWithBook(review Review, book Book) ReviewWithBook {
	return ReviewWithBook{
		Review: review,
		Book:   book,
	}
}
