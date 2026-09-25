package books_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetReviewResponse []ReviewDTOResponse

func (h *BooksHTTPHandler) GetReviews(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id from request")
		return
	}

	limit, offset, err := core_http_request.GetLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset param from request")
		return
	}

	domainReviews, err := h.booksService.GetReviews(ctx, id, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get reviews")
		return
	}

	response := GetReviewResponse(reviewsDTOFromDomains(domainReviews))
	responseHandler.JSONResponse(http.StatusOK, response)
}
