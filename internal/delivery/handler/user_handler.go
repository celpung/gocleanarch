package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/delivery/httperr"
	"github.com/celpung/gocleanarch/internal/delivery/middleware"
	usecaseport "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	userdto "github.com/celpung/gocleanarch/internal/usecase/user/dto"
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
	var req dto.RegisterRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	if _, err := h.usecase.Create(r.Context(), userdto.CreateUserRequest(req)); err != nil {
		respondError(w, err)
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
		respondError(w, err)
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
		respondError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

	users, total, err := h.usecase.List(r.Context(), page, limit)
	if err != nil {
		respondError(w, err)
		return
	}

	resp := dto.ListUsersResponse{
		Data: users,
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

	if req.Name == nil && req.Email == nil && req.Role == nil && req.Password == nil {
		writeBadRequest(w, "no changes to update")
		return
	}

	updateReq := userdto.UpdateUserRequest{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	if _, err := h.usecase.Update(r.Context(), updateReq); err != nil {
		respondError(w, err)
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
		respondError(w, err)
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

func respondError(w http.ResponseWriter, err error) {
	httpErr := httperr.MapError(err)
	httpx.WriteJSON(w, httpErr.Status, map[string]string{"message": httpErr.Message})
}
