package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pack "k8s.client.app/pack"
)

func main() {
	fmt.Println(pack.Name)

	r := gin.Default()
	r.GET(
		"/metrics",
		func(c *gin.Context) {
			// do something
		},
		gin.WrapH(promhttp.Handler()),
	)
	r.Run()
}
