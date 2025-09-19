package delivery_impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/celpung/gocleanarch/application/user/domain/entity"
	"github.com/celpung/gocleanarch/application/user/domain/usecase"
	"github.com/celpung/gocleanarch/delivery/dto"
	delivery "github.com/celpung/gocleanarch/delivery/std/chi/user"
	"github.com/celpung/gocleanarch/delivery/std/chi/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/mapper"
	"github.com/celpung/gocleanarch/infrastructure/validation"
	"github.com/go-chi/chi/v5"
)

type UserDeliveryStruct struct {
	UserUsecase usecase.UserUsecase
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (d *UserDeliveryStruct) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	var e entity.User
	if err := mapper.CopyTo(&req, &e); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to map request",
			"error":   err.Error(),
		})
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	user, err := d.UserUsecase.Create(&e)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	var res dto.UserResponse
	if err := mapper.CopyTo(user, &res); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to map response",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Register success",
		"user":    res,
	})
}

func (d *UserDeliveryStruct) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid login data",
			"error":   err.Error(),
		})
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	token, err := d.UserUsecase.Login(req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Login failed",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Login success",
		"token":   token,
	})
}

func (d *UserDeliveryStruct) GetAllUserData(w http.ResponseWriter, r *http.Request) {
	// whos do this action ?
	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// the one who with this id
	fmt.Println(id)

	const (
		defaultPage  int64 = 1
		defaultLimit int64 = 10
		maxLimit     int64 = 100
	)

	// defaults
	page := defaultPage
	limit := defaultLimit

	// parse query
	if v := r.URL.Query().Get("page"); v != "" {
		if pv, err := strconv.ParseInt(v, 10, 32); err == nil && pv >= 1 {
			page = pv
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if lv, err := strconv.ParseInt(v, 10, 32); err == nil && lv >= 1 {
			limit = lv
		}
	}
	// clamp limit
	if limit > maxLimit {
		limit = maxLimit
	}

	users, total, err := d.UserUsecase.Read(uint(page), uint(limit))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to fetch user data",
			"error":   err.Error(),
		})
		return
	}

	res, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
		return
	}

	var totalPage int64
	if limit > 0 {
		totalPage = (total + limit - 1) / limit
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Users fetched successfully",
		"data": map[string]any{
			"users":        res,
			"count":        total,
			"current_page": page,
			"total_page":   totalPage,
		},
	})
}

func (d *UserDeliveryStruct) SearchUser(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var page, limit int64
	var err error

	if pageStr != "" {
		page, err = strconv.ParseInt(pageStr, 10, 32)
		if err != nil || page < 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid page parameter",
			})
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limit < 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid limit parameter",
			})
			return
		}
	}

	keyword := r.URL.Query().Get("q")

	users, total, err := d.UserUsecase.Search(uint(page), uint(limit), keyword)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Data not found",
		})
		return
	}

	res, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Users fetched successfully",
		"data": map[string]any{
			"users":        res,
			"count":        total,
			"current_page": page,
			"total_page":   (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (d *UserDeliveryStruct) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UserUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid update data",
			"error":   err.Error(),
		})
		return
	}

	var payload entity.UpdateUserPayload
	if err := mapper.CopyTo(&req, &payload); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to map update payload",
			"error":   err.Error(),
		})
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	user, err := d.UserUsecase.Update(&payload)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
		return
	}

	var resp dto.UserResponse
	if err := mapper.CopyTo(user, &resp); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to map response",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "User updated successfully",
		"user":    resp,
	})
}

func (d *UserDeliveryStruct) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Missing id parameter",
		})
		return
	}

	if err := d.UserUsecase.SoftDelete(userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "User deleted successfully",
	})
}

func NewUserDelivery(usecase usecase.UserUsecase) delivery.UserDelivery {
	return &UserDeliveryStruct{
		UserUsecase: usecase,
	}
}
