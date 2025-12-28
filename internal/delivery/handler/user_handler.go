package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase"
	"github.com/celpung/gocleanarch/pkg/helpers/httpx"
	"github.com/celpung/gocleanarch/pkg/helpers/validator"
	"github.com/go-chi/chi/v5"
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

	user := entity.User{
		Name:     req.Name,
		Email:    strings.TrimSpace(req.Email),
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.usecase.Register(r.Context(), user); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user registered"})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

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

	token, err := h.usecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ChangePasswordRequest

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

	if err := h.usecase.ChangePassword(r.Context(), req.Email, req.Password); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := 1, 10

	if v := r.URL.Query().Get("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	users, total, err := h.usecase.UserLists(r.Context(), page, limit)
	if err != nil {
		respondUserError(w, err)
		return
	}

	resp := dto.ListUsersResponse{
		Data: usersToResponses(users),
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
	if id == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "user id is required"})
		return
	}

	var req dto.UpdateUserRequest
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

	if req.Name == "" && req.Email == "" && req.Role == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "no field to update"})
		return
	}

	input := &entity.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	if err := h.usecase.UpdateUser(r.Context(), id, input); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "user id is required"})
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func respondUserError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case err == nil:
		return
	case strings.Contains(err.Error(), "email already exists"):
		status = http.StatusConflict
		message = "email already exists"
	case strings.Contains(err.Error(), "email not found"):
		status = http.StatusNotFound
		message = "email not found"
	case strings.Contains(err.Error(), "user not found"):
		status = http.StatusNotFound
		message = "user not found"
	case strings.Contains(err.Error(), "password not match"):
		status = http.StatusUnauthorized
		message = "password not match"
	case strings.Contains(err.Error(), "no field to update"):
		status = http.StatusBadRequest
		message = "no field to update"
	}

	httpx.WriteJSON(w, status, map[string]string{"message": message})
}

func usersToResponses(users []entity.User) []dto.UserResponse {
	result := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, dto.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		})
	}
	return result
}
