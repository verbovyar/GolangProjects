package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/app"
	"log"
)

func main() {
	config, err := config.LoadConfig("././config")
	if err != nil {
		log.Fatalf("%v", err)
	}

	app.Run(config)
}
