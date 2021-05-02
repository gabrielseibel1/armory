package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Provider interface {
	Hello() echo.HandlerFunc
	Home() echo.HandlerFunc
}

type ArmoryProvider struct{}

func NewArmoryProvider() *ArmoryProvider {
	return &ArmoryProvider{}
}

func (a *ArmoryProvider) Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Armory: Hello, World!")
	}
}

func (a *ArmoryProvider) Home() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File("index.html")
	}
}
