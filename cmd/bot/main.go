package main

import "bot/internal/app/bot"

func main() {
	app, err := bot.NewApp()
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
