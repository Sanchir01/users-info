package api

import "errors"

var (
	ErrQueryString  = errors.New("query not created, check your query string")
	ErrNotFoundById = errors.New("not found by id")
)
