package handler

import (
	"github.com/gabrielseibel1/armory/data"
	"github.com/gabrielseibel1/armory/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Provider interface {
	Home() echo.HandlerFunc

	CharactersRanking() echo.HandlerFunc
	CharactersThatAchievedMoreThanGuilds() echo.HandlerFunc
	ItemsWithHigherLevelThanHighestPlayerILvl() echo.HandlerFunc

	Character() echo.HandlerFunc
	CharacterAndMounts() echo.HandlerFunc
	CharacterAndEquipments() echo.HandlerFunc
	CharacterAndAchievements() echo.HandlerFunc
	CharacterAndAchievementsPanel() echo.HandlerFunc
	CharacterAchievementPointsPerMonth() echo.HandlerFunc

	GuildAchievements() echo.HandlerFunc
	UnearnedAchievements() echo.HandlerFunc
	AchievementsWithRequirements() echo.HandlerFunc
}

type provider struct {
	q data.Queryer
}

func NewArmoryProvider(q data.Queryer) Provider {
	return &provider{q: q}
}

func (p *provider) Home() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File("index.html")
	}
}

func (p *provider) CharactersRanking() echo.HandlerFunc {
	return func(c echo.Context) error {
		chars, err := p.q.CharactersRanking()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, chars, "\t")
	}
}

func (p *provider) CharactersThatAchievedMoreThanGuilds() echo.HandlerFunc {
	return func(c echo.Context) error {
		chars, err := p.q.CharactersThatAchievedMoreThanGuilds()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, chars, "\t")
	}
}

func (p *provider) ItemsWithHigherLevelThanHighestPlayerILvl() echo.HandlerFunc {
	return func(c echo.Context) error {
		equips, err := p.q.ItemsWithHigherLevelThanHighestPlayerILvl()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, equips, "\t")
	}
}

func (p *provider) Character() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		char, err := p.q.Character(ctx.FormValue("character"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSONPretty(http.StatusOK, char, "\t")
	}
}

func (p *provider) CharacterAndMounts() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		char, err := p.q.CharacterAndMounts(ctx.FormValue("character"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSONPretty(http.StatusOK, char, "\t")
	}
}

func (p *provider) CharacterAndEquipments() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		char, err := p.q.CharacterAndEquipments(ctx.FormValue("character"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSONPretty(http.StatusOK, char, "\t")
	}
}

func (p *provider) CharacterAndAchievements() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		charName := ctx.FormValue("character")
		achType := ctx.FormValue("type")
		achSubtype := ctx.FormValue("subtype")

		var char model.Character
		var err error
		if achType != "" && achSubtype != "" {
			char, err = p.q.CharacterAndAchievementsPerCategory(charName, achType, achSubtype)
		} else {
			char, err = p.q.CharacterAndAchievements(charName)
		}
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSONPretty(http.StatusOK, char, "\t")
	}
}

func (p *provider) CharacterAndAchievementsPanel() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		char, err := p.q.CharacterAndAchievementsPanel(ctx.FormValue("character"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSONPretty(http.StatusOK, char, "\t")
	}
}

func (p *provider) CharacterAchievementPointsPerMonth() echo.HandlerFunc {
	return func(c echo.Context) error {
		char := c.FormValue("character")
		minPoints, err := strconv.Atoi(c.FormValue("minPoints"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		months, err := p.q.CharacterAchievementPointsPerMonth(char, minPoints)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, months, "\t")
	}
}

func (p *provider) GuildAchievements() echo.HandlerFunc {
	return func(c echo.Context) error {
		a, err := p.q.GuildAchievements(c.FormValue("guild"))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, a, "\t")
	}
}

func (p *provider) UnearnedAchievements() echo.HandlerFunc {
	return func(c echo.Context) error {
		a, err := p.q.UnearnedAchievements()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, a, "\t")
	}
}

func (p *provider) AchievementsWithRequirements() echo.HandlerFunc {
	return func(c echo.Context) error {
		a, err := p.q.AchievementsWithRequirements()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSONPretty(http.StatusOK, a, "\t")
	}
}
