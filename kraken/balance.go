package kraken

import (
	"log"

	"github.com/primexz/KrakenDCA/config"
)

func GetFiatBalance() float64 {
	var balanceKey string
	if config.Currency == "AUD" {
		balanceKey = "Z"
	} else {
		balanceKey = config.FiatPrefix + config.Currency
	}

	return getKrakenBalance(balanceKey)
}

func GetBtcAmount() float64 {
	return getKrakenBalance("XXBT")
}

func getKrakenBalance(currency string) float64 {
	balances, err := getApi().GetAccountBalances()
	if err != nil {
		log.Fatalln("failed to query account balance")
	}

	balance, ok := balances[currency]
	if !ok {
		log.Fatalln("Failed to get balance for", currency)
	}

	ret, _ := balance.Float64()
	return ret
}
