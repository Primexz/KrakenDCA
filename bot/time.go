package bot

import (
	"fmt"
	"time"
)

func fmtDuration(d time.Duration) string {
	hour := int(d.Hours())
	minute := int(d.Minutes()) % 60
	second := int(d.Seconds()) % 60

	return fmt.Sprintf("%02dh %02dm %02ds", hour, minute, second)
}
