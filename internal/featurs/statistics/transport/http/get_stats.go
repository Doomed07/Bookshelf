package statistics_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

const (
	userIDQueryParam = "user_id"
	bookIDQueryParam = "book_id"
	topQueryParam    = "top"
)

type (
	GetStatsResponse     StatisticsDTO
	GetUserStatsResponse UserStatsDTO
	GetBookStatsResponse BookStatsDTO
)

// GetStats     godoc
// @Summary     Get statistics
// @Description Get system-wide statistics by default; pass user_id for per-user statistics or book_id for per-book statistics (mutually exclusive, response shape returned is GetStatsResponse/GetUserStatsResponse/GetBookStatsResponse respectively)
// @Tags        statistics
// @Produce     json
// @Param       user_id query int    false "Get statistics for this user instead of system-wide (mutually exclusive with book_id)"
// @Param       book_id query int    false "Get statistics for this book instead of system-wide (mutually exclusive with user_id)"
// @Param       top     query int    false "Number of top books to return (system statistics only)"
// @Param       from    query string false "Filter from date, format YYYY-MM-DD"
// @Param       to      query string false "Filter to date, format YYYY-MM-DD"
// @Success     200 {object} GetStatsResponse "System statistics"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "User or book not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /stats [get]
func (h *StatsHTTPHandler) GetStats(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParam)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: user_id")
		return
	}

	bookID, err := core_http_request.GetIntQueryParam(r, bookIDQueryParam)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: book_id")
		return
	}

	top, err := core_http_request.GetIntQueryParam(r, topQueryParam)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: top")
		return
	}

	from, to, err := core_http_request.GetFromToQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: from/to")
		return
	}

	switch {
	case userID != nil && bookID != nil:
		err := fmt.Errorf("query params 'user_id' and 'book_id' can't be used together: %w",
			core_errors.ErrInvalidArgument)
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return

	case userID != nil:
		userStats, err := h.statsService.GetUserStats(ctx, *userID, from, to)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get user statistics")
			return
		}

		response := GetUserStatsResponse(userStatsDTOFromDomain(userStats))
		responseHandler.JSONResponse(http.StatusOK, response)

	case bookID != nil:
		bookStats, err := h.statsService.GetBookStats(ctx, *bookID, from, to)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get book statistics")
			return
		}

		response := GetBookStatsResponse(bookStatsDTOFromDomain(bookStats))
		responseHandler.JSONResponse(http.StatusOK, response)

	default:
		stats, err := h.statsService.GetSystemStats(ctx, from, to, top)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get system statistics")
			return
		}

		response := GetStatsResponse(statisticsDTOFromDomain(stats))
		responseHandler.JSONResponse(http.StatusOK, response)
	}

}
