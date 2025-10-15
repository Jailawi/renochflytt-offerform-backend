package services

import (
	"fmt"
	"request-offer/pkg/models"

	"github.com/sirupsen/logrus"
)

const (
	PRICE_PACKNING                  = 1200 // SEK
	PRICE_MONTERING                 = 300  // SEK
	PRICE_PER_KM_OVER_THRESHOLD     = 8    // SEK
	PRICE_PER_KM_UNDER_THRESHOLD    = 20   // SEK
	PRICE_PER_KM_OVER_MAX_THRESHOLD = 23   // SEK
	RADIUS                          = 6    // km
	DISTANCE_THRESHOLD              = 60   // km
	DISTANCE_THRESHOLD_MAX          = 1000 // km
)

type PriceCalculator struct {
	logger *logrus.Entry
}

func NewPriceCalculator(logger *logrus.Entry) *PriceCalculator {
	return &PriceCalculator{
		logger: logger,
	}
}

func (pc *PriceCalculator) CalculatePrice(booking *models.Booking, route *models.RouteDistances) (int, error) {
	estimatedPrice := 0

	if route == nil {
		pc.logger.Warn("Distance information is nil, cannot calculate price")
		return estimatedPrice, fmt.Errorf("distance information is nil")
	}

	distanceCost := pc.calculateDistanceCost(route.FullRoute)
	pc.logger.Infof("Add distance cost: %d SEK, calculated distance: %d km", distanceCost, route.FullRoute)

	for _, service := range booking.Services {
		pc.logger.Infof("Calculating cost for service: %s", service)
		switch service {
		case "Flytthjälp":
			movingCost := pc.calculateMovingResidenceCost(*booking.CurrentResidence.LivingArea)
			pc.logger.Infof("Add moving cost: %d SEK for living area: %d kvm", movingCost, *booking.CurrentResidence.LivingArea)
			accessCost := pc.calculateAccessCost(booking.CurrentResidence.Accessibility, *booking.CurrentResidence.Floor)
			pc.logger.Infof("Add cost: %d SEK for accessibility: %s, floor: %d", accessCost, booking.CurrentResidence.Accessibility, *booking.CurrentResidence.Floor)
			estimatedPrice += distanceCost + movingCost + accessCost
			pc.logger.Infof("subtotal: %d SEK", estimatedPrice)
		case "Flyttstädning":
			if route.ToCurrentAddress <= 120 {
				cleaningCost := pc.calculateCleaningCost(*booking.CurrentResidence.LivingArea)
				cleaningDistanceCost := int(float64(pc.calculateDistanceCost(route.ToCurrentAddress)) * 0.5)
				pc.logger.Infof("Add cleaning cost: %d SEK for living area: %d kvm and distance cost: %d SEK", cleaningCost, *booking.CurrentResidence.LivingArea, cleaningDistanceCost)
				estimatedPrice += cleaningCost + cleaningDistanceCost
			} else {
				// 30kr/kvm
				cleaningCost := (*booking.CurrentResidence.LivingArea) * 30
				estimatedPrice += cleaningCost
				pc.logger.Infof("Add cleaning cost without distance cost: %d SEK for living area: %d kvm", cleaningCost, *booking.CurrentResidence.LivingArea)
			}
			pc.logger.Infof("subtotal: %d SEK", estimatedPrice)
		case "Packning":
			estimatedPrice += PRICE_PACKNING
			pc.logger.Infof("Add packing cost: %d SEK", PRICE_PACKNING)
		case "Montering":
			estimatedPrice += PRICE_MONTERING
			pc.logger.Infof("Add assembly cost: %d SEK", PRICE_MONTERING)
		}
	}
	return estimatedPrice, nil
}

func (pc *PriceCalculator) calculateDistanceCost(distanceKm int) int {
	if distanceKm <= RADIUS*3 {
		return 0
	}

	if distanceKm <= DISTANCE_THRESHOLD {
		return distanceKm * PRICE_PER_KM_UNDER_THRESHOLD
	}

	if distanceKm <= DISTANCE_THRESHOLD_MAX {
		return ((distanceKm - DISTANCE_THRESHOLD) * PRICE_PER_KM_OVER_THRESHOLD) + (DISTANCE_THRESHOLD * PRICE_PER_KM_UNDER_THRESHOLD)
	}

	return ((distanceKm - DISTANCE_THRESHOLD - DISTANCE_THRESHOLD_MAX) * PRICE_PER_KM_OVER_THRESHOLD) + (DISTANCE_THRESHOLD * PRICE_PER_KM_UNDER_THRESHOLD) + (PRICE_PER_KM_OVER_MAX_THRESHOLD * DISTANCE_THRESHOLD_MAX)
}

func (pc *PriceCalculator) calculateMovingResidenceCost(livingArea int) int {
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

func (pc *PriceCalculator) calculateCleaningCost(livingArea int) int {
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
func (pc *PriceCalculator) calculateAccessCost(access string, floor int) int {
	switch access {
	case "Trappor", "Liten Hiss":
		switch {
		case floor <= 1:
			return 0
		default:
			return floor * 200
		}
	default:
		return 0
	}
}
