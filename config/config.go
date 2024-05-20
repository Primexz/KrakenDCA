package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var (
	KrakenPublicKey  string
	KrakenPrivateKey string
	Currency         string
	KrakenOrderSize  float64
	CheckDelay       float64
	CryptoPrefix     string
	FiatPrefix       string

	LimitOrderMode       bool
	LimitOrderRetryCount int

	MetricPort int
)

var logger = log.WithFields(log.Fields{
	"prefix": "bot",
})

func LoadConfiguration() {
	logger.Info("Loading configuration..")

	KrakenPublicKey = loadRequiredEnvVariable("KRAKEN_PUBLIC_KEY")
	KrakenPrivateKey = loadRequiredEnvVariable("KRAKEN_PRIVATE_KEY")
	Currency = loadFallbackEnvVariable("CURRENCY", "USD")
	KrakenOrderSize = loadFloatEnvVariableWithFallback("KRAKEN_ORDER_SIZE", 0.0001) // https://support.kraken.com/hc/en-us/articles/205893708-Minimum-order-size-volume-for-trading
	CheckDelay = loadFloatEnvVariableWithFallback("CHECK_DELAY", 60)
	LimitOrderMode = loadBoolEnvVariableWithFallback("LIMIT_ORDER_MODE", false)
	LimitOrderRetryCount = int(loadFloatEnvVariableWithFallback("LIMIT_ORDER_RETRY_COUNT", 8))

	MetricPort = int(loadFloatEnvVariableWithFallback("METRIC_PORT", 3000))

	if Currency == "USD" || Currency == "EUR" || Currency == "GBP" {
		CryptoPrefix = "X"
		FiatPrefix = "Z"
	}
}

func loadRequiredEnvVariable(envVar string) string {
	envData := os.Getenv(envVar)

	if envData == "" {
		logger.WithField("var", envVar).Fatal("Required environment variable missing.")
	}

	return envData
}

func loadFallbackEnvVariable(envVar string, fallback string) string {
	envData := os.Getenv(envVar)

	if envData == "" {
		envData = fallback
	}

	return envData
}

func loadFloatEnvVariableWithFallback(envVar string, fallback float64) float64 {
	envData := os.Getenv(envVar)

	if envData == "" {
	} else if s, err := strconv.ParseFloat(envData, 32); err == nil {
		return s
	} else {
		logger.WithField("var", envVar).Fatal("Failed to parse float environment variable.")
	}

	return fallback
}

func loadBoolEnvVariableWithFallback(envVar string, fallback bool) bool {
	envData := os.Getenv(envVar)

	if envData == "" {
	} else if s, err := strconv.ParseBool(envData); err == nil {
		return s
	} else {
		logger.WithFields(log.Fields{
			"var": envVar,
		}).Error("Failed to parse bool environment variable.")

		logger.Fatal("Failed to parse bool environment variable.", envVar)
	}

	return fallback
}
