package bot

import (
	"fmt"
	"time"

	"github.com/primexz/KrakenDCA/prometheus"
)

func computeNextFiatDepositDay() time.Time {
	date := addMonthsToTime(1, time.Now().UTC())

	//set metrics
	prometheus.TimeOfEmptyFiat.Set(float64(date.UnixMilli()))

	return date
}

func addMonthsToTime(months int, time time.Time) time.Time {
	return time.AddDate(0, months, 0)
}

func fmtDuration(d time.Duration) string {
	hour := int(d.Hours())
	minute := int(d.Minutes()) % 60
	second := int(d.Seconds()) % 60

	return fmt.Sprintf("%02dh %02dm %02ds", hour, minute, second)
}
