package books_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetRecentReviewsResponse []RecentReviewDTOResponse

// GetRecentReviews godoc
// @Summary     List recent reviews
// @Description Get the latest written reviews of all users together with their books, newest first (by read date)
// @Tags        books
// @Produce     json
// @Param       limit query int false "Limit"
// @Param       offset query int false "Offset"
// @Success     200 {object} GetRecentReviewsResponse "Reviews found"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /reviews [get]
func (h *BooksHTTPHandler) GetRecentReviews(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: limit/offset")
		return
	}

	domainReviews, err := h.booksService.GetRecentReviews(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get recent reviews")
		return
	}

	response := GetRecentReviewsResponse(recentReviewsDTOFromDomains(domainReviews))
	responseHandler.JSONResponse(http.StatusOK, response)
}
