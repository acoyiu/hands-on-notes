package restartCount

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

var (
	RestartTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pod_restart_total_count",
		Help: "Total number of restart count of application pods",
	}, []string{"namespace", "pod_name", "under_deployment"})
)

func init() {
	// prometheus.MustRegister(RestartTotal)
	prometheus.Register(RestartTotal)

}

func StartLoop() {
	loopNamespace("fa")
	loopNamespace("hw")
	loopNamespace("mc2")
}

func loopNamespace(namespaceName string) {

	caCert, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	req, err := http.NewRequest("GET", "https://kubernetes.default/api/v1/namespaces/"+namespaceName+"/pods", nil)
	if err != nil {
		log.Fatal(err)
	}

	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("%v%v", "Bearer ", string(token)),
	)

	result, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	jsonString := string(body)

	podList := gjson.Get(jsonString, "items")
	podList.ForEach(func(key, podValue gjson.Result) bool {

		RestartTotal.With(prometheus.Labels{
			"namespace":        namespaceName,
			"pod_name":         podValue.Map()["metadata"].Map()["name"].String(),
			"under_deployment": podValue.Map()["metadata"].Map()["ownerReferences"].Array()[0].Map()["name"].String(),
		}).Set(podValue.Map()["status"].Map()["containerStatuses"].Array()[0].Map()["restartCount"].Float())

		return true
	})
}
