package users_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers     godoc
// @Summary     List users
// @Description Get a paginated list of users
// @Tags        users
// @Produce     json
// @Param       limit  query int false "Limit"
// @Param       offset query int false "Offset"
// @Success     200 {object} GetUsersResponse "Users found"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users [get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query param")
		return
	}

	usersDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(usersDomains))

	responseHandler.JSONResponse(http.StatusOK, response)
}
