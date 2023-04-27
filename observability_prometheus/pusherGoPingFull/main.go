package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/samber/lo"
)

// ##
// # cmd <pushgateway-location> <instance_name> <job_name> [<url> <url> ...]
// #
// # To Debug: go run main.go http://localhost:9091 instance_name job_name http://123.57.136.251 https://google.com
// #
// # To Build: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
// #
// # To Run  : ./ping_metrics_to_gateway http://localhost:9091 instance_name job_name http://123.57.136.251 https://google.com
// ##

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func main() {

	if len(os.Args) < 5 {
		panic("There is no enough arguments entered")
	}

	pushTo = os.Args[1]
	instanceName = os.Args[2]
	jobName = os.Args[3]
	endpoints = os.Args[4:]

	if pushTo == "" {
		panic("There is no PushGateway endpoint entered")
	} else if instanceName == "" {
		panic("There is no instance name entered")
	} else if jobName == "" {
		panic("There is no job name entered")
	}

	fmt.Println("Now start fetching....")

	PingPusher()
}

var (
	pushTo       = ""
	instanceName = ""
	jobName      = ""
	endpoints    = []string{}

	PingTimeToSvc = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ping_time_to_service",
		Help: "Total time used to ping target service",
	}, []string{"endpoint_url"})
)

func PingPusher() {

	cList := lo.Map(endpoints, func(_ string, _ int) chan bool {
		return make(chan bool)
	})

	lo.ForEach(endpoints, func(ep string, i int) {

		go func(c chan bool) {
			client := http.Client{
				Timeout: 5 * time.Second,
			}

			sTime := time.Now().UnixMilli()

			fmt.Printf("Start: fetcing endpoint %s \n", ep)

			_, err := client.Get(ep)

			if err != nil {
				PingTimeToSvc.With(prometheus.Labels{"endpoint_url": ep}).Set(5000)
				fmt.Printf("Error: endpoint %s is not callable with 5000ms\n", ep)
			} else {
				eTime := time.Now().UnixMilli()
				delta := float64(eTime - sTime)
				PingTimeToSvc.With(prometheus.Labels{"endpoint_url": ep}).Set(delta)
				fmt.Printf("Succ: endpoint %s called with time %vms\n", ep, delta)
			}

			c <- true
			close(c)

		}(cList[i])
	})

	lo.ForEach(cList, func(finished chan bool, _ int) {
		<-finished
	})

	fmt.Printf("Start push to gateway '%s' with job name '%s'\n", pushTo, jobName)

	err := push.New(pushTo, jobName).
		Collector(PingTimeToSvc).
		Grouping("instance", instanceName).
		Grouping("usage", "monitoring").
		Push()

	if err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	} else {
		fmt.Println("Push successful and finished !")
	}
}
