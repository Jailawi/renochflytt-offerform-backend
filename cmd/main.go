package main

import (
	"fmt"
	"os"
	"os/signal"
	"request-offer/pkg/server"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

type DBConfig struct {
    User     string
    Password string
    Host     string
    Port     int
}

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
			Name:     "host",
			Usage:    "Database host",
			Value:    "localhost",
		},
		&cli.Int64Flag{
			Name:     "HTTPport",
			Usage:    "http port",
			Value:    8080,
		},
		&cli.Int64Flag{
			Name:     "DBport",
			Usage:    "Database host",
			Value:    5432,
		},
		&cli.StringFlag{
			Name:     "user",
			Usage:    "Database user",
			EnvVars:  []string{"DB_USER"},
			Value:   "postgres",
		},
		&cli.StringFlag{
			Name:     "db-name",
			Usage:    "Database name",
			EnvVars:  []string{"DB_NAME"},
			Value:   "renochflytt",
		},
		&cli.StringFlag{
			Name:     "db-password",
			Usage:    "Database password",
			Value:   "admin",
			EnvVars:  []string{"DB_PASSWORD"},
		},
	}

	app.Action  = func(c *cli.Context) error {
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
	fmt.Println("Starting server...")

	server.Start(c, log)
	// port := c.Int64("HTTPport")
	// if err != nil {
	// 	log.Errorf("Failed to start server: %v", err)
	// }
	// db := database.Connect()

	// defer db.Close()

	
	// uuid := "059984fc-860b-4a66-a1e2-768cded2aa1d"
	// userType := "customer"
	// ssn := "123456789"
	// name := "John Doe"
	// email := "john@mail.com"
	// phone := "1234567890"
	

	//res, err := db.Exec("INSERT INTO public.customers (id, type, ssn, name, email, phone) VALUES ($1, $2, $3, $4, $5, $6)", uuid,  userType, ssn, name, email, phone)
	
	// res, err := db.Exec("SELECT * FROM public.customers")
	
	// if err != nil {
	// 	panic(err)
	// }

	// rowsAffected, err := res.RowsAffected()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Successfully connected! Rows affected: %d\n", rowsAffected)
}
