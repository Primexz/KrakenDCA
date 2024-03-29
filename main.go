package main

import (
	"log"

	"github.com/primexz/KrakenDCA/bot"
	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/prometheus"
)

func main() {
	log.Println("Starting Kraken DCA Bot")

	config.LoadConfiguration()
	go bot.StartBot()
	prometheus.StartServer()
}
