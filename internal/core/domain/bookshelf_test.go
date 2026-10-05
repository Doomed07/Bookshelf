package core_domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func TestShelfBookPatch_Validate(t *testing.T) {
	tests := []struct {
		name    string
		read    Nullable[bool]
		rating  Nullable[int]
		review  Nullable[string]
		wantErr bool
	}{
		{
			name:    "empty patch",
			wantErr: false,
		},
		{
			name:    "valid read, rating and review",
			read:    Nullable[bool]{Set: true, Value: new(true)},
			rating:  Nullable[int]{Set: true, Value: new(80)},
			review:  Nullable[string]{Set: true, Value: new("Отличная книга")},
			wantErr: false,
		},
		{
			name:    "read set to NULL",
			read:    Nullable[bool]{Set: true, Value: nil},
			wantErr: true,
		},
		{
			name:    "rating set to NULL is allowed",
			rating:  Nullable[int]{Set: true, Value: nil},
			wantErr: false,
		},
		{
			name:    "rating lower bound",
			rating:  Nullable[int]{Set: true, Value: new(1)},
			wantErr: false,
		},
		{
			name:    "rating upper bound",
			rating:  Nullable[int]{Set: true, Value: new(100)},
			wantErr: false,
		},
		{
			name:    "rating below range",
			rating:  Nullable[int]{Set: true, Value: new(0)},
			wantErr: true,
		},
		{
			name:    "rating above range",
			rating:  Nullable[int]{Set: true, Value: new(101)},
			wantErr: true,
		},
		{
			name:    "review set to NULL is allowed",
			review:  Nullable[string]{Set: true, Value: nil},
			wantErr: false,
		},
		{
			name:    "review only whitespace",
			review:  Nullable[string]{Set: true, Value: new("   ")},
			wantErr: true,
		},
		{
			name:    "review max length",
			review:  Nullable[string]{Set: true, Value: new(strings.Repeat("я", 5000))},
			wantErr: false,
		},
		{
			name:    "review too long",
			review:  Nullable[string]{Set: true, Value: new(strings.Repeat("я", 5001))},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewShelfBookPatch(tt.read, tt.rating, tt.review)
			gotErr := p.Validate()

			if tt.wantErr {
				if gotErr == nil {
					t.Fatal("Validate() succeeded unexpectedly")
				}
				if !errors.Is(gotErr, core_errors.ErrInvalidArgument) {
					t.Errorf("Validate() error = %v, want wrapped core_errors.ErrInvalidArgument", gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("Validate() failed: %v", gotErr)
			}
		})
	}
}

func TestShelfBook_ApplyPatch(t *testing.T) {
	addedAt := time.Date(2026, 1, 10, 9, 0, 0, 0, time.UTC)
	readAt := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	reviewedAt := time.Date(2026, 1, 20, 18, 30, 0, 0, time.UTC)

	// функции, а не переменные: в ShelfBook есть указатели (Rating, Review,
	// ReadAt, ReviewedAt), и простая копия структуры делила бы их между кейсами
	unreadBook := func() ShelfBook {
		return ShelfBook{UserID: 1, BookID: 1, Version: 1, AddedAt: addedAt}
	}
	readBook := func() ShelfBook {
		return ShelfBook{
			UserID: 1, BookID: 1, Version: 1,
			Read:    true,
			Rating:  new(70),
			Review:  new("Хорошо"),
			AddedAt: addedAt,
			ReadAt:  new(readAt),
			// рецензия есть — значит, дата её публикации обязана быть (как в CHECK таблицы)
			ReviewedAt: new(reviewedAt),
		}
	}
	// прочитана и оценена, но без рецензии: ReviewedAt пуст
	ratedBook := func() ShelfBook {
		return ShelfBook{
			UserID: 1, BookID: 1, Version: 1,
			Read:    true,
			Rating:  new(70),
			AddedAt: addedAt,
			ReadAt:  new(readAt),
		}
	}

	tests := []struct {
		name      string
		base      func() ShelfBook
		patch     ShelfBookPatch
		wantErr   error
		checkBook func(t *testing.T, got ShelfBook)
	}{
		{
			name: "mark unread book as read",
			base: unreadBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{Set: true, Value: new(true)},
				Nullable[int]{},
				Nullable[string]{},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if !got.Read {
					t.Error("Read = false, want true")
				}
				if got.ReadAt == nil {
					t.Fatal("ReadAt = nil, want time")
				}
				if got.ReadAt.Before(addedAt) {
					t.Errorf("ReadAt = %v is before AddedAt = %v", got.ReadAt, addedAt)
				}
			},
		},
		{
			name: "mark read book as unread clears rating, review, readAt and reviewedAt",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{Set: true, Value: new(false)},
				Nullable[int]{},
				Nullable[string]{},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.Read || got.Rating != nil || got.Review != nil || got.ReadAt != nil || got.ReviewedAt != nil {
					t.Errorf("got %+v, want unread book without rating/review/readAt/reviewedAt", got)
				}
			},
		},
		{
			name: "read=true on already read book keeps readAt",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{Set: true, Value: new(true)},
				Nullable[int]{},
				Nullable[string]{},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.ReadAt == nil || !got.ReadAt.Equal(readAt) {
					t.Errorf("ReadAt = %v, want %v", got.ReadAt, readAt)
				}
			},
		},
		{
			name: "rating and review updated, review trimmed",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{Set: true, Value: new(95)},
				Nullable[string]{Set: true, Value: new("  Шедевр  ")},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.Rating == nil || *got.Rating != 95 {
					t.Errorf("Rating = %v, want 95", got.Rating)
				}
				if got.Review == nil || *got.Review != "Шедевр" {
					t.Errorf("Review = %v, want %q", got.Review, "Шедевр")
				}
			},
		},
		{
			name: "rating and review cleared with NULL",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{Set: true, Value: nil},
				Nullable[string]{Set: true, Value: nil},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.Rating != nil || got.Review != nil {
					t.Errorf("Rating = %v, Review = %v, want both nil", got.Rating, got.Review)
				}
				if got.ReviewedAt != nil {
					t.Errorf("ReviewedAt = %v, want nil", got.ReviewedAt)
				}
				if !got.Read {
					t.Error("Read = false, want true")
				}
			},
		},
		{
			name: "review published for the first time sets reviewedAt",
			base: ratedBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{},
				Nullable[string]{Set: true, Value: new("Рецензия")},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.ReviewedAt == nil {
					t.Fatal("ReviewedAt = nil, want time")
				}
				if got.ReviewedAt.Before(readAt) {
					t.Errorf("ReviewedAt = %v is before ReadAt = %v", got.ReviewedAt, readAt)
				}
			},
		},
		{
			// правка уже опубликованной рецензии не должна «поднимать» её как новую
			name: "review edit keeps reviewedAt",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{},
				Nullable[string]{Set: true, Value: new("Исправленный текст")},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.ReviewedAt == nil || !got.ReviewedAt.Equal(reviewedAt) {
					t.Errorf("ReviewedAt = %v, want %v", got.ReviewedAt, reviewedAt)
				}
			},
		},
		{
			name: "rating change keeps reviewedAt of existing review",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{Set: true, Value: new(90)},
				Nullable[string]{},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.ReviewedAt == nil || !got.ReviewedAt.Equal(reviewedAt) {
					t.Errorf("ReviewedAt = %v, want %v", got.ReviewedAt, reviewedAt)
				}
			},
		},
		{
			name: "review removed clears reviewedAt",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{},
				Nullable[string]{Set: true, Value: nil},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.Review != nil || got.ReviewedAt != nil {
					t.Errorf("Review = %v, ReviewedAt = %v, want both nil", got.Review, got.ReviewedAt)
				}
			},
		},
		{
			name: "rating only does not set reviewedAt",
			base: ratedBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{Set: true, Value: new(90)},
				Nullable[string]{},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.ReviewedAt != nil {
					t.Errorf("ReviewedAt = %v, want nil", got.ReviewedAt)
				}
			},
		},
		{
			name: "mark read and review in one patch sets readAt and reviewedAt",
			base: unreadBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{Set: true, Value: new(true)},
				Nullable[int]{},
				Nullable[string]{Set: true, Value: new("Рецензия")},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if got.ReadAt == nil || got.ReviewedAt == nil {
					t.Fatalf("ReadAt = %v, ReviewedAt = %v, want both set", got.ReadAt, got.ReviewedAt)
				}
				// CHECK таблицы требует reviewed_at >= read_at
				if got.ReviewedAt.Before(*got.ReadAt) {
					t.Errorf("ReviewedAt = %v is before ReadAt = %v", got.ReviewedAt, got.ReadAt)
				}
			},
		},
		{
			name: "mark read and rate in one patch",
			base: unreadBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{Set: true, Value: new(true)},
				Nullable[int]{Set: true, Value: new(60)},
				Nullable[string]{},
			),
			checkBook: func(t *testing.T, got ShelfBook) {
				if !got.Read || got.Rating == nil || *got.Rating != 60 {
					t.Errorf("got %+v, want read book with rating 60", got)
				}
			},
		},
		{
			name: "rating on unread book — conflict",
			base: unreadBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{Set: true, Value: new(50)},
				Nullable[string]{},
			),
			wantErr: core_errors.ErrConflict,
		},
		{
			name: "review on unread book — conflict",
			base: unreadBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{},
				Nullable[string]{Set: true, Value: new("Рецензия")},
			),
			wantErr: core_errors.ErrConflict,
		},
		{
			name: "unread and rate in one patch — conflict",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{Set: true, Value: new(false)},
				Nullable[int]{Set: true, Value: new(50)},
				Nullable[string]{},
			),
			wantErr: core_errors.ErrConflict,
		},
		{
			name: "invalid rating — invalid argument",
			base: readBook,
			patch: NewShelfBookPatch(
				Nullable[bool]{},
				Nullable[int]{Set: true, Value: new(150)},
				Nullable[string]{},
			),
			wantErr: core_errors.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := tt.base()
			gotErr := book.ApplyPatch(tt.patch)

			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("ApplyPatch() error = %v, want %v", gotErr, tt.wantErr)
				}
				// при ошибке книга не должна измениться
				before := tt.base()
				if book.Read != before.Read ||
					!equalPtr(book.Rating, before.Rating) ||
					!equalPtr(book.Review, before.Review) ||
					!equalPtr(book.ReadAt, before.ReadAt) ||
					!equalPtr(book.ReviewedAt, before.ReviewedAt) {
					t.Errorf("ApplyPatch() mutated book on error: got %+v, want %+v", book, before)
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("ApplyPatch() failed: %v", gotErr)
			}
			tt.checkBook(t, book)
		})
	}
}

func TestShelfBook_ReadDuration(t *testing.T) {
	addedAt := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		book ShelfBook
		want *time.Duration
	}{
		{
			name: "unread book",
			book: ShelfBook{AddedAt: addedAt},
			want: nil,
		},
		{
			name: "read book without readAt",
			book: ShelfBook{Read: true, AddedAt: addedAt},
			want: nil,
		},
		{
			name: "read book",
			book: ShelfBook{Read: true, AddedAt: addedAt, ReadAt: new(addedAt.Add(48 * time.Hour))},
			want: new(48 * time.Hour),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.book.ReadDuration()
			if !equalPtr(got, tt.want) {
				t.Errorf("ReadDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

// equalPtr сравнивает значения под указателями: два nil равны,
// nil и не-nil — нет
func equalPtr[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
