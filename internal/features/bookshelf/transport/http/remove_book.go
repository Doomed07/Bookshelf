package bookshelf_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

// RemoveBook   godoc
// @Summary     Remove book from bookshelf
// @Description Remove a book from the user's bookshelf
// @Tags        bookshelf
// @Param       user_id path int true "User ID"
// @Param       book_id path int true "Book ID"
// @Success     204 "Successed to remove book from bookshelf"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Book not found on bookshelf"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{user_id}/bookshelf/{book_id} [delete]
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
