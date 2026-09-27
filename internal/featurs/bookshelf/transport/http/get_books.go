package bookshelf_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetBooksResponse []ShelfBookWithBookDTOResponse

func (h *BookshelfHTTPHandler) GetBooks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathParam(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID")
		return
	}

	read, err := core_http_request.GetReadQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: read")
		return
	}

	limit, offset, err := core_http_request.GetLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: limit/offset")
		return
	}

	booksFBS, err := h.bookshelfService.GetBooks(ctx, userID, read, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get books from bookshelf")
		return
	}

	response := GetBooksResponse(booksFBSDTOFromDomain(booksFBS))
	responseHandler.JSONResponse(http.StatusOK, response)
}
