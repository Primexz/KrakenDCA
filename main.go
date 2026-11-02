package main

import (
	"runtime"

	"github.com/primexz/KrakenDCA/bot"
	"github.com/primexz/KrakenDCA/metrics"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat:  "2006/01/02 - 15:04:05",
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		SpacePadding:     45,
	})

	log.SetReportCaller(true)
}

func main() {
	log.WithFields(log.Fields{
		"commit":  commit,
		"runtime": runtime.Version(),
		"arch":    runtime.GOARCH,
	}).Infof("Kraken DCA 🐙 %s modded by Jag", version)

	metrics.StartServer()

	bot.NewBot().StartBot()
}
