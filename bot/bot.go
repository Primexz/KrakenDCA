package bot

import (
	"time"

	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/kraken"
	"github.com/primexz/KrakenDCA/logger"
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
	for {
		run()

		initialRun = false
		time.Sleep(time.Duration(config.CheckDelay) * time.Second)
	}

}

var log *logger.Logger

func init() {
	log = logger.NewLogger("bot")
}

func run() {
	log.Info("Starting bot")

	if fiatAmnt, err := kraken.GetFiatBalance(); err == nil {
		fiatAmount = fiatAmnt
	} else {
		log.Error("Error getting fiat balance: ", err)
		return
	}

	if fiatAmount == 0 {
		log.Warn("No remaining fiat balance found. It's time to top up your account ;)")
		return
	}

	newFiatMoney := fiatAmount > lastFiatBalance
	if newFiatMoney || initialRun {
		if initialRun {
			log.Info("Initial run. Calculating next fiat deposit day..")
		} else {
			log.Info("New fiat deposit found. 💰")
		}

		timeOfEmptyFiat = computeNextFiatDepositDay()

		log.Info("Next Fiat deposit required at", timeOfEmptyFiat)

		lastFiatBalance = fiatAmount
	}

	if fiatPrice, err := kraken.GetCurrentBtcFiatPrice(); err == nil {
		lastBtcFiatPrice = fiatPrice
	} else {
		log.Error("Error getting current btc price:", err)
		return
	}

	if initialRun {
		calculateTimeOfNextOrder()
		logNextOrder()
	}

	if (timeOfNextOrder.Before(time.Now()) || newFiatMoney) && !initialRun {
		log.Info("Placing bitcoin purchase order. ₿")

		kraken.BuyBtc()
		calculateTimeOfNextOrder()
		logNextOrder()
	}
}

func logNextOrder() {
	log.Info("Next order in", fmtDuration(time.Until(timeOfNextOrder)), timeOfNextOrder)
}
