package users_transport_http

import (
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_response "github.com/Doomed07/Bookshelf/internal/core/transport/http/response"
	core_http_utils "github.com/Doomed07/Bookshelf/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCtx(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke DeleteUser handler")

	id, err := core_http_utils.GetIntPathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id path param")
		return
	}

	if err := h.usersService.DeleteUser(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.StatusCodeResponse(http.StatusNoContent)

}
