package main

import (
	"bot/internal/app/bot"
	"fmt"
	"os"

	pkg_config "bot/internal/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "run: %s\n", err.Error()) //nolint:revive
		os.Exit(1)
	}
}

func run() error {
	config, err := pkg_config.NewConfig()
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}

	app := bot.NewApp(config)

	err = app.Run()
	if err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}
