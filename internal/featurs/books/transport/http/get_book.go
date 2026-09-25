package books_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetBookResponse BookDTOResponse

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

	response := GetBookResponse(bookDTOFromDomain(bookDomain))

	responseHandler.JSONResponse(http.StatusOK, response)
}
