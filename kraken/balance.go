package kraken

import (
	"fmt"

	"github.com/primexz/KrakenDCA/config"
)

func GetFiatBalance() (float64, error) {
	var balanceKey string
	if config.Currency == "AUD" {
		balanceKey = "Z"
	} else {
		balanceKey = config.FiatPrefix + config.Currency
	}

	return getKrakenBalance(balanceKey)
}

func getKrakenBalance(currency string) (float64, error) {
	balances, err := getApi().GetAccountBalances()
	if err != nil {
		return 0, err
	}

	balance, ok := balances[currency]
	if !ok {
		return 0, fmt.Errorf("no balance found for currency %s", currency)
	}

	ret, _ := balance.Float64()
	return ret, nil
}
