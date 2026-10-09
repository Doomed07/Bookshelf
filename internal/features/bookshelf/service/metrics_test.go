package bookshelf_service

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_metrics "github.com/Doomed07/Bookshelf/internal/core/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var errBoom = errors.New("db is down")

// fakeRepo — репозиторий в памяти на одну запись полки.
type fakeRepo struct {
	shelf      core_domain.ShelfBookWithBook
	patchCalls int

	addErr, getErr, patchErr, removeErr error
}

func (f *fakeRepo) AddBook(ctx context.Context, userID, bookID int) (core_domain.ShelfBook, error) {
	return core_domain.ShelfBook{UserID: userID, BookID: bookID}, f.addErr
}

func (f *fakeRepo) GetBook(ctx context.Context, userID, bookID int) (core_domain.ShelfBookWithBook, error) {
	return f.shelf, f.getErr
}

func (f *fakeRepo) GetBooks(ctx context.Context, userID int, read *bool, limit, offset int) ([]core_domain.ShelfBookWithBook, error) {
	return nil, nil
}

func (f *fakeRepo) UserExists(ctx context.Context, userID int) (bool, error) { return true, nil }

func (f *fakeRepo) GetUserActivity(ctx context.Context, userID int, limit, offset int) ([]core_domain.Event, error) {
	return nil, nil
}

func (f *fakeRepo) PatchBook(ctx context.Context, shelfBook core_domain.ShelfBook) (core_domain.ShelfBook, error) {
	f.patchCalls++
	if f.patchErr != nil {
		return core_domain.ShelfBook{}, f.patchErr
	}
	f.shelf.ShelfBook = shelfBook
	return shelfBook, nil
}

func (f *fakeRepo) RemoveBook(ctx context.Context, userID, bookID int) error { return f.removeErr }

type shelfCounters struct{ added, removed, markedRead, reviews float64 }

func readShelfCounters() shelfCounters {
	return shelfCounters{
		added:      testutil.ToFloat64(core_metrics.ShelfBooksAdded),
		removed:    testutil.ToFloat64(core_metrics.ShelfBooksRemoved),
		markedRead: testutil.ToFloat64(core_metrics.BooksMarkedRead),
		reviews:    testutil.ToFloat64(core_metrics.ReviewsPublished),
	}
}

func (c shelfCounters) since(before shelfCounters) shelfCounters {
	return shelfCounters{
		added:      c.added - before.added,
		removed:    c.removed - before.removed,
		markedRead: c.markedRead - before.markedRead,
		reviews:    c.reviews - before.reviews,
	}
}

func TestBookshelfService_AddBook_Metrics(t *testing.T) {
	tests := []struct {
		name    string
		addErr  error
		want    float64
		wantErr bool
	}{
		{name: "added", want: 1},
		{name: "already on the shelf", addErr: core_errors.ErrConflict, wantErr: true},
		{name: "database failure", addErr: errBoom, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookshelfService(&fakeRepo{addErr: tt.addErr})
			before := readShelfCounters()

			_, err := s.AddBook(context.Background(), 1, 2)

			if (err != nil) != tt.wantErr {
				t.Fatalf("AddBook() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := readShelfCounters().since(before); got != (shelfCounters{added: tt.want}) {
				t.Errorf("counters grew by %+v, want only added=%v", got, tt.want)
			}
		})
	}
}

