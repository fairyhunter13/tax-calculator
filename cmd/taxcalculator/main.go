package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fairyhunter13/tax-calculator/internal/app"
)

const (
	configPath = "/configs/config.ini"
)

var (
	osSignal = make(chan os.Signal, 1)
)

func init() {
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT)
}

func main() {
	application := app.NewApp()
	defer application.Close()
	defer close(osSignal)
	appConfig := new(app.Config)
	err := application.ParseConfig(configPath, appConfig)
	if err != nil {
		log.Fatalf("[App] Failed to parse the config: %s", err)
	}
	application.SetConfig(appConfig)
	err = application.Init()
	if err != nil {
		log.Fatalf("[App] Failed to init the application: %s", err)
	}
	err = application.Migrate()
	if err != nil {
		log.Printf("[App] Failed to migrate the database: %s", err)
	}
	err = application.Run(osSignal)
	if err != nil {
		log.Fatalf("[App] Failed to run the application: %s", err)
	}
	log.Printf("[App] Shutting down!")
}
