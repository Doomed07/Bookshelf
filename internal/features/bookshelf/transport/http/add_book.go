package bookshelf_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type AddBookRequest struct {
	BookID int `json:"book_id" validate:"required,min=1" example:"1"`
}

type AddBookResponse ShelfBookWithBookDTOResponse

// AddBook      godoc
// @Summary     Add book to bookshelf
// @Description Add a book to the user's bookshelf
// @Tags        bookshelf
// @Accept      json
// @Produce     json
// @Param       user_id path int            true "User ID"
// @Param       request body AddBookRequest true "AddBook request body"
// @Success     201 {object} AddBookResponse "Successed to add book to bookshelf"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "User or book not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{user_id}/bookshelf [post]
func (h *BookshelfHTTPHandler) AddBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathParam(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID")
		return
	}

	var request AddBookRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request: get book_id")
		return
	}

	shelfbook, err := h.bookshelfService.AddBook(ctx, userID, request.BookID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to add book to your bookshelf")
		return
	}

	response := AddBookResponse(shelfBookWithBookDTOFromDomain(shelfbook))

	responseHandler.JSONResponse(http.StatusCreated, response)

}
