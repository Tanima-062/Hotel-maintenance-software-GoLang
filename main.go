package main

import (
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/Adventureinc/hotel-hm-api/src/common/app"
	"github.com/Adventureinc/hotel-hm-api/src/common/infra"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/joho/godotenv"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

const location = "Asia/Tokyo"

func main() {

	err := godotenv.Load()
		if err != nil {
		panic(".env not loaded")
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
	e := echo.New()
	// log level
	e.Logger.SetLevel(log.INFO)

	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// debug用cors設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("ALLOW_HOST")},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// 最初に発行しておきたいやつ
	hotelDB, err := infra.DBCon()
	if err != nil {
		e.Logger.Fatal(err)
	}
	// ルーティング
	app.Route(e, hotelDB)

	e.Logger.Fatal(e.Start(":1323"))
}
