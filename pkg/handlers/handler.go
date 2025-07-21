package handlers

import (
	"net/http"
	"request-offer/pkg/services"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	bookingService *services.BookingService
	logger         *logrus.Entry
}

func NewHandler(bookingService *services.BookingService, logger *logrus.Entry) *Handler {
	return &Handler{
		bookingService: bookingService,
		logger:         logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.homeHandler)
	// mux.HandleFunc("/offer-estimate", h.calc.EstimateHandler)
	mux.HandleFunc("/booking", h.bookingHandler)
}

func (h *Handler) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Ren och Flytt API!"))
}

func (h *Handler) bookingHandler(w http.ResponseWriter, r *http.Request) {
	h.bookingService.CreateBooking(w, r)
}
