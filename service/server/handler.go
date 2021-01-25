package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"magpie-gateway/configuration"
	"magpie-gateway/logger"
	"net/http"
	"time"
)

func ServiceHandler() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	if configuration.GlobalConfiguration.Debug {
		e.Use(gin.Logger())
	} else {
		fileLogger := logger.NewFileLogger(fmt.Sprintf("logs/magpie-services-%s.log", time.Now().Format("2006-01-02-15-04-05")))
		lc := fileLogger.GetLoggerConfig()
		e.Use(gin.LoggerWithConfig(lc))
		gin.DefaultWriter = lc.Output
	}

	return e
}

