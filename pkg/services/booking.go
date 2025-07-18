package services

import (
	"encoding/json"
	"io"
	"net/http"
	"log"
	"request-offer/pkg/models"
)

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	// This function will handle the booking creation logic
	// For now, we will just return a success message
	var booking models.Booking
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &booking)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received booking request: %+v", booking)
	w.WriteHeader(http.StatusCreated)
}