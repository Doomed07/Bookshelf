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
	Read   core_http_types.Nullable[bool]   `json:"read"`
	Rating core_http_types.Nullable[int]    `json:"rating"`
	Review core_http_types.Nullable[string] `json:"review"`
}

func shelfBookPatchDomainFromRequest(r PatchShelfBookRequest) core_domain.ShelfBookPatch {
	return core_domain.NewShelfBookPatch(
		r.Read.ToDomain(),
		r.Rating.ToDomain(),
		r.Review.ToDomain(),
	)
}

type PatchBookResponse ShelfBookWithBookDTOResponse

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
