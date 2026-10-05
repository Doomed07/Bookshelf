package statistics_repository_postgres

import core_postgres_pool "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool"

type StatsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatsRepository(pool core_postgres_pool.Pool) *StatsRepository {
	return &StatsRepository{
		pool: pool,
	}
}
