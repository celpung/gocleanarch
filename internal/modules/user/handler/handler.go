package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal/modules/user/domain/entity"
	"github.com/celpung/gocleanarch/internal/modules/user/handler/dto"
	usecase_dto "github.com/celpung/gocleanarch/internal/modules/user/usecase/dto"
	"github.com/celpung/gocleanarch/internal/modules/user/usecase/port"
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

	if err := h.usecase.CreateUser(r.Context(), user); err != nil {
		var vErr port.ValidationError
		status := http.StatusInternalServerError
		resp := map[string]string{
			"message": "internal server error",
		}

		if errors.As(err, &vErr) {
			status = http.StatusBadRequest
			resp["message"] = vErr.Error()
		} else if errors.Is(err, port.ErrEmailAlreadyExists) {
			status = http.StatusConflict
			resp["message"] = "email already exists"
		}

		httpx.WriteJSON(w, status, resp)
		return
	}

	resp := dto.CreateUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request body",
		})
		return
	}

	user, err := h.usecase.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		var vErr port.ValidationError
		status := http.StatusInternalServerError
		resp := map[string]string{
			"message": "internal server error",
		}

		switch {
		case errors.As(err, &vErr):
			status = http.StatusBadRequest
			resp["message"] = vErr.Error()
		case errors.Is(err, port.ErrInvalidCredentials):
			status = http.StatusUnauthorized
			resp["message"] = "invalid credentials"
		default:
		}

		httpx.WriteJSON(w, status, resp)
		return
	}

	resp := dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := 1, 10

	if v := r.URL.Query().Get("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			page = parsed
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			limit = parsed
		}
	}

	users, total, err := h.usecase.ListUsers(r.Context(), page, limit)
	if err != nil {
		var vErr port.ValidationError
		status := http.StatusInternalServerError
		resp := map[string]string{
			"message": "internal server error",
		}

		if errors.As(err, &vErr) {
			status = http.StatusBadRequest
			resp["message"] = vErr.Error()
		}

		httpx.WriteJSON(w, status, resp)
		return
	}

	items := make([]dto.UserItem, 0, len(users))
	for _, u := range users {
		items = append(items, dto.UserItem{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	resp := dto.ListUsersResponse{
		Data: items,
		Meta: dto.PagingMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
		},
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"message": "invalid request body",
		})
		return
	}

	input := usecase_dto.UpdateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.usecase.UpdateUser(r.Context(), id, input); err != nil {
		var vErr port.ValidationError
		status := http.StatusInternalServerError
		resp := map[string]string{
			"message": "internal server error",
		}

		switch {
		case errors.As(err, &vErr):
			status = http.StatusBadRequest
			resp["message"] = vErr.Error()
		case errors.Is(err, port.ErrEmailAlreadyExists):
			status = http.StatusConflict
			resp["message"] = "email already exists"
		case errors.Is(err, port.ErrUserNotFound):
			status = http.StatusNotFound
			resp["message"] = "user not found"
		case errors.Is(err, port.ErrNoFieldsToUpdate):
			status = http.StatusBadRequest
			resp["message"] = "no fields to update"
		}

		httpx.WriteJSON(w, status, resp)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.usecase.DeleteUser(r.Context(), id); err != nil {
		var vErr port.ValidationError
		status := http.StatusInternalServerError
		resp := map[string]string{
			"message": "internal server error",
		}

		switch {
		case errors.As(err, &vErr):
			status = http.StatusBadRequest
			resp["message"] = vErr.Error()
		case errors.Is(err, port.ErrUserNotFound):
			status = http.StatusNotFound
			resp["message"] = "user not found"
		}

		httpx.WriteJSON(w, status, resp)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func NewUserHandler(usecase port.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}
