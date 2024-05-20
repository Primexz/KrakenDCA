package metrics

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/primexz/KrakenDCA/config"
	log "github.com/sirupsen/logrus"
)

type metric struct {
	NextOrder int64 `json:"nextOrder"`
}

var (
	Metrics = metric{}

	logger = log.WithFields(log.Fields{
		"prefix": "prometheus",
	})
)

func StartServer() {
	go func() {
		port := config.MetricPort

		logger.WithField("port", port).Info("Starting Prometheus metrics server started.")

		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(Metrics); err != nil {
				logger.WithError(err).Error("Failed to encode metrics.")
			}
		})

		server := &http.Server{
			Addr:              ":" + strconv.Itoa(port),
			ReadHeaderTimeout: 3 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
			logger.WithError(err).Error("Failed to start metrics server.")
		}
	}()
}
