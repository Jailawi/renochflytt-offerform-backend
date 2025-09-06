package server

import (
	"fmt"
	"net/http"
	"request-offer/pkg/handlers"
	"request-offer/pkg/middleware"
	"request-offer/pkg/models"
	"request-offer/pkg/services"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Start(c *cli.Context, log *logrus.Entry, envs *models.Envs, bookingService *services.BookingService) {
	mux := http.NewServeMux()
	
	middlewares := chain(mux, middleware.ApikeyMiddleware(envs.APIKey), middleware.RateLimiterMiddleware)
	
	handler := handlers.NewHandler(bookingService, log)
	handler.RegisterRoutes(mux)

	port := fmt.Sprintf(":%d", c.Int64("http-port"))
	log.Infof("Starting server on port: %s", port)
	http.ListenAndServe(port, middlewares)

}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
