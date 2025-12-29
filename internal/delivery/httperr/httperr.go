package httperr

import (
	"errors"
	"net/http"

	"github.com/celpung/gocleanarch/internal/entity"
)

// HTTPError represents a normalized HTTP error response.
type HTTPError struct {
	Status  int
	Message string
	Err     error
}

// MapError translates domain/usecase errors into HTTP status codes and messages.
func MapError(err error) HTTPError {
	if err == nil {
		return HTTPError{Status: http.StatusInternalServerError, Message: "internal server error"}
	}

	switch {
	case errors.Is(err, entity.ErrEmailExists):
		return HTTPError{Status: http.StatusConflict, Message: "email already exists", Err: err}
	case errors.Is(err, entity.ErrEmailNotFound):
		return HTTPError{Status: http.StatusNotFound, Message: "email not found", Err: err}
	case errors.Is(err, entity.ErrUserNotFound):
		return HTTPError{Status: http.StatusNotFound, Message: "user not found", Err: err}
	case errors.Is(err, entity.ErrPasswordMismatch):
		return HTTPError{Status: http.StatusUnauthorized, Message: "wrong password", Err: err}
	case errors.Is(err, entity.ErrPasswordRequired):
		return HTTPError{Status: http.StatusBadRequest, Message: "password is required", Err: err}
	case errors.Is(err, entity.ErrUserIDRequired):
		return HTTPError{Status: http.StatusBadRequest, Message: "user id is required", Err: err}
	case errors.Is(err, entity.ErrInvalidInput):
		return HTTPError{Status: http.StatusBadRequest, Message: "invalid input", Err: err}
	case errors.Is(err, entity.ErrNoChanges):
		return HTTPError{Status: http.StatusBadRequest, Message: "no changes to update", Err: err}
	case errors.Is(err, entity.ErrNameRequired):
		return HTTPError{Status: http.StatusBadRequest, Message: "name is required", Err: err}
	case errors.Is(err, entity.ErrEmailRequired):
		return HTTPError{Status: http.StatusBadRequest, Message: "email is required", Err: err}
	case errors.Is(err, entity.ErrRoleRequired):
		return HTTPError{Status: http.StatusBadRequest, Message: "role is required", Err: err}
	}

	return HTTPError{Status: http.StatusInternalServerError, Message: "internal server error", Err: err}
}
