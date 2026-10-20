package gitverse

import (
	"errors"
	"fmt"
)

type APIError struct {
	HTTPStatus       int    `json:"-"`
	Code             string `json:"code"`
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("gitverse: %d %s: %s", e.HTTPStatus, e.Code, e.Message)
}

var (
	ErrNotFound    = errors.New("resource not found")
	ErrConflict    = errors.New("resource already exists")
	ErrRateLimited = errors.New("rate limited")
)

func (e *APIError) Classify() error {
	switch e.HTTPStatus {
	case 404:
		return fmt.Errorf("%w: %s", ErrNotFound, e.Message)
	case 409:
		return fmt.Errorf("%w: %s", ErrConflict, e.Message)
	case 429:
		return fmt.Errorf("%w: %s", ErrRateLimited, e.Message)
	default:
		return e
	}
}
