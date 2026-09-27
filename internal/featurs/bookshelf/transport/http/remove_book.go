package bookshelf_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

func (h *BookshelfHTTPHandler) RemoveBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathParam(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID")
		return
	}
	bookID, err := core_http_request.GetIntPathParam(r, "book_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get book ID")
		return
	}

	if err := h.bookshelfService.RemoveBook(ctx, userID, bookID); err != nil {
		responseHandler.ErrorResponse(err, "failed to remove book from shelf")
		return
	}

	responseHandler.StatusCodeResponse(http.StatusNoContent)
}
