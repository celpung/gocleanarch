package handler

import (
	"encoding/json"
	"net/http"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/usecase"
	"github.com/celpung/gocleanarch/pkg/helpers/httpx"
	"github.com/celpung/gocleanarch/pkg/helpers/validator"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	if err := validator.Validate(req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

}
