package handler

import (
	"encoding/json"
	"net/http"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/domain/entity"
	errs "github.com/celpung/gocleanarch/internal/domain/errors"
	"github.com/celpung/gocleanarch/internal/usecase/port"
	"github.com/celpung/gocleanarch/pkg/helpers/httpx"
	"github.com/celpung/gocleanarch/pkg/utilities/validator"
)

type UserHandler struct {
	usecase port.UserUsecase
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// 1) Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	// 3) Validate
	if err := validator.Validate(req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
		return
	}

	// 4) DTO -> Entity
	in := entity.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Role:     req.Role,
	}

	created, err := h.usecase.Create(r.Context(), in)
	if err != nil {
		switch err {
		case errs.ErrEmailExists:
			httpx.WriteJSON(w, http.StatusConflict, map[string]any{
				"message": err.Error(),
			})
			return

		case errs.ErrInvalidInput,
			errs.ErrNameRequired,
			errs.ErrInvalidEmail,
			errs.ErrEmailRequired,
			errs.ErrInvalidRole,
			errs.ErrRoleRequired,
			errs.ErrWeakPassword,
			errs.ErrPasswordRequired:
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"message": err.Error(),
			})
			return

		default:
			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Internal server error",
			})
			return
		}
	}

	// 5) Entity -> Response DTO
	resp, err := dto.ToUserResponse(created)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Internal server error",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "User registered",
		"data":    resp,
	})
}

func NewUserHandler(usecase port.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}
