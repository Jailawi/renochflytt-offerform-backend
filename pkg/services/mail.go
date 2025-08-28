package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"request-offer/pkg/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
}

// EmailService handles email operations
type EmailService struct {
	envs       *models.Envs
	logger     *logrus.Entry
	mongoClient *mongo.Client
}

// NewEmailService creates a new email service
func NewEmailService(mongoClient *mongo.Client, envs *models.Envs, logger *logrus.Entry) *EmailService {
	return &EmailService{
		envs:       envs,
		logger:     logger,
		mongoClient: mongoClient,
	}
}


// SendEmail sends an email using the EmailService
func (s *EmailService) SendEmail(emailMsg *EmailMessage, bookingID primitive.ObjectID) error {
	s.logger.Infof("Sending email from: %s to: %v", s.envs.FromEmail, emailMsg.To)

	// Build the complete email message
	msg := s.buildEmailMessage(emailMsg)

	// Setup authentication
	auth := smtp.PlainAuth("", s.envs.FromEmail, s.envs.SMTPPass, s.envs.SMTPHost)

	// Send email via SMTP
	serverAddr := fmt.Sprintf("%s:%s", s.envs.SMTPHost, s.envs.SMTPPort)
	if err := smtp.SendMail(serverAddr, auth, s.envs.FromEmail, emailMsg.To, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.mongoClient.Database("renochflytt").Collection("bookings").UpdateOne(
		context.Background(),
		bson.M{"_id": bookingID},
		bson.M{"$set": bson.M{"email_sent": true}},
	)

	s.logger.Infof("Email sent successfully!")
	return nil
}

// buildEmailMessage creates a properly formatted HTML email message
// buildEmailMessage creates a properly formatted HTML email message
func (s *EmailService) buildEmailMessage(emailMsg *EmailMessage) []byte {
	var buf bytes.Buffer

	// Generate unique message ID and timestamp
	messageID := fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), "gomail", "renochflytt.se")
	date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")

	// Write headers for simple HTML email
	fmt.Fprintf(&buf, "Message-ID: %s\r\n", messageID)
	fmt.Fprintf(&buf, "Date: %s\r\n", date)
	fmt.Fprintf(&buf, "From: \"%s\" <%s>\r\n", s.envs.FromName, s.envs.FromEmail)
	fmt.Fprintf(&buf, "To: %s\r\n", emailMsg.To[0])
	fmt.Fprintf(&buf, "Subject: %s\r\n", emailMsg.Subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: text/html; charset=UTF-8\r\n")
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: 8bit\r\n")
	fmt.Fprintf(&buf, "X-Mailer: Ren-Flytt-System\r\n")
	fmt.Fprintf(&buf, "\r\n")

	// Write HTML body
	buf.WriteString(emailMsg.Body)

	return buf.Bytes()
}

// SendBookingConfirmationEmail sends a booking confirmation email
func (s *EmailService) SendBookingConfirmationEmail(to []string, booking *models.Booking, bookingID primitive.ObjectID) error {
	// Load template from file
	templatePath := "templates/customer-booking.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load email template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, booking)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	emailMsg := &EmailMessage{
		To:      to,
		Subject: "Bokningsbekräftelse - Ren & Flytt",
		Body:    buf.String(),
	}

	return s.SendEmail(emailMsg, bookingID)
}

// SendTestEmail sends a test email - convenience function for testing
func (s *EmailService) SendTestEmail(to []string, booking *models.Booking) error {
	// Load template from file
	templatePath := "../templates/customer-booking.html"
	companyEmail := s.envs.FromEmail
	receivers := append(to, companyEmail)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load email template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, booking)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	emailMsg := &EmailMessage{
		To:      receivers,
		Subject: "Bokningsbekräftelse - Ren & Flytt",
		Body:    buf.String(),
	}

	return s.SendEmail(emailMsg, booking.ID)
}
