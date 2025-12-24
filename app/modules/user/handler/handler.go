package handler

import (
	"encoding/json"
	"net/http"

	"github.com/celpung/gocleanarch/app/modules/user/domain/entity"
	"github.com/celpung/gocleanarch/app/modules/user/handler/dto"
	"github.com/celpung/gocleanarch/app/modules/user/usecase/port"
	"github.com/celpung/gocleanarch/pkg/helpers/httpx"
)

type UserHandler struct {
	usecase port.UserUsecase
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request body",
		})
		return
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.usecase.CreateUser(user); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	resp := dto.CreateUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func NewUserHandler(usecase port.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}
