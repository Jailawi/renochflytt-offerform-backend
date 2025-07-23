package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"request-offer/pkg/models"
	"request-offer/pkg/utils"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
}

// EmailService handles email operations
type EmailService struct {
	config *EmailConfig
	logger *logrus.Entry
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// NewEmailService creates a new email service
func NewEmailService(logger *logrus.Entry) *EmailService {
	emailConfig, err := LoadEmailConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load email configuration: %v", err))
	}
	return &EmailService{
		config: emailConfig,
		logger: logger,
	}
}

// LoadEmailConfig loads email configuration from environment variables
func LoadEmailConfig() (*EmailConfig, error) {

	config := &EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
		FromName:     os.Getenv("FROM_NAME"),
	}

	// Validate required fields
	if config.SMTPHost == "" || config.SMTPPassword == "" || config.FromEmail == "" {
		return nil, fmt.Errorf("missing required email configuration")
	}

	return config, nil
}

// SendEmail sends an email using the EmailService
func (s *EmailService) SendEmail(emailMsg *EmailMessage) error {
	s.logger.Infof("Sending email from: %s to: %v", s.config.FromEmail, emailMsg.To)

	// Build the complete email message
	msg := s.buildEmailMessage(emailMsg)

	// Setup authentication
	auth := smtp.PlainAuth("", s.config.FromEmail, s.config.SMTPPassword, s.config.SMTPHost)

	// Send email via SMTP
	serverAddr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	if err := smtp.SendMail(serverAddr, auth, s.config.FromEmail, emailMsg.To, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Infof("Email sent successfully!")
	return nil
}

// buildEmailMessage creates a properly formatted email message
func (s *EmailService) buildEmailMessage(emailMsg *EmailMessage) []byte {
	var builder strings.Builder

	// Generate unique message ID and timestamp
	messageID := fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), "gomail", "renochflytt.se")
	date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")

	// Write headers consistently using fmt.Fprintf
	fmt.Fprintf(&builder, "Message-ID: %s\r\n", messageID)
	fmt.Fprintf(&builder, "Date: %s\r\n", date)
	fmt.Fprintf(&builder, "From: \"%s\" <%s>\r\n", s.config.FromName, s.config.FromEmail)
	fmt.Fprintf(&builder, "To: %s\r\n", emailMsg.To[0])
	fmt.Fprintf(&builder, "Subject: %s\r\n", emailMsg.Subject)
	fmt.Fprintf(&builder, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&builder, "Content-Type: text/html; charset=UTF-8\r\n")
	fmt.Fprintf(&builder, "X-Mailer: Ren-Flytt-System\r\n")
	fmt.Fprintf(&builder, "\r\n") // Empty line separates headers from body
	builder.WriteString(emailMsg.Body)

	return []byte(builder.String())
}

// SendTestEmail sends a test email - convenience function for testing
func (s *EmailService) SendTestEmail(to []string, booking *models.Booking) error {
	tmpl, _ := template.New("email").Parse(utils.BookingEmailTemplate)
	var buf bytes.Buffer
	tmpl.Execute(&buf, booking)
	emailMsg := &EmailMessage{
		To:      to,
		Subject: "Important: Message from Ren & Flytt",
		Body:    buf.String(),
	}

	return s.SendEmail(emailMsg)
}


