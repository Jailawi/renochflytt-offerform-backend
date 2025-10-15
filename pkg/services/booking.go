package services

import (
	"encoding/json"
	"io"
	"net/http"

	// "request-offer/pkg/app"
	"request-offer/pkg/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmailSender interface for sending emails
type EmailSender interface {
	SendTestEmail(to []string, booking *models.Booking) error
}

type BookingService struct {
	mongoClient     *mongo.Client
	logger          *logrus.Entry
	emailSender     EmailSender
	routeService    *DistanceService
	priceCalculator *PriceCalculator
	envs            *models.Envs
}

func NewBookingService(db *mongo.Client, emailSender EmailSender, routeService *DistanceService, priceCalculator *PriceCalculator, logger *logrus.Entry, envs *models.Envs) *BookingService {
	return &BookingService{mongoClient: db, emailSender: emailSender, routeService: routeService, priceCalculator: priceCalculator, logger: logger, envs: envs}
}

func (s *BookingService) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var booking models.Booking
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Errorf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &booking)
	if err != nil {
		s.logger.Errorf("Error unmarshaling JSON: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Set the created_at field to current time
	booking.ID = primitive.NewObjectID() // Generate a new ObjectID
	booking.CreatedAt = time.Now()
	booking.EmailSent = false // Default to false, will be set to true after email is sent

	result, err := s.mongoClient.Database("renochflytt").Collection("bookings").InsertOne(r.Context(), booking)
	if err != nil {
		s.logger.Errorf("Error inserting booking into database: %v", err)
		http.Error(w, "Failed to save booking", http.StatusInternalServerError)
		return
	}

	s.logger.Infof("Booking inserted successfully with ID: %v", result.InsertedID)
	s.logger.Infof("Received booking request: %+v", booking)

	go s.emailSender.SendTestEmail([]string{booking.Contact.Email}, &booking) // Send email notification
	w.WriteHeader(http.StatusCreated)
}

// Define the structure for the request body


func (s *BookingService) EstimateBooking(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Errorf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var booking *models.Booking
	err = json.Unmarshal(body, &booking)
	if err != nil {
		s.logger.Errorf("Error unmarshaling JSON: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	routeDistances, err := s.routeService.GetRouteDistances(
		booking.CurrentResidence.Address,
		booking.NewResidence.Address,
	)
	if err != nil {
		s.logger.Errorf("Failed to get route distances: %v", err)
		http.Error(w, "Failed to get route distances: "+err.Error(), http.StatusInternalServerError)
		return
	}

	estimatedPrice, err := s.priceCalculator.CalculatePrice(booking, routeDistances)
	if err != nil {
		s.logger.Errorf("Failed to calculate estimated price: %v", err)
		http.Error(w, "Failed to calculate estimated price: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Infof("Estimated price: %d SEK for booking", estimatedPrice)
	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]int{"estimated_price": estimatedPrice}); err != nil {
		s.logger.Errorf("Failed to write response: %v", err)
	}
}
