package server

import (
	"fmt"
	"net/http"
	"request-offer/pkg/handlers"
	"request-offer/pkg/services"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Start(c *cli.Context, log *logrus.Entry, bookingService *services.BookingService) {
	mux := http.NewServeMux()

	handler := handlers.NewHandler(bookingService, log)

	handler.RegisterRoutes(mux)

	port := fmt.Sprintf(":%d", c.Int64("HTTPport"))
	log.Infof("Starting server on port: %s", port)
	http.ListenAndServe(port, mux)

}
