package app

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
	billDelivery "github.com/fairyhunter13/tax-calculator/internal/bill/delivery"
	billRepository "github.com/fairyhunter13/tax-calculator/internal/bill/repository"
	billUsecase "github.com/fairyhunter13/tax-calculator/internal/bill/usecase"

	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
	taxDelivery "github.com/fairyhunter13/tax-calculator/internal/taxobj/delivery"
	taxRepository "github.com/fairyhunter13/tax-calculator/internal/taxobj/repository"
	taxUsecase "github.com/fairyhunter13/tax-calculator/internal/taxobj/usecase"

	"github.com/labstack/echo"
	ini "gopkg.in/ini.v1"

	//Using the pq library for the database.
	_ "github.com/lib/pq"
)

//App defines the group of connection, config, repository, usecase, and etc.
type App struct {
	config    *Config
	pool      *sql.DB
	billRepo  bill.Repository
	billUcase bill.Usecase
	taxRepo   taxobj.Repository
	taxUcase  taxobj.Usecase
	echoMux   *echo.Echo
}

//Config define all configs needed to store configured variables.
type Config struct {
	Database
	Server
}

//Database define the config for conection string.
type Database struct {
	ConnectionString string `ini:"connection"`
}

//Server define the config for server port to start the apps.
type Server struct {
	Port string `ini:"port"`
}

//NewApp return the new defined App
func NewApp() *App {
	return new(App)
}

//Migrate runs the migration script for the application.
func (app *App) Migrate() (err error) {
	err = app.taxRepo.Migrate()
	if err != nil {
		return
	}
	err = app.billUcase.LoadData()
	return
}

//ParseConfig parse config defined in the path to the given struct.
func (app *App) ParseConfig(configPath string, appConfig *Config) (err error) {
	err = ini.MapTo(appConfig, configPath)
	if err != nil {
		return
	}
	return
}

//SetConfig set the parsed config to the app.
func (app *App) SetConfig(config *Config) {
	app.config = config
}

//Init begin the initialization of application.
//This process initialize all connection, usecase, repositories, config, and etc.
func (app *App) Init() (err error) {
	app.pool, err = sql.Open("postgres", app.config.Database.ConnectionString)
	if err != nil {
		return
	}
	app.billRepo = billRepository.NewCacheRepository()
	app.taxRepo = taxRepository.NewPqRepository(app.pool)
	app.billUcase = billUsecase.NewBillUsecase(app.billRepo, app.taxRepo)
	app.taxUcase = taxUsecase.NewTaxObjectUsecase(app.taxRepo, app.billRepo)
	app.echoMux = echo.New()
	billDelivery.NewHTTPBillHandler(app.echoMux, app.billUcase)
	taxDelivery.NewTaxObjectHandler(app.echoMux, app.taxUcase)
	return
}

//Run run the application using graceful shutdown
func (app *App) Run(osSignal chan os.Signal) (err error) {
	go func() {
		if err = app.echoMux.Start(app.config.Server.Port); err != nil {
			return
		}
	}()
	<-osSignal

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = app.echoMux.Shutdown(ctx)
	return
}

//Close closes the app and all connections.
func (app *App) Close() {
	app.taxRepo.Close()
	app.pool.Close()
}
