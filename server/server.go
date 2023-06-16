package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"ticker-server/config"
	"ticker-server/handler"
	"time"
)

func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		glog.Info(fmt.Sprintf("| %3d | %13v | %15s | %s | %s |", statusCode, latencyTime, clientIP, reqMethod, reqUri))
	}
}

func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(LoggerToFile())
	apiPrefix := "/"
	g := r.Group(apiPrefix)

	coinHandler := handler.NewCoinHandler()

	g.GET("/price/:currency", coinHandler.GetPrice)

	r.Run(":" + config.Config().Server.Port)
}
