package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Provider interface {
	Hello() echo.HandlerFunc
	Home() echo.HandlerFunc
}

type provider struct{}

func NewArmoryProvider() Provider {
	return &provider{}
}

func (a *provider) Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Armory: Hello, World!")
	}
}

func (a *provider) Home() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File("index.html")
	}
}
