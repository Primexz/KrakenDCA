package bot

import (
	"time"

	"github.com/primexz/KrakenDCA/config"
)

func computeNextFiatDepositDay() time.Time {
	now := time.Now()

	year := now.Year()
	month := now.Month() + 1

	if month > 12 {
		month = 1
		year++
	}

	firstOfNextMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return firstOfNextMonth
}

func (b *Bot) calculateTimeOfNextOrder() {
	orderAmountUntilRefill := b.getOrderAmountUntilRefill()

	if orderAmountUntilRefill < 1 {
		b.log.Error("Fiat balance is too low to make an order.")
		b.timeOfNextOrder = time.Now().AddDate(0, 1, 0)
		return
	}

	now := time.Now().UnixMilli()
	b.timeOfNextOrder = time.UnixMilli((b.timeOfEmptyFiat.UnixMilli()-now)/int64(orderAmountUntilRefill) + now)
}

func (b *Bot) getOrderAmountUntilRefill() float64 {
	fiatValueInBtc := b.fiatAmount / b.lastBtcFiatPrice

	return fiatValueInBtc / config.C.KrakenOrderSize
}
