package propagationSample

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"com.aco.go.otel/initOtel"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

//

func Start() {

	initOtel.Engage("go-propagation")

	r := gin.Default()

	r.GET("/c01", c01)

	r.GET("/c02", c02)

	println(":: curl localhost:6060/c01")

	r.Run(":6060")
}

//

func c01(c *gin.Context) {

	println("c01 called")

	traceProvider := initOtel.TraceProvider

	tracer := traceProvider.Tracer("go-test-trace-parent")

	parentCtx, span := tracer.Start(c, "span-name-parent")

	defer span.End()

	time.Sleep(250 * time.Millisecond)
	println(
		fetchWithInjection("http://localhost:6060/c02", parentCtx),
	)
	time.Sleep(250 * time.Millisecond)

	c.JSON(200, gin.H{
		"message": "c01",
	})
}

//

func fetchWithInjection(url string, ctx context.Context) string {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// req.Header is basically map[string]string
	req.Header.Set("aco-t-1", "header---aco") // same
	req.Header.Add("aco-t-2", "header---yiu")

	// create context for injection
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	clt := http.Client{}

	resp, err := clt.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

//

func c02(c *gin.Context) {

	println("c02 called")

	// log http header
	for k, v := range c.Request.Header {
		fmt.Println(k, v)
	}

	traceProvider := initOtel.TraceProvider

	tracer := traceProvider.Tracer("go-test-tracer-child")

	// 通过http header，提取span元数据信息
	ctx := c.Request.Context()
	_, span := tracer.Start(
		otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)),
		"span-name-child",
	)

	defer span.End()

	time.Sleep(250 * time.Millisecond)

	c.JSON(200, gin.H{
		"message": "c02",
	})
}
