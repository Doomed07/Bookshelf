package statistics_service

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
)

func day(d int) time.Time {
	return time.Date(2026, 1, d, 12, 0, 0, 0, time.UTC)
}

func TestInPeriod(t *testing.T) {
	from := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		t    time.Time
		from *time.Time
		to   *time.Time
		want bool
	}{
		{name: "no bounds", t: day(1), want: true},
		{name: "before from", t: from.Add(-time.Second), from: &from, want: false},
		{name: "exactly from", t: from, from: &from, want: true},
		{name: "inside period", t: day(15), from: &from, to: &to, want: true},
		// 'to' — дата включительно: весь день 20 января входит в период
		{name: "end of 'to' day", t: to.Add(24*time.Hour - time.Second), to: &to, want: true},
		{name: "start of next day after 'to'", t: to.Add(24 * time.Hour), to: &to, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inPeriod(tt.t, tt.from, tt.to); got != tt.want {
				t.Errorf("inPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFavorite(t *testing.T) {
	tests := []struct {
		name   string
		counts map[string]int
		want   *string
	}{
		{name: "empty map", counts: map[string]int{}, want: nil},
		{name: "single entry", counts: map[string]int{"Фэнтези": 1}, want: new("Фэнтези")},
		{name: "max wins", counts: map[string]int{"Фэнтези": 1, "Драма": 3}, want: new("Драма")},
		// при равенстве — первый по алфавиту, чтобы ответ был стабильным
		{name: "tie broken alphabetically", counts: map[string]int{"Б": 2, "А": 2, "В": 2}, want: new("А")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := favorite(tt.counts)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("favorite() = %v, want %v", deref(got), deref(tt.want))
			}
		})
	}
}

func TestCalcBookStatistics(t *testing.T) {
	from := day(10)
	to := day(20)

	shelves := []core_domain.ShelfBook{
		// добавил и прочитал в периоде, оценка + рецензия
		{AddedAt: day(11), Read: true, ReadAt: new(day(12)), Rating: new(80), Review: new("Отлично")},
		// добавил и прочитал в периоде, только оценка
		{AddedAt: day(11), Read: true, ReadAt: new(day(13)), Rating: new(60)},
		// добавил в периоде, но не прочитал
		{AddedAt: day(15)},
		// добавил до периода, прочитал в периоде — считается в users_read, но не в users_on_shelf
		{AddedAt: day(1), Read: true, ReadAt: new(day(14)), Rating: new(80)},
		// добавил и прочитал до периода — не считается вообще
		{AddedAt: day(1), Read: true, ReadAt: new(day(2)), Rating: new(10), Review: new("Плохо")},
	}

	tests := []struct {
		name    string
		shelves []core_domain.ShelfBook
		from    *time.Time
		to      *time.Time
		want    core_domain.BookStats
	}{
		{
			name:    "no shelves",
			shelves: nil,
			want: core_domain.BookStats{
				RatingDistribution: []core_domain.RatingCount{},
			},
		},
		{
			name:    "whole period",
			shelves: shelves,
			want: core_domain.BookStats{
				UsersAddedOnShelf: 5,
				UsersRead:         4,
				BookAverageRate:   new((80.0 + 60 + 80 + 10) / 4),
				RatingDistribution: []core_domain.RatingCount{
					{Rating: 10, Count: 1},
					{Rating: 60, Count: 1},
					{Rating: 80, Count: 2},
				},
				ReviewCount: 2,
			},
		},
		{
			name:    "filtered by period",
			shelves: shelves,
			from:    &from,
			to:      &to,
			want: core_domain.BookStats{
				UsersAddedOnShelf: 3,
				UsersRead:         3,
				BookAverageRate:   new((80.0 + 60 + 80) / 3),
				RatingDistribution: []core_domain.RatingCount{
					{Rating: 60, Count: 1},
					{Rating: 80, Count: 2},
				},
				ReviewCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcBookStatistics(tt.shelves, tt.from, tt.to)
			// reflect.DeepEqual сравнивает структуры целиком, включая
			// значения под указателями и содержимое слайсов
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcBookStatistics() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestCalcUserStatistics(t *testing.T) {
	from := day(10)
	to := day(20)

	book := func(author string, genres ...string) core_domain.Book {
		return core_domain.Book{Author: author, Genres: genres}
	}

	shelf := []core_domain.ShelfBookWithBook{
		{
			// добавил и прочитал в периоде за 24 часа, есть рецензия
			ShelfBook: core_domain.ShelfBook{AddedAt: day(11), Read: true, ReadAt: new(day(12)), Review: new("Ок")},
			Book:      book("Толкин", "Фэнтези", "Приключения"),
		},
		{
			// добавил и прочитал в периоде за 48 часов
			ShelfBook: core_domain.ShelfBook{AddedAt: day(11), Read: true, ReadAt: new(day(13))},
			Book:      book("Толкин", "Фэнтези"),
		},
		{
			// добавил в периоде, не прочитал
			ShelfBook: core_domain.ShelfBook{AddedAt: day(15)},
			Book:      book("Оруэлл", "Антиутопия"),
		},
		{
			// добавил до периода, прочитал до периода
			ShelfBook: core_domain.ShelfBook{AddedAt: day(1), Read: true, ReadAt: new(day(3))},
			Book:      book("Оруэлл", "Антиутопия"),
		},
		{
			// read=true, но ReadAt пустой: на полке, но не считается прочитанной
			ShelfBook: core_domain.ShelfBook{AddedAt: day(11), Read: true, ReadAt: nil},
			Book:      book("Кинг", "Ужасы"),
		},
	}

	tests := []struct {
		name  string
		books []core_domain.ShelfBookWithBook
		from  *time.Time
		to    *time.Time
		want  core_domain.UserStats
	}{
		{
			name:  "empty shelf — nullable fields are nil",
			books: nil,
			want:  core_domain.UserStats{},
		},
		{
			name:  "filtered by period",
			books: shelf,
			from:  &from,
			to:    &to,
			want: core_domain.UserStats{
				BooksOnShelf:            4,
				BooksRead:               2,
				BooksReadRate:           new(2.0 / 4 * 100),
				BooksAverageReadTime:    new(36 * time.Hour), // (24ч + 48ч) / 2
				BooksReadFavoriteGenre:  new("Фэнтези"),
				BooksReadFavoriteAuthor: new("Толкин"),
				UserReviews:             1,
			},
		},
		{
			name:  "whole period",
			books: shelf,
			want: core_domain.UserStats{
				BooksOnShelf:            5,
				BooksRead:               3,
				BooksReadRate:           new(3.0 / 5 * 100),
				BooksAverageReadTime:    new(40 * time.Hour), // (24ч + 48ч + 48ч) / 3
				BooksReadFavoriteGenre:  new("Фэнтези"),
				BooksReadFavoriteAuthor: new("Толкин"),
				UserReviews:             1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcUserStatistics(tt.books, tt.from, tt.to)

			// float64 нельзя сравнивать через ==: 2.0/3*100, посчитанное
			// компилятором как константа, и то же выражение в рантайме
			// могут отличаться в последнем бите. Сравниваем с погрешностью,
			// а остальные поля — как обычно.
			if !almostEqualPtr(got.BooksReadRate, tt.want.BooksReadRate) {
				t.Errorf("BooksReadRate = %v, want %v", deref(got.BooksReadRate), deref(tt.want.BooksReadRate))
			}
			got.BooksReadRate, tt.want.BooksReadRate = nil, nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcUserStatistics() = %s, want %s", formatUserStats(got), formatUserStats(tt.want))
			}
		})
	}
}

func almostEqualPtr(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return math.Abs(*a-*b) < 1e-9
}

func formatUserStats(s core_domain.UserStats) string {
	return fmt.Sprintf("{BooksOnShelf:%d BooksRead:%d AvgReadTime:%v FavoriteGenre:%v FavoriteAuthor:%v UserReviews:%d}",
		s.BooksOnShelf, s.BooksRead, deref(s.BooksAverageReadTime),
		deref(s.BooksReadFavoriteGenre), deref(s.BooksReadFavoriteAuthor), s.UserReviews)
}

func deref[T any](p *T) any {
	if p == nil {
		return nil
	}
	return *p
}
