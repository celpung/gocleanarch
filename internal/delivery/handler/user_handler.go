package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/entity"
	usecaseport "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	"github.com/celpung/gocleanarch/pkg/helpers/httpx"
	"github.com/celpung/gocleanarch/pkg/helpers/validator"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	usecase usecaseport.UserUsecase
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserDto
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
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
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
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
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	if err := h.usecase.ChangePassword(r.Context(), req.Email, req.Password); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

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
		writeBadRequest(w, "user id is required")
		return
	}

	var req dto.UpdateUserRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	if req.Name == "" && req.Email == "" && req.Role == "" {
		writeBadRequest(w, "no field to update")
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
		writeBadRequest(w, "user id is required")
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func NewUserHandler(usecase usecaseport.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func decodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errors.New("invalid request body")
	}
	return validator.Validate(dst)
}

func parsePagination(r *http.Request) (page, limit int) {
	page, limit = 1, 10

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

	return page, limit
}

func writeBadRequest(w http.ResponseWriter, msg string) {
	httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": msg})
}

func respondUserError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case err == nil:
		return
	case errors.Is(err, entity.ErrEmailExists):
		status = http.StatusConflict
		message = "email already exists"
	case errors.Is(err, entity.ErrEmailNotFound):
		status = http.StatusNotFound
		message = "email not found"
	case errors.Is(err, entity.ErrUserNotFound):
		status = http.StatusNotFound
		message = "user not found"
	case errors.Is(err, entity.ErrPasswordMismatch):
		status = http.StatusUnauthorized
		message = "password not match"
	case errors.Is(err, entity.ErrNoFieldToUpdate):
		status = http.StatusBadRequest
		message = "no field to update"
	case errors.Is(err, entity.ErrPasswordRequired):
		status = http.StatusBadRequest
		message = "password is required"
	case errors.Is(err, entity.ErrPasswordHash), errors.Is(err, entity.ErrIDGeneration):
		status = http.StatusInternalServerError
		message = "internal server error"
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
