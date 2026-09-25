package domain

import "time"

type Review struct {
	UserID   int
	Username string
	Rating   *int
	Review   *string
	ReadAt   time.Time
}

func NewReview(
	id int,
	username string,
	rating *int,
	review *string,
	readAt time.Time,
) Review {
	return Review{
		UserID:   id,
		Username: username,
		Rating:   rating,
		Review:   review,
		ReadAt:   readAt,
	}
}
