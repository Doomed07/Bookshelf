package core_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	LoginSuccess = "success"
	LoginInvalid = "invalid"
)

var (
	Registrations = promauto.With(Registry).NewCounter(prometheus.CounterOpts{
		Name: "shelfmate_registrations_total",
		Help: "Успешные регистрации.",
	})

	Logins = promauto.With(Registry).NewCounterVec(prometheus.CounterOpts{
		Name: "shelfmate_logins_total",
		Help: "Попытки входа по результату.",
	}, []string{"result"})

	ShelfBooksAdded = promauto.With(Registry).NewCounter(prometheus.CounterOpts{
		Name: "shelfmate_shelf_books_added_total",
		Help: "Книги, добавленные на полку.",
	})

	BooksMarkedRead = promauto.With(Registry).NewCounter(prometheus.CounterOpts{
		Name: "shelfmate_books_marked_read_total",
		Help: "Книги, отмеченные прочитанными.",
	})

	ReviewsPublished = promauto.With(Registry).NewCounter(prometheus.CounterOpts{
		Name: "shelfmate_reviews_published_total",
		Help: "Опубликованные рецензии (первая непустая рецензия на книгу).",
	})

	ShelfBooksRemoved = promauto.With(Registry).NewCounter(prometheus.CounterOpts{
		Name: "shelfmate_shelf_books_removed_total",
		Help: "Книги, убранные с полки.",
	})
)
