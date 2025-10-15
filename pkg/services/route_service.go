package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"request-offer/pkg/models"

	"github.com/sirupsen/logrus"
)

const (
	ORIGIN = "Helsingborg C"
    URL    = "https://routes.googleapis.com/directions/v2:computeRoutes"
)

type RouteDistances struct {
	ToCurrentAddress int // Distance from origin to current address in km
	Total            int // Total distance from origin --> current address --> new address --> origin in km
}

type DistanceService struct {
	envs   *models.Envs
	logger *logrus.Entry
}

func NewDistanceService(envs *models.Envs, logger *logrus.Entry) *DistanceService {
	return &DistanceService{
		envs:   envs,
		logger: logger,
	}
}

func (ds *DistanceService) GetRouteDistances(currentAddress, newAddress string) (*models.RouteDistances, error) {
	
	var reqBody *models.RouteRequest

	reqBody = &models.RouteRequest{
		Origin:        models.Address{Address: ORIGIN},
		Intermediates: []models.Address{{Address: currentAddress}, {Address: newAddress}},
		Destination:   models.Address{Address: ORIGIN},
		TravelMode:    "DRIVE",
	}
	return ds.doRouteRequest(reqBody)
}

func (ds *DistanceService) doRouteRequest(reqBody *models.RouteRequest) (*models.RouteDistances, error) {
	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		ds.logger.Errorf("Error marshalling request body: %v\n", err)
		return nil, err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequestWithContext(context.Background(), "POST", URL, bytes.NewReader(bodyBytes))
	if err != nil {
		ds.logger.Errorf("Error creating request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", ds.envs.MapsAPIKey)
	req.Header.Set("X-Goog-FieldMask", "routes.duration,routes.distanceMeters,routes.legs.distanceMeters,routes.legs.duration")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ds.logger.Errorf("Error sending request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ds.logger.Errorf("Error reading response body: %v\n", err)
		return nil, err
	}

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		ds.logger.Errorf("API request failed with status %d: %s\n", resp.StatusCode, string(data))
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Define the structure for the response
	var response *models.RouteResponse

	// Unmarshal the response body into the response structure
	if err := json.Unmarshal(data, &response); err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	fullDistanceMeters := response.Routes[0].DistanceMeters
	fullRouteKm := int(math.Round(float64(fullDistanceMeters) / 1000))

	toCurrentAddress := response.Routes[0].Legs[0].DistanceMeters
	toCurrentAddressKm := int(math.Round(float64(toCurrentAddress) / 1000))

	return &models.RouteDistances{
		ToCurrentAddress: toCurrentAddressKm,
		FullRoute:        fullRouteKm,
	}, nil
}
