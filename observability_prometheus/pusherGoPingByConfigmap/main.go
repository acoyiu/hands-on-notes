package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/samber/lo"
)

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// Create Metrics
var completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "db_backup_last_completion_timestamp_seconds",
	Help: "The timestamp of the last successful completion of a DB backup.",
})

func GaugePusher() {

	// set metrics to currenttime, set可以设置任意值（float64）, e.g.:completionTime.Set(200)
	completionTime.SetToCurrentTime()

	// Grouping label, 不是 metric‘s label，给指标添加标签，可以添加多个
	// instance 係必定會有的字段，default 係 empty String
	err := push.New("http://pushgateway-service:9091", "job_name").
		Collector(completionTime).
		Grouping("db", "customers").
		Grouping("instance", "1.1.1.1").
		Push()

	if err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	} else {
		fmt.Println("Push success!")
	}
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

var (
	pushTo    = "http://pushgateway-service:9091"
	jobName   = "job_name"
	endpoints = []string{
		"http://123.57.136.251",
		"https://api.funacademycn.com",
		"https://highlights.milkcargo.cn",
	}
	PingTimeToSvc = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ping_time_to_service",
		Help: "Total time used to ping target service",
	}, []string{"endpoint_url"})
)

var wg = sync.WaitGroup{}

func PingPusher() {

	wg.Add(len(endpoints))

	lo.ForEach(endpoints, func(ep string, _ int) {

		go func() {
			client := http.Client{
				Timeout: 5 * time.Second,
			}

			sTime := time.Now().UnixMilli()

			_, err := client.Get(ep)

			if err != nil {
				PingTimeToSvc.
					With(prometheus.Labels{"endpoint_url": ep}).
					Set(5000)
			} else {
				eTime := time.Now().UnixMilli()

				PingTimeToSvc.
					With(prometheus.Labels{"endpoint_url": ep}).
					Set(float64(eTime - sTime))
			}

			wg.Done()
		}()
	})

	wg.Wait()

	err := push.New(pushTo, jobName).
		Collector(PingTimeToSvc).
		Grouping("instance", "tester").
		Push()

	if err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	} else {
		fmt.Println("Push success!")
	}
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func main() {
	// GaugePusher()
	PingPusher()
}
