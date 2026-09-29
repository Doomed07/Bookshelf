package core_domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

type ShelfBook struct {
	UserID, BookID, Version int
	Read                    bool
	Rating                  *int
	Review                  *string
	AddedAt                 time.Time
	ReadAt                  *time.Time
}

func NewShelfBook(
	userID, bookID, version int,
	read bool,
	rating *int,
	review *string,
	addedAt time.Time,
	readAt *time.Time,
) ShelfBook {
	return ShelfBook{
		UserID:  userID,
		BookID:  bookID,
		Version: version,
		Read:    read,
		Rating:  rating,
		Review:  review,
		AddedAt: addedAt,
		ReadAt:  readAt,
	}
}

func (book *ShelfBook) ReadDuration() *time.Duration {
	if !book.Read {
		return nil
	}

	if book.ReadAt == nil {
		return nil
	}

	duration := book.ReadAt.Sub(book.AddedAt)

	return &duration
}

type ShelfBookWithBook struct {
	ShelfBook ShelfBook
	Book      Book
}

func NewShelfBookWithBook(
	userID, bookID, version int,
	read bool,
	rating *int,
	review *string,
	addedAt time.Time,
	readAt *time.Time,
	id int,
	title string,
	author string,
	year int,
	pages int,
	genres []string,
	description string,
	score *int,
	readsCount int,
) ShelfBookWithBook {
	return ShelfBookWithBook{
		ShelfBook: NewShelfBook(userID, bookID, version, read,
			rating, review, addedAt, readAt),
		Book: NewBook(id, title, author, year, pages, genres,
			description, score, readsCount),
	}
}

type ShelfBookPatch struct {
	Read   Nullable[bool]
	Rating Nullable[int]
	Review Nullable[string]
}

func NewShelfBookPatch(
	read Nullable[bool],
	rating Nullable[int],
	review Nullable[string],
) ShelfBookPatch {
	return ShelfBookPatch{
		Read:   read,
		Rating: rating,
		Review: review,
	}
}

func (p *ShelfBookPatch) Validate() error {
	if p.Read.Set && p.Read.Value == nil {
		return fmt.Errorf("'read' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument)
	}

	if p.Rating.Set && p.Rating.Value != nil {
		if r := *p.Rating.Value; r < 1 || r > 100 {
			return fmt.Errorf("invalid 'rating': %d: %w",
				r, core_errors.ErrInvalidArgument)
		}
	}

	if p.Review.Set && p.Review.Value != nil {
		reviewLen := utf8.RuneCountInString(strings.TrimSpace(*p.Review.Value))
		if reviewLen < 1 || reviewLen > 5000 {
			return fmt.Errorf("invalid 'review' len: %d: %w",
				reviewLen, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

func (s *ShelfBook) ApplyPatch(patch ShelfBookPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate shelfbook patch: %w", err)
	}

	temp := *s

	if patch.Read.Set {
		switch {
		case *patch.Read.Value && !temp.Read:
			temp.markRead()
		case !*patch.Read.Value && temp.Read:
			temp.markUnread()
		}
	}

	if patch.Rating.Set {
		temp.Rating = patch.Rating.Value
	}

	if patch.Review.Set {
		if patch.Review.Value != nil {
			review := strings.TrimSpace(*patch.Review.Value)
			temp.Review = &review
		} else {
			temp.Review = patch.Review.Value
		}
	}

	if !temp.Read && (temp.Rating != nil || temp.Review != nil) {
		return fmt.Errorf("rating and review allowed only for read book: %w",
			core_errors.ErrConflict)
	}

	*s = temp

	return nil
}

func (s *ShelfBook) markRead() {
	readAt := time.Now().UTC()
	if readAt.Before(s.AddedAt) {
		readAt = s.AddedAt
	}
	s.Read = true
	s.ReadAt = &readAt
}

func (s *ShelfBook) markUnread() {
	s.Read = false
	s.Rating = nil
	s.Review = nil
	s.ReadAt = nil
}

type Event struct {
	Name   string
	BookID int
	Title  string
	Author string
	Rating *int
	Review *string
	At     time.Time
}

func NewEvent(name, title, author string,
	bookID int,
	rating *int,
	review *string,
	at time.Time,
) Event {
	return Event{
		Name:   name,
		BookID: bookID,
		Title:  title,
		Author: author,
		Rating: rating,
		Review: review,
		At:     at,
	}
}
