package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

func getIntQueryParam(r *http.Request, key string) (*int, error) {
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

func GetStrQueryParam(r *http.Request, key string) *string {
	param := strings.TrimSpace(r.URL.Query().Get(key))
	if param == "" {
		return nil
	}

	return &param
}

func GetLimitOffsetQueryParam(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	limit, err := getIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := getIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return limit, offset, nil
}

func getBoolQueryParam(r *http.Request, key string) (*bool, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	queryParam, err := strconv.ParseBool(param)
	if err != nil {
		return nil, fmt.Errorf("invalid bool param: %s; key: %s; err: %v; typeErr:%w",
			param, key, err, core_errors.ErrInvalidArgument)
	}

	return &queryParam, nil
}

func GetReadQueryParam(r *http.Request) (*bool, error) {
	const readQueryParam = "read"

	read, err := getBoolQueryParam(r, readQueryParam)
	if err != nil {
		return nil, fmt.Errorf("get read query param: %w", err)
	}

	return read, nil
}
