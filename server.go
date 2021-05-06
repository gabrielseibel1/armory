package main

import (
	"database/sql"
	"github.com/gabrielseibel1/armory/data"
	"github.com/gabrielseibel1/armory/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	"os"
)

func main() {
	dsn, ok := os.LookupEnv("ARMORY_DSN")
	if !ok {
		panic("No connection information on variable ARMORY_DSN")
	}

	conn, err := pq.NewConnector(dsn)
	if err != nil {
		panic(err)
	}
	db:= sql.OpenDB(conn)
	defer db.Close()

	q := data.NewQueryer(db)

	h := handler.NewArmoryProvider(q)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", h.Home())
	e.GET("/hello", h.Hello())
	e.GET("/tables", h.Tables())
	e.GET("/mounts", h.CharacterMounts())
	e.GET("/achievements", h.AllCharacterAchievements())

	e.Logger.Fatal(e.Start(":8742"))
}