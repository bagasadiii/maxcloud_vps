package utils

import (
	"errors"
	"net/http"
)

var (
	ErrExists     = errors.New("already exists")
	ErrNotFound   = errors.New("no data found")
	ErrDatabase   = errors.New("database error")
	ErrBadRequest = errors.New("bad request")
	ErrInternal		= errors.New("internal error")
)

func ErrCheck(err error) (int) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrDatabase):
		return http.StatusInternalServerError
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrExists):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}