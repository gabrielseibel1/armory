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
	CharacterMounts() echo.HandlerFunc
	AllCharacterAchievements() echo.HandlerFunc
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

func (a *provider) CharacterMounts() echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := a.q.CharacterMounts(c.FormValue("character")) //todo: get by parameter
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, m, "\t")
	}
}

func (a *provider) AllCharacterAchievements() echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := a.q.AllCharacterAchievements(c.FormValue("character"))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, m, "\t")
	}
}

