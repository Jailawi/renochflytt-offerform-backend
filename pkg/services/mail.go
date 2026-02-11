package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"request-offer/pkg/models"
	"request-offer/pkg/util"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mailgun/mailgun-go/v4"
)

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
}

// EmailService handles email operations
type EmailService struct {
	envs        *models.Envs
	logger      *logrus.Entry
	mongoClient *mongo.Client
}

// NewEmailService creates a new email service
func NewEmailService(mongoClient *mongo.Client, envs *models.Envs, logger *logrus.Entry) *EmailService {
	return &EmailService{
		envs:        envs,
		logger:      logger,
		mongoClient: mongoClient,
	}
}

// SendEmail sends an email using the EmailService
func (s *EmailService) SendEmail(emailMsg *EmailMessage, bookingID primitive.ObjectID) error {
	s.logger.Infof("Sending email from: %s to: %v", s.envs.FromEmail, emailMsg.To)

	// Build the complete email message
	_, err := s.SendWithMailGun("renochflytt.se", emailMsg) // --- IGNORE ---
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.mongoClient.Database(s.envs.Database).Collection("bookings").UpdateOne(
		context.Background(),
		bson.M{"_id": bookingID},
		bson.M{"$set": bson.M{"email_sent": true}},
	)

	s.logger.Infof("Email sent successfully!")
	return nil
}

// SendTestEmail sends a test email - convenience function for testing
func (s *EmailService) SendTestEmail(to []string, booking *models.Booking) error {
	s.logger.Infof("Loading email template and sending test email to: %v", to)
	// Load template from file
	templatePath := "templates/company-booking.html"
	companyEmail := s.envs.FromEmail
	receivers := append(to, companyEmail)
	tmpl := template.New("company-booking.html").Funcs(template.FuncMap{
		"formatSEK": util.FormatSEK,
	})
	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load email template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, booking)
	if err != nil {
		s.logger.Errorf("Error executing template: %v", err)
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	emailMsg := &EmailMessage{
		To:      receivers,
		Subject: "Bokningsbekräftelse - Ren & Flytt",
		Body:    buf.String(),
	}

	return s.SendEmail(emailMsg, booking.ID)
}

func (s *EmailService) SendWithMailGun(domain string, emailMsg *EmailMessage) (string, error) {
	s.logger.Infof("Sending email with MailGun to: %v", emailMsg.To)
	mg := mailgun.NewMailgun(domain, s.envs.MailgunAPIKey)
	//When you have an EU-domain, you must specify the endpoint:
	mg.SetAPIBase("https://api.eu.mailgun.net/v3")
	msg := mailgun.NewMessage(
		fmt.Sprintf("Ren & Flytt <%s>", s.envs.FromEmail),
		emailMsg.Subject,
		"",
		strings.Join(emailMsg.To, ", "),
	)

	msg.SetHTML(emailMsg.Body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, msg)
	s.logger.Infof("Sent email with MailGun to: %v", emailMsg.To)
	return id, err
}
