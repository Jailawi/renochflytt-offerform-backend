package main

import (
	"fmt"
	"os"
	"os/signal"
	"request-offer/database"
	"request-offer/pkg/server"
	"request-offer/pkg/services"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
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
		&cli.StringFlag{
			Name:  "host",
			Usage: "Database host",
			Value: "localhost",
		},
		&cli.Int64Flag{
			Name:  "HTTPport",
			Usage: "http port",
			Value: 8080,
		},
		&cli.Int64Flag{
			Name:  "DBport",
			Usage: "Database host",
			Value: 5432,
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "Database user",
			EnvVars: []string{"DB_USER"},
			Value:   "postgres",
		},
		&cli.StringFlag{
			Name:    "db-name",
			Usage:   "Database name",
			EnvVars: []string{"DB_NAME"},
			Value:   "renochflytt",
		},
		&cli.StringFlag{
			Name:    "db-password",
			Usage:   "Database password",
			Value:   "admin",
			EnvVars: []string{"DB_PASSWORD"},
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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Warnf("Error loading .env file: %v", err)
	}

	// Initialize database connection
	db, err := database.ConnectToMongoDB(log)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	emailService := services.NewEmailService(log)

	bookingService := services.NewBookingService(db, emailService, log)

	// Start the server
	server.Start(c, log, bookingService)
	log.Infof("Application started successfully")
}
