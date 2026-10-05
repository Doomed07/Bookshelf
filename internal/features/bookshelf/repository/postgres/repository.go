package bookshelf_repository_postgres

import core_postgres_pool "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool"

type BookshelfRepository struct {
	pool core_postgres_pool.Pool
}

func NewBookshelfRepository(pool core_postgres_pool.Pool) *BookshelfRepository {
	return &BookshelfRepository{
		pool: pool,
	}
}
