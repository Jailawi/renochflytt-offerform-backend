package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	mongoClient *mongo.Client
	logger      *logrus.Entry
	emailSender EmailSender
	envs        *models.Envs
}

func NewBookingService(db *mongo.Client, emailSender EmailSender, logger *logrus.Entry, envs *models.Envs) *BookingService {
	return &BookingService{mongoClient: db, emailSender: emailSender, logger: logger, envs: envs}
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
type Address struct {
	Address string `json:"address"`
}

type RouteRequest struct {
	Origin        Address   `json:"origin"`
	Intermediates []Address `json:"intermediates,omitempty"`
	Destination   Address   `json:"destination"`
	TravelMode    string    `json:"travelMode"`
}

type RouteResponse struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	DistanceMeters int    `json:"distanceMeters"`
	Duration       string `json:"duration"`
}

const ORIGIN = "Helsingborg C"
const RADIUS = 6 // km
const PRICE_PER_KM_UNDER_THRESHOLD = 20
const PRICE_PER_KM_OVER_THRESHOLD = 8
const PRICE_THRESHOLD_KM = 60 // km
const URL = "https://routes.googleapis.com/directions/v2:computeRoutes"

func (s *BookingService) EstimateBooking(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Errorf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var booking models.Booking
	err = json.Unmarshal(body, &booking)
	if err != nil {
		s.logger.Errorf("Error unmarshaling JSON: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	distance, err := s.calculateDistance(booking.CurrentAddress.Address, booking.NewAddress.Address)
	if err != nil {
		s.logger.Errorf("Error calculating distance: %v", err)
		http.Error(w, "Failed to calculate distance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	distanceCost := s.calculateDistanceCost(distance)
	s.logger.Infof("Calculated distance: %d km, distance cost: %d SEK", distance, distanceCost)
	movingCost := s.calculateMovingResidenceCost(*booking.CurrentAddress.LivingArea)
	s.logger.Infof("Calculated moving cost: %d SEK for living area: %d kvm", movingCost, *booking.CurrentAddress.LivingArea)
	estimatedPrice := 0

	for _, service := range booking.Services {
		s.logger.Infof("Calculating cost for service: %s", service)
		switch service {
		case "Flytthjälp":
			accessCost := s.calculateAccessCost(booking.CurrentAddress.Accessibility, *booking.CurrentAddress.Floor)
			s.logger.Infof("Calculated access cost: %d SEK for accessibility: %s, floor: %d", accessCost, booking.CurrentAddress.Accessibility, *booking.CurrentAddress.Floor)
			estimatedPrice += distanceCost + movingCost + accessCost
			s.logger.Infof("Subtotal 1: %d SEK", estimatedPrice)
		case "Flyttstädning":
			cleaningCost := s.calculateCleaningCost(*booking.CurrentAddress.LivingArea)
			cleaningDistanceCost := int(float64(s.calculateDistanceCost(distance)) * 0.5)
			s.logger.Infof("Calculated cleaning cost: %d SEK for living area: %d kvm and distance cost: %d SEK", cleaningCost, *booking.CurrentAddress.LivingArea, cleaningDistanceCost)
			estimatedPrice += cleaningCost + cleaningDistanceCost
			s.logger.Infof("Subtotal 2: %d SEK", estimatedPrice)
		case "Packning":
			estimatedPrice += 1200 // Fixed price for packning
		case "Montering":
			estimatedPrice += 300 // Fixed price for montering
		}
	}

	s.logger.Infof("Estimated price: %d SEK for booking", estimatedPrice)
	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]int{"estimated_price": estimatedPrice}); err != nil {
		s.logger.Errorf("Failed to write response: %v", err)
	}
}

func (s *BookingService) calculateDistance(currentAddress, newAddress string) (int, error) {
	// Placeholder implementation - replace with actual distance calculation logic
	// Define the request body
	reqBody := RouteRequest{
		Origin:        Address{Address: ORIGIN},
		Intermediates: []Address{{Address: currentAddress}, {Address: newAddress}},
		Destination:   Address{Address: ORIGIN},
		TravelMode:    "DRIVE",
	}

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		return 0, err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequestWithContext(context.Background(), "POST", URL, bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", s.envs.MapsAPIKey)
	req.Header.Set("X-Goog-FieldMask", "routes.duration,routes.distanceMeters")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return 0, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return 0, err
	}

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API request failed with status %d: %s\n", resp.StatusCode, string(data))
		return 0, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Define the structure for the response
	var response RouteResponse

	// Unmarshal the response body into the response structure
	if err := json.Unmarshal(data, &response); err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return 0, err
	}

	distanceMeters := response.Routes[0].DistanceMeters
	distanceKm := int(math.Round(float64(distanceMeters) / 1000))

	return distanceKm, nil
}

func (s *BookingService) calculateDistanceCost(distanceKm int) int {
	if distanceKm <= RADIUS*3 {
		return 0
	}

	if distanceKm <= PRICE_THRESHOLD_KM {
		return distanceKm * PRICE_PER_KM_UNDER_THRESHOLD
	}

	return ((distanceKm - PRICE_THRESHOLD_KM) * PRICE_PER_KM_OVER_THRESHOLD) + (PRICE_THRESHOLD_KM * PRICE_PER_KM_UNDER_THRESHOLD)
}

func (s *BookingService) calculateMovingResidenceCost(livingArea int) int {
	switch {
	case livingArea <= 50:
		return 2900
	case livingArea > 50 && livingArea <= 70:
		return 3400
	case livingArea > 70 && livingArea <= 100:
		return 4400
	case livingArea > 100 && livingArea <= 140:
		return 5000
	case livingArea > 140:
		return 5000 + (livingArea-140)*30
	default:
		return 0
	}
}

func (s *BookingService) calculateCleaningCost(livingArea int) int {
	switch {
	case livingArea <= 50:
		return 2300
	case livingArea > 50 && livingArea <= 70:
		return 3000
	case livingArea > 70 && livingArea <= 100:
		return 3500
	case livingArea > 100 && livingArea <= 140:
		return 4000
	case livingArea > 140:
		return 4000 + (livingArea-140)*40
	default:
		return 0
	}
}

func (s *BookingService) calculateAccessCost(access string, floor int) int {
	switch access {
	case "Trappor", "Liten Hiss":
		switch {
		case floor <= 1:
			return 0
		default:
			s.logger.Infof("Add cost for access type: %s, floor: %d, cost: %d SEK", access, floor, floor*200)
			return floor * 200
		}
	default:
		return 0
	}
}
