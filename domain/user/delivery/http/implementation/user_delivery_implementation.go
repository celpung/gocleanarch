package user_delivery_implementation

import (
	"encoding/json"
	"net/http"

	user_delivery "github.com/celpung/gocleanarch/domain/user/delivery/http"
	user_usecase "github.com/celpung/gocleanarch/domain/user/usecase"
	"github.com/celpung/gocleanarch/entity"
	"github.com/celpung/gocleanarch/utils"
)

type UserDeliveryStruct struct {
	UserUsecase user_usecase.UserUsecaseInterface
}

// Register implements user_delivery.UserDeliveryInterface.
func (d *UserDeliveryStruct) Register(w http.ResponseWriter, r *http.Request) {
	if !utils.RequestMethodCheck(w, r, http.MethodPost) {
		return
	}

	var reg entity.User
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, "Failed binding json data: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := d.UserUsecase.Create(&reg)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Register success!",
		"user":    user,
	})
}

// Login implements user_delivery.UserDeliveryInterface.
func (d *UserDeliveryStruct) Login(w http.ResponseWriter, r *http.Request) {
	if !utils.RequestMethodCheck(w, r, http.MethodPost) {
		return
	}

	type UserLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var login UserLogin
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Failed to bind login data: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := d.UserUsecase.Login(login.Email, login.Password)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login success",
		"token":   token,
	})
}

// GetAllUserData implements user_delivery.UserDeliveryInterface.
func (d *UserDeliveryStruct) GetAllUserData(w http.ResponseWriter, r *http.Request) {
	if !utils.RequestMethodCheck(w, r, http.MethodGet) {
		return
	}

	user, err := d.UserUsecase.Read()
	if err != nil {
		http.Error(w, "Failed to fetch user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Success fetch user data!",
		"users":   user,
	})
}

func (d *UserDeliveryStruct) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !utils.RequestMethodCheck(w, r, http.MethodPatch) {
		return
	}

	var updateData *entity.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Failed to bind update data: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := d.UserUsecase.Update(updateData)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Success fetch user data!",
		"user":    user,
	})
}

func NewUserDelivery(usecase user_usecase.UserUsecaseInterface) user_delivery.UserDeliveryInterface {
	return &UserDeliveryStruct{
		UserUsecase: usecase,
	}
}
