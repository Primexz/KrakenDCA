package bot

import (
	"fmt"
	"time"
)

func computeNextFiatDepositDay() time.Time {
	date := addMonthsToTime(1, time.Now())

	//get the first day of the month
	return date.AddDate(0, 0, -date.Day()+1)
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
