package config

import (
	"log"
	"os"
	"strconv"
)

var (
	KrakenPublicKey  string
	KrakenPrivateKey string
	Currency         string
	KrakenOrderSize  float64
	CheckDelay       float64
	CryptoPrefix     string
	FiatPrefix       string
)

func LoadConfiguration() {
	log.Println("Loading configuration..")

	KrakenPublicKey = loadRequiredEnvVariable("KRAKEN_PUBLIC_KEY")
	KrakenPrivateKey = loadRequiredEnvVariable("KRAKEN_PRIVATE_KEY")
	Currency = loadFallbackEnvVariable("❌", "USD")
	KrakenOrderSize = loadFloatEnvVariableWithFallback("KRAKEN_ORDER_SIZE", 0.0001) // https://support.kraken.com/hc/en-us/articles/205893708-Minimum-order-size-volume-for-trading
	CheckDelay = loadFloatEnvVariableWithFallback("CHECK_DELAY", 60)

	if Currency == "USD" || Currency == "EUR" || Currency == "GBP" {
		CryptoPrefix = "X"
		FiatPrefix = "Z"
	}
}

func loadRequiredEnvVariable(envVar string) string {
	envData := os.Getenv(envVar)

	if envData == "" {
		log.Fatalln("Required environment variable", envVar, "missing.")
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
		log.Fatalln("Failed to parse float environment variable.")
	}

	return fallback
}
