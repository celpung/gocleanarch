package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/delivery/middleware"
	apperrors "github.com/celpung/gocleanarch/internal/domain/errors"
	usecasedto "github.com/celpung/gocleanarch/internal/usecase/dto"
	usecaseport "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	"github.com/celpung/gocleanarch/pkg/httpx"
	"github.com/celpung/gocleanarch/pkg/mapper"
	"github.com/celpung/gocleanarch/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	usecase usecaseport.UserUsecase
	copier  mapper.Copier
}

func NewUserHandler(usecase usecaseport.UserUsecase, copier mapper.Copier) *UserHandler {
	if copier == nil {
		copier = mapper.DefaultCopier{}
	}
	return &UserHandler{usecase: usecase, copier: copier}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	user := usecasedto.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	if _, err := h.usecase.Create(r.Context(), user); err != nil {
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

	mapped, err := mapper.MapStructListDTOWith[usecasedto.UserResponse, dto.UserResponse](h.copier, users)
	if err != nil {
		writeBadRequest(w, "failed to map user response")
		return
	}

	respUsers := make([]dto.UserResponse, 0, len(mapped))
	for _, u := range mapped {
		if u != nil {
			respUsers = append(respUsers, *u)
		}
	}

	resp := dto.ListUsersResponse{
		Data: respUsers,
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

	updateReq := usecasedto.UpdateUserRequest{
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
	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, apperrors.ErrEmailExists):
		status, msg = http.StatusConflict, "email already exists"
	case errors.Is(err, apperrors.ErrEmailNotFound):
		status, msg = http.StatusNotFound, "email not found"
	case errors.Is(err, apperrors.ErrUserNotFound):
		status, msg = http.StatusNotFound, "user not found"
	case errors.Is(err, apperrors.ErrPasswordMismatch):
		status, msg = http.StatusUnauthorized, "wrong password"
	case errors.Is(err, apperrors.ErrPasswordRequired):
		status, msg = http.StatusBadRequest, "password is required"
	case errors.Is(err, apperrors.ErrUserIDRequired):
		status, msg = http.StatusBadRequest, "user id is required"
	case errors.Is(err, apperrors.ErrInvalidInput):
		status, msg = http.StatusBadRequest, "invalid input"
	case errors.Is(err, apperrors.ErrNoChanges):
		status, msg = http.StatusBadRequest, "no changes to update"
	case errors.Is(err, apperrors.ErrNameRequired):
		status, msg = http.StatusBadRequest, "name is required"
	case errors.Is(err, apperrors.ErrEmailRequired):
		status, msg = http.StatusBadRequest, "email is required"
	case errors.Is(err, apperrors.ErrRoleRequired):
		status, msg = http.StatusBadRequest, "role is required"
	}

	httpx.WriteJSON(w, status, map[string]string{"message": msg})
}

func (h *UserHandler) toDeliveryUserResponse(user usecasedto.UserResponse) (dto.UserResponse, error) {
	var resp dto.UserResponse
	if err := h.copier.Copy(&resp, &user); err != nil {
		return dto.UserResponse{}, err
	}
	return resp, nil
}
