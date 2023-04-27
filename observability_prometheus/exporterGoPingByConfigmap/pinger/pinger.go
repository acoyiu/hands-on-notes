package pinger

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/lo"
)

var (
	endpoints = []string{
		"http://123.57.136.251",
		"https://api.funacademycn.com",
		"https://highlights.milkcargo.cn",
	}
	wg            = sync.WaitGroup{}
	PingTimeToSvc = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ping_time_to_service",
		Help: "Total time used to ping target service",
	}, []string{"endpoint_url"})
)

func init() {
	prometheus.Register(PingTimeToSvc)
}

func StartPing() {

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
}
