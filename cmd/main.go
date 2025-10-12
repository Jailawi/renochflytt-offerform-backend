package main

import (
	"fmt"
	"os"
	"os/signal"
	"request-offer/database"
	"request-offer/pkg/middleware"
	"request-offer/pkg/models"
	"request-offer/pkg/server"
	"request-offer/pkg/services"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

const (
	HTTPPort      = "http-port"
	Origins       = "origins"
	MongoUri      = "mongo-uri"
	FromEmail     = "from-email"
	FromName      = "from-name"
	APIKey        = "api-key"
	MapsAPIKey    = "maps-api-key"
	MailgunAPIKey = "mailgun-api-key"
)

func main() {
	app := createApp()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func createApp() *cli.App {
	fmt.Println("Creating app...")
	app := cli.NewApp()
	app.Name = "Booking App"
	app.Usage = "A booking app for moving services"
	app.Flags = []cli.Flag{
		&cli.Int64Flag{
			Name:    HTTPPort,
			Usage:   "HTTP port",
			Value:   8080,
			EnvVars: []string{"PORT"},
		},
		&cli.StringFlag{
			Name:    Origins,
			Usage:   "Allowed origins",
			Value:   "*",
			EnvVars: []string{"ORIGINS"},
		},
		&cli.StringFlag{
			Name:    MongoUri,
			Usage:   "MongoDB URI",
			EnvVars: []string{"MONGO_URI"},
		},
		&cli.StringFlag{
			Name:    FromEmail,
			Usage:   "From email address",
			Value:   "info@renochflytt.se",
			EnvVars: []string{"FROM_EMAIL"},
		},
		&cli.StringFlag{
			Name:    FromName,
			Usage:   "From name",
			Value:   "Ren & Flytt",
			EnvVars: []string{"FROM_NAME"},
		},
		&cli.StringFlag{
			Name:    APIKey,
			Usage:   "API key",
			EnvVars: []string{"X-API-KEY"},
		},
		&cli.StringFlag{
			Name:    MapsAPIKey,
			Usage:   "Maps API key",
			EnvVars: []string{"MAPS_API_KEY"},
		},
		&cli.StringFlag{
			Name:    MailgunAPIKey,
			Usage:   "Mailgun API key",
			EnvVars: []string{"MAILGUN_API_KEY"},
		},
	}

	app.Action = func(c *cli.Context) error {
		log := logrus.WithFields(logrus.Fields{})
		done := createTerminationHandler(log)
		go start(c, log)
		<-done
		return nil
	}

	return app
}

func createTerminationHandler(log *logrus.Entry) chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Warnf("Received signal: %s. Shutting down...", sig)
		done <- true
	}()
	return done
}

func start(c *cli.Context, log *logrus.Entry) {
	log.Infof("Starting application...")

	envs := &models.Envs{
		Origins :      c.String(Origins),
		MongoUri:      c.String(MongoUri),
		FromEmail:     c.String(FromEmail),
		FromName:      c.String(FromName),
		APIKey:        c.String(APIKey),
		MapsAPIKey:    c.String(MapsAPIKey),
		MailgunAPIKey: c.String(MailgunAPIKey),
	}

	// Initialize database connection
	db, err := database.ConnectToMongoDB(envs.MongoUri, log)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	emailService := services.NewEmailService(db, envs, log)

	bookingService := services.NewBookingService(db, emailService, log, envs)

	// Start cleanup of inactive limiters
	middleware.CleanupInactiveLimiters()

	// Start the server
	server.Start(c, log, envs, bookingService)
	log.Infof("Application started successfully")
}
