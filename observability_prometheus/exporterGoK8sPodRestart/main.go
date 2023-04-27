package main

import (
	"prometheus/restartCount"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := gin.Default()
	r.GET(
		"/metrics",
		func(c *gin.Context) { restartCount.StartLoop() },
		gin.WrapH(promhttp.Handler()),
	)
	r.Run()
}
