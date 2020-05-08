package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type app struct{}

// NewApp is a constructor for application server
func NewApp() *app {
	return &app{}
}

func (a app) Run() {
	router := initializeRouter()
	logrus.Infof("Server started, listening on %v port", servicePort())

	err := fasthttp.ListenAndServe(":"+servicePort(), router.Handler)
	logrus.Errorf("Could not run server: %v", err)
}

func servicePort() string {
	port := os.Getenv("HTTP_PORT")
	if len(port) == 0 {
		return "80"
	}

	return port
}
