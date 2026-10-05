package books_service

import (
	"context"
	"errors"
	"testing"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

// fakeBooksRepository — подставной репозиторий: запоминает, с каким limit и offset его
// вызвали, и возвращает заранее заданный ответ. База для теста сервиса не нужна.
type fakeBooksRepository struct {
	gotLimit  int
	gotOffset int
	reviews   []core_domain.ReviewWithBook
	err       error
}

func (f *fakeBooksRepository) GetBooks(ctx context.Context, title, author *string, limit, offset int) ([]core_domain.Book, error) {
	return nil, nil
}

func (f *fakeBooksRepository) GetBook(ctx context.Context, id int) (core_domain.Book, error) {
	return core_domain.Book{}, nil
}

func (f *fakeBooksRepository) GetReviews(ctx context.Context, id, limit, offset int) ([]core_domain.Review, error) {
	return nil, nil
}

func (f *fakeBooksRepository) GetRecentReviews(ctx context.Context, limit, offset int) ([]core_domain.ReviewWithBook, error) {
	f.gotLimit, f.gotOffset = limit, offset
	return f.reviews, f.err
}

func TestBooksService_GetRecentReviews(t *testing.T) {
	repoErr := errors.New("db is down")

	tests := []struct {
		name       string
		limit      *int
		offset     *int
		repoErr    error
		wantLimit  int
		wantOffset int
		wantErr    error
	}{
		{name: "default limit", limit: nil, wantLimit: core_domain.DefaultLimit},
		{name: "custom limit", limit: new(6), wantLimit: 6},
		{name: "zero limit", limit: new(0), wantErr: core_errors.ErrInvalidArgument},
		{name: "limit above max", limit: new(core_domain.MaxLimit + 1), wantErr: core_errors.ErrInvalidArgument},
		{name: "repository error is wrapped", limit: new(6), repoErr: repoErr, wantErr: repoErr},
		{name: "custom limit and offset", limit: new(6), offset: new(12), wantLimit: 6, wantOffset: 12},
		{name: "negative offset", offset: new(-1), wantErr: core_errors.ErrInvalidArgument},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeBooksRepository{
				reviews: []core_domain.ReviewWithBook{{}},
				err:     tt.repoErr,
			}
			s := NewBooksService(repo)

			got, gotErr := s.GetRecentReviews(context.Background(), tt.limit, tt.offset)

			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Errorf("GetRecentReviews() error = %v, want wrapped %v", gotErr, tt.wantErr)
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("GetRecentReviews() failed: %v", gotErr)
			}
			if repo.gotLimit != tt.wantLimit {
				t.Errorf("repository limit = %d, want %d", repo.gotLimit, tt.wantLimit)
			}
			if repo.gotOffset != tt.wantOffset {
				t.Errorf("repository offset = %d, want %d", repo.gotOffset, tt.wantOffset)
			}
			if len(got) != 1 {
				t.Errorf("GetRecentReviews() returned %d reviews, want 1", len(got))
			}
		})
	}
}
