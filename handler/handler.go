package handler

import (
	"github.com/gabrielseibel1/armory/data"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Provider interface {
	Hello() echo.HandlerFunc
	Home() echo.HandlerFunc
	Tables() echo.HandlerFunc
}

type provider struct{
	q data.Queryer
}

func NewArmoryProvider(q data.Queryer) Provider {
	return &provider{q: q}
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

func (a *provider) Tables() echo.HandlerFunc {
	return func(c echo.Context) error {
		t, err := a.q.Tables()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, strings.Join(t, "\n"))
	}
}