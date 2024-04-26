package main

import (
	"log"

	"github.com/primexz/KrakenDCA/bot"
	"github.com/primexz/KrakenDCA/config"
)

func main() {
	log.Println("Starting Kraken DCA Bot")

	config.LoadConfiguration()
	bot.StartBot()
}
