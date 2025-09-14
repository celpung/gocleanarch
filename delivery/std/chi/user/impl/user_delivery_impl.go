package delivery_impl

import (
	"encoding/json"
	"net/http"

	"github.com/celpung/gocleanarch/application/user/domain/entity"
	"github.com/celpung/gocleanarch/application/user/domain/usecase"
	"github.com/celpung/gocleanarch/delivery/dto"
	delivery "github.com/celpung/gocleanarch/delivery/std/chi/user"
	"github.com/celpung/gocleanarch/infrastructure/mapper"
	"github.com/celpung/gocleanarch/infrastructure/validation"
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
	users, err := d.UserUsecase.Read()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to fetch user data",
			"error":   err.Error(),
		})
		return
	}

	resp, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Success fetch user data",
		"users":   resp,
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
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Missing user_id parameter",
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
