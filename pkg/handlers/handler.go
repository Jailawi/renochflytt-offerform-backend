package handlers

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Ren och Flytt API!"))
	})

	mux.HandleFunc("/booking", CreateBooking)

}
