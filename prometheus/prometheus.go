package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MetricFiatAmount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fiat_amount",
		Help: "Current fiat amount",
	})
)

func StartServer() {
	log.Println("Starting Prometheus server..")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
