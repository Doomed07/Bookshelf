package books_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetBookResponse BookDTOResponse

// GetBook      godoc
// @Summary     Get book
// @Description Get book by ID
// @Tags        books
// @Produce     json
// @Param       id path int true "Book ID"
// @Success     200 {object} GetBookResponse "Book found"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Book not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /books/{id} [get]
func (h *BooksHTTPHandler) GetBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id from request")
		return
	}

	bookDomain, err := h.booksService.GetBook(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get book")
		return
	}

	response := GetBookResponse(BookDTOFromDomain(bookDomain))

	responseHandler.JSONResponse(http.StatusOK, response)
}
