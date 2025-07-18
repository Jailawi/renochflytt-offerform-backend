package server

import (
	"fmt"
	"net/http"
	"request-offer/pkg/handlers"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Start(c *cli.Context, log *logrus.Entry) {
	mux := http.NewServeMux()

	handlers.RegisterRoutes(mux)

	port := fmt.Sprintf(":%d", c.Int64("HTTPport"))
	log.Infof("Starting server on port: %s", port)
	http.ListenAndServe(port, mux)

}