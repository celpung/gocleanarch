package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/delivery/middleware"
	"github.com/celpung/gocleanarch/internal/entity"
	usecaseport "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	"github.com/celpung/gocleanarch/pkg/httpx"
	"github.com/celpung/gocleanarch/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	usecase usecaseport.UserUsecase
}

func NewUserHandler(usecase usecaseport.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserDto
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
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
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || strings.TrimSpace(userID) == "" {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	if err := h.usecase.ChangePassword(r.Context(), userID, req.Password); err != nil {
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
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeBadRequest(w, "user id is required")
		return
	}

	var req dto.UpdateUserRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	if req.Name == nil && req.Email == nil && req.Role == nil {
		writeBadRequest(w, "no changes to update")
		return
	}

	input := &entity.UpdateUser{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	if err := h.usecase.UpdateUser(r.Context(), id, input); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeBadRequest(w, "user id is required")
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		respondUserError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}

// ---------- helpers ----------

func decodeAndValidate(r *http.Request, dst any) error {
	// optional: enforce content-type json (kalau strict)
	// ct := r.Header.Get("Content-Type")
	// if ct != "" && !strings.Contains(ct, "application/json") {
	// 	return errors.New("content-type must be application/json")
	// }

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return errors.New("invalid request body")
	}

	if dec.More() {
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
	if err == nil {
		return
	}

	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
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
		message = "wrong password"

	case errors.Is(err, entity.ErrPasswordRequired):
		status = http.StatusBadRequest
		message = "password is required"

	case errors.Is(err, entity.ErrUserIDRequired):
		status = http.StatusBadRequest
		message = "user id is required"

	case errors.Is(err, entity.ErrInvalidInput):
		status = http.StatusBadRequest
		message = "invalid input"

	case errors.Is(err, entity.ErrNoChanges):
		status = http.StatusBadRequest
		message = "no changes to update"
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