func TestBookshelfService_RemoveBook_Metrics(t *testing.T) {
	tests := []struct {
		name      string
		removeErr error
		want      float64
	}{
		{name: "removed", want: 1},
		{name: "not on the shelf", removeErr: core_errors.ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewBookshelfService(&fakeRepo{removeErr: tt.removeErr})
			before := readShelfCounters()

			_ = s.RemoveBook(context.Background(), 1, 2)

			if got := readShelfCounters().since(before); got != (shelfCounters{removed: tt.want}) {
				t.Errorf("counters grew by %+v, want only removed=%v", got, tt.want)
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }

func set[T any](v T) core_domain.Nullable[T] { return core_domain.Nullable[T]{Value: &v, Set: true} }

func setNull[T any]() core_domain.Nullable[T] { return core_domain.Nullable[T]{Set: true} }

func shelfEntry(read bool, review *string) core_domain.ShelfBookWithBook {
	added := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	entry := core_domain.ShelfBook{UserID: 1, BookID: 2, Version: 1, AddedAt: added, Read: read}
	if read {
		readAt := added.Add(time.Hour)
		entry.ReadAt = &readAt
		if review != nil {
			entry.Review = review
			reviewedAt := readAt.Add(time.Hour)
			entry.ReviewedAt = &reviewedAt
		}
	}
	return core_domain.ShelfBookWithBook{ShelfBook: entry}
}

// «Прочитано» и «рецензия» считаются только при реальном переходе, а не при каждом PATCH.
func TestBookshelfService_PatchBook_Metrics(t *testing.T) {
	tests := []struct {
		name       string
		shelf      core_domain.ShelfBookWithBook
		patch      core_domain.ShelfBookPatch
		patchErr   error
		wantErr    bool
		wantRead   float64
		wantReview float64
	}{
		{
			name:  "unread -> read",
			shelf: shelfEntry(false, nil), patch: core_domain.ShelfBookPatch{Read: set(true)},
			wantRead: 1,
		},
		{
			name:  "unread -> read with a review in one request",
			shelf: shelfEntry(false, nil), patch: core_domain.ShelfBookPatch{Read: set(true), Review: set("great book")},
			wantRead: 1, wantReview: 1,
		},
		{
			name:  "first review on a read book",
			shelf: shelfEntry(true, nil), patch: core_domain.ShelfBookPatch{Review: set("great book")},
			wantReview: 1,
		},
		{
			name:  "rating alone is not a review",
			shelf: shelfEntry(true, nil), patch: core_domain.ShelfBookPatch{Rating: set(80)},
		},
		{
			name:  "editing an existing review is not a new one",
			shelf: shelfEntry(true, ptr("great book")), patch: core_domain.ShelfBookPatch{Review: set("really great book")},
		},
		{
			name:  "marking an already read book again",
			shelf: shelfEntry(true, nil), patch: core_domain.ShelfBookPatch{Read: set(true)},
		},
		{
			name:  "removing a review",
			shelf: shelfEntry(true, ptr("great book")), patch: core_domain.ShelfBookPatch{Review: setNull[string]()},
		},
		{
			name:  "read -> unread",
			shelf: shelfEntry(true, ptr("great book")), patch: core_domain.ShelfBookPatch{Read: set(false)},
		},
		{
			name:  "invalid rating",
			shelf: shelfEntry(false, nil), patch: core_domain.ShelfBookPatch{Read: set(true), Rating: set(101)},
			wantErr: true,
		},
		{
			name:  "rating on an unread book",
			shelf: shelfEntry(false, nil), patch: core_domain.ShelfBookPatch{Rating: set(5)},
			wantErr: true,
		},
		{
			name:  "database failure on save",
			shelf: shelfEntry(false, nil), patch: core_domain.ShelfBookPatch{Read: set(true), Review: set("great book")},
			patchErr: errBoom, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{shelf: tt.shelf, patchErr: tt.patchErr}
			s := NewBookshelfService(repo)
			before := readShelfCounters()

			_, err := s.PatchBook(context.Background(), 1, 2, tt.patch)

			if (err != nil) != tt.wantErr {
				t.Fatalf("PatchBook() error = %v, wantErr %v", err, tt.wantErr)
			}
			want := shelfCounters{markedRead: tt.wantRead, reviews: tt.wantReview}
			if got := readShelfCounters().since(before); got != want {
				t.Errorf("counters grew by %+v, want %+v", got, want)
			}
		})
	}
}

func TestHasText(t *testing.T) {
	if hasText(nil) || hasText(ptr("")) {
		t.Error("nil and empty must not count as text")
	}
	if !hasText(ptr("a")) {
		t.Error(`"a" must count as text`)
	}
}
