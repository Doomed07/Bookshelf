package bookshelf_transport_http

import (
	"net/http"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
	core_http_types "github.com/Doomed07/Bookshelf/internal/core/transport/http/types"
)

type PatchShelfBookRequest struct {
	Read   core_http_types.Nullable[bool]   `json:"read" swaggertype:"boolean" example:"true"`
	Rating core_http_types.Nullable[int]    `json:"rating" swaggertype:"integer" example:"80"`
	Review core_http_types.Nullable[string] `json:"review" swaggertype:"string" example:"Одна из лучших книг, что я читал."`
}

func shelfBookPatchDomainFromRequest(r PatchShelfBookRequest) core_domain.ShelfBookPatch {
	return core_domain.NewShelfBookPatch(
		r.Read.ToDomain(),
		r.Rating.ToDomain(),
		r.Review.ToDomain(),
	)
}

type PatchBookResponse ShelfBookWithBookDTOResponse

// PatchBook    godoc
// @Summary     Patch bookshelf book
// @Description Partially update a book on the user's bookshelf (read status, rating, review)
// @Description ###Logic of updating fields:
// @Description 1. **If a field is provided in the request**: it will be updated with the new value.
// @Description 2. **If a field is not provided in the request**: it will remain unchanged.
// @Description 3. **If `read` is provided with a null value**: it will result in a 400 error.
// @Description 4. **If `rating` or `review` is provided with a null value**: the field will be cleared.
// @Description 5. **`rating`/`review` can only be set when the book is read** (`read=true`), otherwise a 409 error is returned.
// @Tags        bookshelf
// @Accept      json
// @Produce     json
// @Param       user_id path int                   true "User ID"
// @Param       book_id path int                   true "Book ID"
// @Param       request body PatchShelfBookRequest true "PatchBook request body"
// @Success     200 {object} PatchBookResponse "Successed to patch book"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Book not found on bookshelf"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{user_id}/bookshelf/{book_id} [patch]
func (h *BookshelfHTTPHandler) PatchBook(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchShelfBookRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request: patch book")
		return
	}

	bookPatch := shelfBookPatchDomainFromRequest(request)

	bfbs, err := h.bookshelfService.PatchBook(ctx, userID, bookID, bookPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get patched book")
		return
	}

	response := PatchBookResponse(shelfBookWithBookDTOFromDomain(bfbs))
	responseHandler.JSONResponse(http.StatusOK, response)
}
