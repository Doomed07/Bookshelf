package users_transport_http

import (
	"net/http"

	"github.com/Doomed07/Bookshelf/internal/core/domain"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Username string `json:"username" example:"book_worm07"`
	Email    string `json:"email" example:"lol@mail.com"`
}

type CreateUserResponse UserDTOResponse

// CreateUser   godoc
// @Summary     Create user
// @Description Create new user in the app
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body CreateUserRequest true "CreateUser request body"
// @Success     201 {object} CreateUserResponse "Successed to create user"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(http.StatusCreated, response)
}

func domainFromDTO(dto CreateUserRequest) core_domain.User {
	return core_domain.NewUserUninitialized(dto.Username, dto.Email)
}
