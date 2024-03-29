package bot

import (
	"log"
	"time"

	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/kraken"
	"github.com/primexz/KrakenDCA/prometheus"
)

var (
	timeOfNextOrder  time.Time
	timeOfEmptyFiat  time.Time
	lastFiatBalance  float64
	lastBtcFiatPrice float64
	fiatAmount       float64
	initialRun       bool = true
)

func StartBot() {
	log.Println("Starting Kraken DCA Bot lifecycle..")

	for {
		run()

		initialRun = false
		time.Sleep(time.Duration(config.CheckDelay) * time.Second)
	}

}

func run() {
	// fiatBalance := kraken.GetFiatBalance()
	kraken.GetFiatBalance()
	fiatAmount = 1000
	prometheus.MetricFiatAmount.Set(fiatAmount)

	if fiatAmount == 0 {
		log.Println("No remaining fiat balance found. It's time to top up your account ;)")
		return
	}

	newFiatMoney := fiatAmount > lastFiatBalance
	if newFiatMoney || initialRun {
		timeOfEmptyFiat = computeNextFiatDepositDay()

		log.Println("Next Fiat deposit required at", timeOfEmptyFiat)

		updateFiatBalance(fiatAmount)
	}

	lastBtcFiatPrice = kraken.GetCurrentBtcFiatPrice()

	if timeOfNextOrder.Before(time.Now()) || newFiatMoney {
		log.Println("Placing order..")
		//kraken.BuyBtc()

		calculateTimeOfNextOrder()
	}

	log.Println("Next order in", fmtDuration(time.Until(timeOfNextOrder)), timeOfNextOrder)

	//TODO: ADD AUTO WITHDRAWAL
}

func calculateTimeOfNextOrder() {
	fiatValueInBtc := fiatAmount / lastBtcFiatPrice
	orderAmountUntilRefill := fiatValueInBtc / config.KrakenOrderSize

	now := time.Now().UnixMilli()
	timeOfNextOrder = time.UnixMilli((timeOfEmptyFiat.UnixMilli()-now)/int64(orderAmountUntilRefill) + now)
}

func updateFiatBalance(fiat float64) {
	lastFiatBalance = fiat
}
