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

func calculateTimeOfNextOrder() {
	orderAmountUntilRefill := getOrderAmountUntilRefill()

	now := time.Now().UnixMilli()
	timeOfNextOrder = time.UnixMilli((timeOfEmptyFiat.UnixMilli()-now)/int64(orderAmountUntilRefill) + now)
}

func getOrderAmountUntilRefill() float64 {
	fiatValueInBtc := fiatAmount / lastBtcFiatPrice

	return fiatValueInBtc / config.KrakenOrderSize
}
