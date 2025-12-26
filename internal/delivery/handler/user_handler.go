package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/usecase/port"
	"github.com/celpung/gocleanarch/pkg/helpers/httpx"
)

type UserHandler struct {
	usecase port.UserUsecase
}

func NewUserHandler(usecase port.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	input := port.CreateUserInput{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		CompanyID: req.CompanyID,
	}

	user, err := h.usecase.CreateUser(r.Context(), input)
	if err != nil {
		respondUserError(w, err)
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
	if user.CompanyID != "" {
		resp.CompanyID = &user.CompanyID
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.usecase.GetUser(r.Context(), id)
	if err != nil {
		respondUserError(w, err)
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
	if user.CompanyID != "" {
		resp.CompanyID = &user.CompanyID
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
		respondUserError(w, err)
		return
	}

	items := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		item := dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
		}
		if user.CompanyID != "" {
			item.CompanyID = &user.CompanyID
		}
		items = append(items, item)
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
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	input := port.UpdateUserInput{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		CompanyID: req.CompanyID,
		IsActive:  req.IsActive,
	}

	if err := h.usecase.UpdateUser(r.Context(), id, input); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.usecase.DeleteUser(r.Context(), id); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func respondUserError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	resp := map[string]string{
		"message": "internal server error",
	}

	var vErr port.ValidationError

	switch {
	case errors.As(err, &vErr):
		status = http.StatusBadRequest
		resp["message"] = vErr.Error()
	case errors.Is(err, port.ErrUserNotFound):
		status = http.StatusNotFound
		resp["message"] = "user not found"
	case errors.Is(err, port.ErrCompanyNotFound):
		status = http.StatusNotFound
		resp["message"] = "company not found"
	case errors.Is(err, port.ErrEmailAlreadyExists):
		status = http.StatusConflict
		resp["message"] = "email already exists"
	case errors.Is(err, port.ErrNoFieldsToUpdate):
		status = http.StatusBadRequest
		resp["message"] = "no fields to update"
	default:
	}

	httpx.WriteJSON(w, status, resp)
}
