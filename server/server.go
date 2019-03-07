package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iodsp/user_center/context"
)

// Start the REST API server using the configuration provided
func Start(conf *context.Config) {
	if conf.HttpServerMode() != "" {
		gin.SetMode(conf.HttpServerMode())
	} else if conf.Debug() == false {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	registerRoutes(app, conf)

	app.Run(fmt.Sprintf("%s:%d", conf.HttpServerHost(), conf.HttpServerPort()))
}
