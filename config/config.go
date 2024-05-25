package config

import (
	"github.com/caarlos0/env/v11"
	log "github.com/sirupsen/logrus"
)

type config struct {
	KrakenPublicKey      string  `env:"KRAKEN_PUBLIC_KEY,required"`
	KrakenPrivateKey     string  `env:"KRAKEN_PRIVATE_KEY,required"`
	Currency             string  `env:"CURRENCY" envDefault:"USD"`
	KrakenOrderSize      float64 `env:"KRAKEN_ORDER_SIZE" envDefault:"0.0001"`
	CheckDelay           float64 `env:"CHECK_DELAY" envDefault:"60"`
	LimitOrderMode       bool    `env:"LIMIT_ORDER_MODE" envDefault:"false"`
	LimitOrderRetryCount int     `env:"LIMIT_ORDER_RETRY_COUNT" envDefault:"8"`
	MetricPort           int     `env:"METRIC_PORT" envDefault:"3000"`

	CryptoPrefix string
	FiatPrefix   string
}

var (
	logger = log.WithFields(log.Fields{
		"prefix": "bot",
	})
	C config
)

func init() {
	loadConfiguration()
}

func loadConfiguration() {
	if config, err := env.ParseAs[config](); err == nil {
		C = config
	} else {
		logger.Fatal(err)
	}

	if C.Currency == "USD" || C.Currency == "EUR" || C.Currency == "GBP" {
		C.CryptoPrefix = "X"
		C.FiatPrefix = "Z"
	}
}
