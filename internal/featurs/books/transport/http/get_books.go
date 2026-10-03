package books_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetBooksResponse []BookDTOResponse

// GetBooks     godoc
// @Summary     List books
// @Description Get a paginated list of books, optionally filtered by title/author
// @Tags        books
// @Produce     json
// @Param       title  query string false "Filter by title"
// @Param       author query string false "Filter by author"
// @Param       limit  query int    false "Limit"
// @Param       offset query int    false "Offset"
// @Success     200 {object} GetBooksResponse "Books found"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /books [get]
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
