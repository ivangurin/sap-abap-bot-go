package main

import (
	"bot/internal/app/bot"
	"log"
)

func main() {
	app, err := bot.NewApp()
	if err != nil {
		log.Default().Fatalf("create app: %s", err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Default().Fatalf("run app: %s", err.Error())
	}
}
