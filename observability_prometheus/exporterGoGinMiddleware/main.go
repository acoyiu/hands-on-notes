package main

import (
	"prometheus/ginprometheus"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := gin.Default()
	r.Use(ginprometheus.Middleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/test", func(c *gin.Context) {})
	r.Run()
}
