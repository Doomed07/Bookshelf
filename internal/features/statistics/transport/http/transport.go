package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
)

type StatsHTTPHandler struct {
	statsService StatsService
}

type StatsService interface {
	GetSystemStats(ctx context.Context, from, to *time.Time, top *int) (core_domain.Statistics, error)
	GetUserStats(ctx context.Context, userID int, from, to *time.Time) (core_domain.UserStats, error)
	GetBookStats(ctx context.Context, bookID int, from, to *time.Time) (core_domain.BookStats, error)
}

func NewStatsHTTPHandler(statsService StatsService) *StatsHTTPHandler {
	return &StatsHTTPHandler{
		statsService: statsService,
	}
}

func (h *StatsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/stats",
			Handler: h.GetStats,
		},
	}
}
