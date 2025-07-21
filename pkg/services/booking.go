package services

import (
	"encoding/json"
	"io"
	"net/http"

	// "request-offer/pkg/app"
	"request-offer/pkg/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
) 

// EmailSender interface for sending emails
type EmailSender interface {
	SendTestEmail(to []string) error
}

type BookingService struct {
	mongoClient *mongo.Client
	logger      *logrus.Entry
	emailSender EmailSender
}

func NewBookingService(db *mongo.Client, emailSender EmailSender, logger *logrus.Entry) *BookingService {
	return &BookingService{mongoClient: db, emailSender: emailSender, logger: logger}
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
	booking.CreatedAt = time.Now()

	result, err := s.mongoClient.Database("renochflytt").Collection("bookings").InsertOne(r.Context(), booking)
	if err != nil {
		s.logger.Errorf("Error inserting booking into database: %v", err)
		http.Error(w, "Failed to save booking", http.StatusInternalServerError)
		return
	}

	s.logger.Infof("Booking inserted successfully with ID: %v", result.InsertedID)
	s.logger.Infof("Received booking request: %+v", booking)
	go s.emailSender.SendTestEmail([]string{"mustafa.al-jailawi@mail.com", "mustafa.aljailawi@gmail.com"}) // Send email notification
	w.WriteHeader(http.StatusCreated)
}
