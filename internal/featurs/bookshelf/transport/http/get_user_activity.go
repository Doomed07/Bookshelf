package bookshelf_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetUserActivityResponse []EventDTO

// GetUserActivity godoc
// @Summary        Get user activity feed
// @Description    Get a paginated list of the user's bookshelf activity events
// @Tags           bookshelf
// @Produce        json
// @Param          user_id path  int true "User ID"
// @Param          limit   query int false "Limit"
// @Param          offset  query int false "Offset"
// @Success        200 {object} GetUserActivityResponse "Activity events found"
// @Failure        400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure        404 {object} core_http_response.ErrorResponse "User not found"
// @Failure        500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router         /users/{user_id}/activity [get]
func (h *BookshelfHTTPHandler) GetUserActivity(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathParam(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID")
		return
	}

	limit, offset, err := core_http_request.GetLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param: limit/offset")
		return
	}

	events, err := h.bookshelfService.GetUserActivity(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user's activity")
		return
	}

	response := GetUserActivityResponse(eventsDTOFromDomain(events))
	responseHandler.JSONResponse(http.StatusOK, response)
}
