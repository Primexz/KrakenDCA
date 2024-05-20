package bot

import (
	"time"

	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/kraken"
	"github.com/primexz/KrakenDCA/metrics"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	timeOfNextOrder  time.Time
	timeOfEmptyFiat  time.Time
	lastFiatBalance  float64
	lastBtcFiatPrice float64
	fiatAmount       float64
	initialRun       bool

	log       *log.Entry
	krakenApi *kraken.KrakenApi
}

func NewBot() *Bot {
	return &Bot{
		initialRun: true,
		log: log.WithFields(log.Fields{
			"prefix": "bot",
		}),
		krakenApi: kraken.NewKrakenApi(),
	}
}

func (b *Bot) StartBot() {
	go func() {
		b.log.Info("Starting bot")

		for {
			b.run()

			b.initialRun = false
			time.Sleep(time.Duration(config.CheckDelay) * time.Second)
		}

	}()
}

func (b *Bot) run() {
	if fiatAmnt, err := b.krakenApi.GetFiatBalance(); err == nil {
		b.fiatAmount = fiatAmnt
	} else {
		b.log.Error("Error getting fiat balance: ", err)
		return
	}

	if b.fiatAmount == 0 {
		b.log.Warn("No remaining fiat balance found. It's time to top up your account ;)")
		return
	}

	newFiatMoney := b.fiatAmount > b.lastFiatBalance
	if newFiatMoney || b.initialRun {
		if b.initialRun {
			b.log.Info("Initial run. Calculating next fiat deposit day..")
		} else {
			b.log.Info("New fiat deposit found. 💰")
		}

		b.timeOfEmptyFiat = computeNextFiatDepositDay()

		b.log.WithFields(log.Fields{
			"time": b.timeOfEmptyFiat,
		}).Info("Next fiat deposit in")

		b.lastFiatBalance = b.fiatAmount
	}

	if fiatPrice, err := b.krakenApi.GetCurrentBtcFiatPrice(); err == nil {
		b.lastBtcFiatPrice = fiatPrice
	} else {
		b.log.Error("Error getting current btc price:", err)
		return
	}

	if b.initialRun {
		b.calculateTimeOfNextOrder()
	}

	if (b.timeOfNextOrder.Before(time.Now()) || newFiatMoney) && !b.initialRun {
		b.log.Info("Placing bitcoin purchase order. ₿")

		b.krakenApi.BuyBtc()
		b.calculateTimeOfNextOrder()
	}

	b.updateMetrics()

	b.log.WithFields(log.Fields{
		"fiat_balance": b.fiatAmount,
		"duration":     fmtDuration(time.Until(b.timeOfNextOrder)),
	}).Info("Next order in")
}

func (b *Bot) updateMetrics() {
	metrics.Metrics.NextOrder = b.timeOfNextOrder.Unix()
}
