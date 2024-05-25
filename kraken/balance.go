package kraken

import (
	"fmt"

	"github.com/primexz/KrakenDCA/config"
)

func (k *KrakenApi) GetFiatBalance() (float64, error) {
	var balanceKey string
	if config.C.Currency == "AUD" {
		balanceKey = "Z"
	} else {
		balanceKey = config.C.FiatPrefix + config.C.Currency
	}

	balances, err := k.api.GetAccountBalances()
	if err != nil {
		return 0, err
	}

	balance, ok := balances[balanceKey]
	if !ok {
		return 0, fmt.Errorf("no balance found for currency %s", balanceKey)
	}

	ret, _ := balance.Float64()
	return ret, nil
}
