package bot

import (
	"time"

	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/kraken"
	"github.com/primexz/KrakenDCA/logger"
)

type Bot struct {
	timeOfNextOrder  time.Time
	timeOfEmptyFiat  time.Time
	lastFiatBalance  float64
	lastBtcFiatPrice float64
	fiatAmount       float64
	initialRun       bool
}

func NewBot() *Bot {
	return &Bot{
		initialRun: true,
	}
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("bot")
}

func (b *Bot) StartBot() {
	go func() {

		for {
			b.run()

			b.initialRun = false
			time.Sleep(time.Duration(config.CheckDelay) * time.Second)
		}

	}()
}

func (b *Bot) run() {
	log.Info("Starting bot")

	if fiatAmnt, err := kraken.GetFiatBalance(); err == nil {
		b.fiatAmount = fiatAmnt
	} else {
		log.Error("Error getting fiat balance: ", err)
		return
	}

	if b.fiatAmount == 0 {
		log.Warn("No remaining fiat balance found. It's time to top up your account ;)")
		return
	}

	newFiatMoney := b.fiatAmount > b.lastFiatBalance
	if newFiatMoney || b.initialRun {
		if b.initialRun {
			log.Info("Initial run. Calculating next fiat deposit day..")
		} else {
			log.Info("New fiat deposit found. 💰")
		}

		b.timeOfEmptyFiat = computeNextFiatDepositDay()

		log.Info("Next Fiat deposit required at", b.timeOfEmptyFiat)

		b.lastFiatBalance = b.fiatAmount
	}

	if fiatPrice, err := kraken.GetCurrentBtcFiatPrice(); err == nil {
		b.lastBtcFiatPrice = fiatPrice
	} else {
		log.Error("Error getting current btc price:", err)
		return
	}

	if b.initialRun {
		b.calculateTimeOfNextOrder()
		b.logNextOrder()
	}

	if (b.timeOfNextOrder.Before(time.Now()) || newFiatMoney) && !b.initialRun {
		log.Info("Placing bitcoin purchase order. ₿")

		kraken.BuyBtc()
		b.calculateTimeOfNextOrder()
		b.logNextOrder()
	}
}

func (b *Bot) logNextOrder() {
	log.Info("Next order in", fmtDuration(time.Until(b.timeOfNextOrder)))
}
