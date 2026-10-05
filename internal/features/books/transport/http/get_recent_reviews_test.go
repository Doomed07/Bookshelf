package books_transport_http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	"go.uber.org/zap"
)

// fakeBooksService — подставной сервис: запоминает limit и offset и возвращает заданный ответ.
type fakeBooksService struct {
	gotLimit  *int
	gotOffset *int
	reviews   []core_domain.ReviewWithBook
	err       error
}

func (f *fakeBooksService) GetBooks(ctx context.Context, title, author *string, limit, offset *int) ([]core_domain.Book, error) {
	return nil, nil
}

func (f *fakeBooksService) GetBook(ctx context.Context, id int) (core_domain.Book, error) {
	return core_domain.Book{}, nil
}

func (f *fakeBooksService) GetReviews(ctx context.Context, id int, limit, offset *int) ([]core_domain.Review, error) {
	return nil, nil
}

func (f *fakeBooksService) GetRecentReviews(ctx context.Context, limit, offset *int) ([]core_domain.ReviewWithBook, error) {
	f.gotLimit, f.gotOffset = limit, offset
	return f.reviews, f.err
}

func TestGetRecentReviews(t *testing.T) {
	// дата публикации позже даты прочтения: именно по ней сортируется лента
	reviewedAt := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	review := core_domain.ReviewWithBook{
		Review: core_domain.Review{
			UserID: 2, Username: "anna_reads", Rating: new(88), Review: new("Отлично"),
			ReadAt:     time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC),
			ReviewedAt: new(reviewedAt),
		},
		Book: core_domain.Book{ID: 156, Title: "Задача трёх тел", Author: "Лю Цысинь", Year: 2008},
	}

	tests := []struct {
		name       string
		url        string
		serviceErr error
		wantStatus int
		wantLimit  *int
		wantOffset *int
	}{
		{name: "default limit", url: "/reviews", wantStatus: http.StatusOK, wantLimit: nil},
		{name: "custom limit", url: "/reviews?limit=6", wantStatus: http.StatusOK, wantLimit: new(6)},
		{name: "limit is not a number", url: "/reviews?limit=abc", wantStatus: http.StatusBadRequest},
		{name: "invalid limit from service", url: "/reviews?limit=0", serviceErr: core_errors.ErrInvalidArgument, wantStatus: http.StatusBadRequest, wantLimit: new(0)},
		{name: "service failure", url: "/reviews", serviceErr: errors.New("db is down"), wantStatus: http.StatusInternalServerError},
		{name: "custom limit and offset", url: "/reviews?limit=6&offset=12", wantStatus: http.StatusOK, wantLimit: new(6), wantOffset: new(12)},
		{name: "offset is not a number", url: "/reviews?offset=abc", wantStatus: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeBooksService{reviews: []core_domain.ReviewWithBook{review}, err: tt.serviceErr}
			h := NewBooksHTTPHandler(service)

			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			// в проде логгер кладёт в контекст middleware; без него FromCtx паникует
			r = r.WithContext(core_logger.ToCtx(r.Context(), &core_logger.Logger{Logger: zap.NewNop()}))
			rec := httptest.NewRecorder()

			h.GetRecentReviews(rec, r)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus != http.StatusOK {
				return
			}
			if !equalPtr(service.gotLimit, tt.wantLimit) {
				t.Errorf("service limit = %v, want %v", service.gotLimit, tt.wantLimit)
			}
			if !equalPtr(service.gotOffset, tt.wantOffset) {
				t.Errorf("service offset = %v, want %v", service.gotOffset, tt.wantOffset)
			}

			var body []struct {
				Review struct {
					Username   string     `json:"username"`
					ReviewedAt *time.Time `json:"reviewed_at"`
				} `json:"review"`
				Book struct {
					Title string `json:"title"`
				} `json:"book"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			// Fatalf: ниже идёт обращение к body[0], при пустом ответе оно упало бы паникой
			if len(body) != 1 {
				t.Fatalf("got %d reviews, want 1", len(body))
			}
			if body[0].Review.Username != "anna_reads" || body[0].Book.Title != "Задача трёх тел" {
				t.Errorf("body = %+v, want a review by anna_reads with its book", body)
			}
			if body[0].Review.ReviewedAt == nil || !body[0].Review.ReviewedAt.Equal(reviewedAt) {
				t.Errorf("reviewed_at = %v, want %v", body[0].Review.ReviewedAt, reviewedAt)
			}

		})
	}
}

func equalPtr[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
