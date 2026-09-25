package books_repository_postgres

import core_postgres_pool "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool"

type BooksRepository struct {
	pool core_postgres_pool.Pool
}

func NewBooksRepository(pool core_postgres_pool.Pool) *BooksRepository {
	return &BooksRepository{
		pool: pool,
	}
}
