package main

import (
	"log"

	"payslip-generation-system/config"

	app "payslip-generation-system/internal/app"
)

func main() {
	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	// init http app
	app := app.NewAppHTTP(cfg)

	// run http app
	app.Run(cfg)
}