package user_delivery

import "net/http"

type UserDeliveryInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetAllUserData(w http.ResponseWriter, r *http.Request)
}
