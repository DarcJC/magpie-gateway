package logger

import "github.com/gin-gonic/gin"

type MagpieLogger interface {
	GetLoggerConfig() gin.LoggerConfig
}
