package config

import (
	"os"
	"strconv"

	"github.com/primexz/KrakenDCA/logger"
)

var (
	KrakenPublicKey      string
	KrakenPrivateKey     string
	Currency             string
	KrakenOrderSize      float64
	CheckDelay           float64
	CryptoPrefix         string
	FiatPrefix           string
	ExperimentalMakerFee bool
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("config")
}

func LoadConfiguration() {
	log.Info("Loading configuration..")

	KrakenPublicKey = loadRequiredEnvVariable("KRAKEN_PUBLIC_KEY")
	KrakenPrivateKey = loadRequiredEnvVariable("KRAKEN_PRIVATE_KEY")
	Currency = loadFallbackEnvVariable("CURRENCY", "USD")
	KrakenOrderSize = loadFloatEnvVariableWithFallback("KRAKEN_ORDER_SIZE", 0.0001) // https://support.kraken.com/hc/en-us/articles/205893708-Minimum-order-size-volume-for-trading
	CheckDelay = loadFloatEnvVariableWithFallback("CHECK_DELAY", 60)
	ExperimentalMakerFee = loadBoolEnvVariableWithFallback("EXPERIMENTAL_MAKER_FEE", false)

	if Currency == "USD" || Currency == "EUR" || Currency == "GBP" {
		CryptoPrefix = "X"
		FiatPrefix = "Z"
	}
}

func loadRequiredEnvVariable(envVar string) string {
	envData := os.Getenv(envVar)

	if envData == "" {
		log.Fatal("Required environment variable", envVar, "missing.")
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
		log.Fatal("Failed to parse float environment variable.", envVar)
	}

	return fallback
}

func loadBoolEnvVariableWithFallback(envVar string, fallback bool) bool {
	envData := os.Getenv(envVar)

	if envData == "" {
	} else if s, err := strconv.ParseBool(envData); err == nil {
		return s
	} else {
		log.Fatal("Failed to parse bool environment variable.", envVar)
	}

	return fallback
}
