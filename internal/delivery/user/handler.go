package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	domain "github.com/celpung/gocleanarch/internal/domain/user"
	usecase "github.com/celpung/gocleanarch/internal/usecase/user"
	"github.com/celpung/gocleanarch/pkg/httpx"
	"github.com/celpung/gocleanarch/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	usecase usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	userInput := usecase.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	if _, err := h.usecase.Create(r.Context(), userInput); err != nil {
		respondError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user registered"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	token, err := h.usecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok || strings.TrimSpace(userID) == "" {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	if err := h.usecase.ChangePassword(r.Context(), userID, req.Password); err != nil {
		respondError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

	users, total, err := h.usecase.List(r.Context(), page, limit)
	if err != nil {
		respondError(w, err)
		return
	}

	respUsers := make([]UserResponse, 0, len(users))
	for _, u := range users {
		respUsers = append(respUsers, UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	resp := ListUsersResponse{
		Data: respUsers,
		Meta: PagingMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
		},
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeBadRequest(w, "user id is required")
		return
	}

	var req UpdateUserRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	if req.Name == nil && req.Email == nil && req.Role == nil && req.Password == nil {
		writeBadRequest(w, "no changes to update")
		return
	}

	updateReq := usecase.UpdateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	if _, err := h.usecase.Update(r.Context(), id, updateReq); err != nil {
		respondError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeBadRequest(w, "user id is required")
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		respondError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}

// ---------- helpers ----------

func decodeAndValidate(r *http.Request, dst any) error {
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

func respondError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, domain.ErrEmailExists):
		status, msg = http.StatusConflict, "email already exists"
	case errors.Is(err, domain.ErrEmailNotFound):
		status, msg = http.StatusNotFound, "email not found"
	case errors.Is(err, domain.ErrUserNotFound):
		status, msg = http.StatusNotFound, "user not found"
	case errors.Is(err, domain.ErrPasswordMismatch):
		status, msg = http.StatusUnauthorized, "wrong password"
	case errors.Is(err, domain.ErrPasswordRequired):
		status, msg = http.StatusBadRequest, "password is required"
	case errors.Is(err, domain.ErrUserIDRequired):
		status, msg = http.StatusBadRequest, "user id is required"
	case errors.Is(err, domain.ErrInvalidInput):
		status, msg = http.StatusBadRequest, "invalid input"
	case errors.Is(err, domain.ErrNoChanges):
		status, msg = http.StatusBadRequest, "no changes to update"
	case errors.Is(err, domain.ErrNameRequired):
		status, msg = http.StatusBadRequest, "name is required"
	case errors.Is(err, domain.ErrEmailRequired):
		status, msg = http.StatusBadRequest, "email is required"
	case errors.Is(err, domain.ErrInvalidEmail):
		status, msg = http.StatusBadRequest, "invalid email"
	case errors.Is(err, domain.ErrRoleRequired):
		status, msg = http.StatusBadRequest, "role is required"
	case errors.Is(err, domain.ErrInvalidRole):
		status, msg = http.StatusBadRequest, "invalid role"
	case errors.Is(err, domain.ErrWeakPassword):
		status, msg = http.StatusBadRequest, "password is too weak"
	}

	httpx.WriteJSON(w, status, map[string]string{"message": msg})
}
