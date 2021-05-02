package main

import (
	"github.com/gabrielseibel1/armory/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	a := handler.NewArmoryProvider()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", a.Hello())
	e.GET("/", a.Home())

	e.Logger.Fatal(e.Start(":8742"))
}