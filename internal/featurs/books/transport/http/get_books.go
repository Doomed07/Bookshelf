package books_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetBooksResponse []BookDTOResponse

func (h *BooksHTTPHandler) GetBooks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	title := core_http_request.GetStrQueryParam(r, "title")
	author := core_http_request.GetStrQueryParam(r, "author")

	limit, offset, err := core_http_request.GetLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: limit/offset")
		return
	}

	books, err := h.booksService.GetBooks(ctx, title, author, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get books")
		return
	}

	response := GetBooksResponse(booksDTOFromDomains(books))

	responseHandler.JSONResponse(http.StatusOK, response)

}
