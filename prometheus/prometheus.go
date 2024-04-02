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

	ExpectedOrderCnt = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "expected_orders_until_refill",
		Help: "Expected order count",
	})

	NextOrderTime = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "next_order_time",
		Help: "Time of next order",
	})

	TimeOfEmptyFiat = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "time_of_empty_fiat",
		Help: "Time of empty fiat",
	})
)

func StartServer() {
	log.Println("Starting Prometheus server..")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
