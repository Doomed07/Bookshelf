package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func GetQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	queryParam, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid int: param: %s; key: %s; err: %v; typeErr: %w",
			param, key, err, core_errors.ErrInvalidArgument,
		)
	}

	return &queryParam, nil
}
