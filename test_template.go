package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"request-offer/pkg/models"
	"time"
)

func main() {
	// Create sample booking data
	sampleBooking := createSampleBooking()

	// Parse and execute the template from file
	tmpl, err := template.ParseFiles("templates/customer-booking.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, sampleBooking)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	// Write the rendered HTML to a file
	outputFile := "rendered_email.html"
	err = os.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("✅ Email template rendered successfully!\n")
	fmt.Printf("📁 Output saved to: %s\n", outputFile)
	fmt.Printf("🌐 Open the file in a web browser to see how it looks\n")
}

func createSampleBooking() *models.Booking {
	// Create sample data
	livingArea := 75
	floor := 2
	message := "Behöver extra hjälp med piano och stora möbler"

	movingDate := models.Date{Time: time.Now().AddDate(0, 0, 14)}   // 2 weeks from now
	cleaningDate := models.Date{Time: time.Now().AddDate(0, 0, 16)} // 16 days from now

	return &models.Booking{
		Services: []string{
			"Flytthjälp",
			"Städning",
			"Packningshjälp",
			"Möbelmontering",
		},
		MovingDate:     &movingDate,
		IsFlexibleDate: true,
		CleaningDate:   &cleaningDate,
		CurrentResidence: models.Residence{
			Address:       "Storgatan 15, 11455 Stockholm",
			ResidenceType: "Lägenhet",
			LivingArea:    &livingArea,
			Accessibility: "Hiss finns",
			Floor:         &floor,
		},
		NewResidence: models.Residence{
			Address:       "Vasagatan 22, 11120 Stockholm",
			ResidenceType: "Lägenhet",
			LivingArea:    &livingArea,
			Accessibility: "Trappor, 3 våningar",
			Floor:         &floor,
		},
		Contact: models.Contact{
			Name:      "Anna Andersson",
			SSN:       "19851205-1234",
			Email:     "anna.andersson@email.se",
			Phone:     "070-123 45 67",
			Rutavdrag: true,
			Message:   &message,
			Consent:   true,
		},
		EmailSent: false,
		CreatedAt: time.Now(),
	}
}
