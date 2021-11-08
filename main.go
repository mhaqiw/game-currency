package main

import (
	"database/sql"
	"fmt"
	"github.com/mhaqiw/game-currency/repository"
	"github.com/mhaqiw/game-currency/service"
	"log"
	"time"

	"github.com/mhaqiw/game-currency/controller"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/mhaqiw/game-currency/util"
)

func main() {
	db := initDB()
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	timeoutContext := time.Duration(10) * time.Second

	cu := repository.NewCurrencyRepository(db)
	co := repository.NewConversionRepository(db)

	currencyService := service.NewCurrencyService(cu, timeoutContext)
	conversionService := service.NewConversionService(co, cu, timeoutContext)

	controller.NewCurrencyHandler(e, currencyService)
	controller.NewConversionHandler(e, conversionService)

	log.Fatal(e.Start(":9090"))
}

func initDB() *sql.DB {
	postgresHost := util.MustHaveEnv("POSTGRES_HOST")
	postgresPort := util.MustHaveEnv("POSTGRES_PORT")
	postgresUser := util.MustHaveEnv("POSTGRES_USER")
	postgresPassword := util.MustHaveEnv("POSTGRES_PASSWORD")
	postgresDbname := util.MustHaveEnv("POSTGRES_DB")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgresHost, postgresPort, postgresUser, postgresPassword, postgresDbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
