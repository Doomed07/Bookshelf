package users_transport_http

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
	core_http_types "github.com/Doomed07/Bookshelf/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Username core_http_types.Nullable[string] `json:"username"`
	Email    core_http_types.Nullable[string] `json:"email"`
}

func (p *PatchUserRequest) Validate() error {
	if p.Username.Set {
		if p.Username.Value == nil {
			return fmt.Errorf("'username' can't be NULL")
		}

		usernameLen := utf8.RuneCountInString(*p.Username.Value)
		if usernameLen < 3 || usernameLen > 30 {
			return fmt.Errorf("'username' must be between 3 and 30 symbols")
		}
	}

	if p.Email.Set {
		if p.Email.Value == nil {
			return fmt.Errorf("'email' can't be NULL")
		}
		emailLen := utf8.RuneCountInString(*p.Email.Value)
		if emailLen > 254 {
			return fmt.Errorf("'email' not more then 254 symbols")
		}
	}

	return nil
}

type PatchedUserResponse UserDTOResponse

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

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Username.ToDomain(),
		request.Email.ToDomain(),
	)
}
