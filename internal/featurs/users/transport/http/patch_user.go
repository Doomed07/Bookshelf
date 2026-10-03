package users_transport_http

import (
	"net/http"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
	core_http_types "github.com/Doomed07/Bookshelf/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Username core_http_types.Nullable[string] `json:"username" swaggertype:"string" example:"booklover"`
	Email    core_http_types.Nullable[string] `json:"email" swaggertype:"string" example:"Max@mail.com"`
}

func userPatchFromRequest(request PatchUserRequest) core_domain.UserPatch {
	return core_domain.NewUserPatch(
		request.Username.ToDomain(),
		request.Email.ToDomain(),
	)
}

type PatchedUserResponse UserDTOResponse

// PatchUser    godoc
// @Summary     Patch user
// @Description Partially update user fields
// @Description ###Logic of updating fields:
// @Description 1. **If a field is provided in the request**: it will be updated with the new value.
// @Description 2. **If a field is not provided in the request**: it will remain unchanged.
// @Description 3. **If a field is provided with a null value**: it will result in a 400 error.
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id      path int              true "User ID"
// @Param       request body PatchUserRequest true "PatchUser request body"
// @Success     200 {object} PatchedUserResponse "Successed to patch user"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "User not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id path param")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, id, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchedUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(http.StatusOK, response)
}
