package utils

import (
	"net/http"
)

// MethodCheck is a reusable function to check if the request method is allowed.
func RequestMethodCheck(w http.ResponseWriter, r *http.Request, allowedMethods ...string) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}

	// If the method is not allowed, respond with a 405 Method Not Allowed status.
	w.Header().Set("Allow", allowedMethods[0])
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return false
}
