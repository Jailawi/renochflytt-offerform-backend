package handlers

import (
	"net/http"
	"request-offer/pkg/services"
)

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	services.CreateBooking(w, r)
}