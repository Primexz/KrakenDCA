package bot

import (
	"time"

	"github.com/primexz/KrakenDCA/config"
)

func computeNextFiatDepositDay() time.Time {
	date := addMonthsToTime(1, time.Now())

	//get the first day of the month
	return date.AddDate(0, 0, -date.Day()+1)
}

func (b *Bot) calculateTimeOfNextOrder() {
	orderAmountUntilRefill := b.getOrderAmountUntilRefill()

	if orderAmountUntilRefill < 1 {
		b.log.Error("Fiat balance is too low to make an order.")
		b.timeOfNextOrder = b.timeOfEmptyFiat
		return
	}

	now := time.Now().UnixMilli()
	b.timeOfNextOrder = time.UnixMilli((b.timeOfEmptyFiat.UnixMilli()-now)/int64(orderAmountUntilRefill) + now)
}

func (b *Bot) getOrderAmountUntilRefill() float64 {
	fiatValueInBtc := b.fiatAmount / b.lastBtcFiatPrice

	return fiatValueInBtc / config.C.KrakenOrderSize
}
