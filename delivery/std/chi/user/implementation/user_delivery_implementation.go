package delivery_implementation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/celpung/gocleanarch/delivery/dto"
	delivery "github.com/celpung/gocleanarch/delivery/std/http/user"
	"github.com/celpung/gocleanarch/domain/user/usecase"
)

type UserDeliveryStruct struct {
	UserUsecase usecase.UserUsecaseInterface
}

func (d *UserDeliveryStruct) Register(w http.ResponseWriter, r *http.Request) {
	// Uncomment the following lines if you want to use context values for user ID, email, and role
	// userID := r.Context().Value(middlewares.ContextKeyUserID)
	// email := r.Context().Value(middlewares.ContextKeyEmail)
	// role := r.Context().Value(middlewares.ContextKeyRole)
	var req dto.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input data: "+err.Error(), http.StatusBadRequest)
		return
	}

	entity := dto.UserCreateRequestDTO(&req)

	user, err := d.UserUsecase.Create(entity)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.UserResponseDTO(user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Register success!",
		"user":    resp,
	})
}

func (d *UserDeliveryStruct) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid login data: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := d.UserUsecase.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login success",
		"token":   token,
	})
}

func (d *UserDeliveryStruct) GetAllUserData(w http.ResponseWriter, r *http.Request) {
	users, err := d.UserUsecase.Read()
	if err != nil {
		http.Error(w, "Failed to fetch user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.UserResponseListDTO(users)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success fetch user data!",
		"users":   resp,
	})
}

func (d *UserDeliveryStruct) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid update data: "+err.Error(), http.StatusBadRequest)
		return
	}

	entity := dto.UserUpdateRequestDTO(&req)

	user, err := d.UserUsecase.Update(entity)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.UserResponseDTO(user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User updated successfully!",
		"user":    resp,
	})
}

func (d *UserDeliveryStruct) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user_id: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = d.UserUsecase.SoftDelete(uint(userID))
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User deleted successfully",
	})
}

func NewUserDelivery(usecase usecase.UserUsecaseInterface) delivery.UserDeliveryInterface {
	return &UserDeliveryStruct{
		UserUsecase: usecase,
	}
}
