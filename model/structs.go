package model

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type Character struct {
	Name              string        `json:"nome"`
	Title             string        `json:"titulo"`
	Spec              string        `json:"especialização"`
	Class             string        `json:"classe"`
	Race              string        `json:"raça"`
	Guild             string        `json:"guilda"`
	Level             string        `json:"nivel"`
	ILvl              int           `json:"i_lvl"`
	Faction           string        `json:"facção"`
	AchievementsScore int           `json:"escore_conquista"`
	IsGuildAdmin      bool          `json:"administrador_de_guilda,omitempty"`
	GuildReputation   int           `json:"reputação_de_guilda,omitempty"`
	Mounts            []Mount       `json:"montarias,omitempty"`
	Equipments        []Equipment   `json:"equipamentos,omitempty"`
	Achievements      []Achievement `json:"conquistas,omitempty"`
}

type CharacterWithAttributes struct {
	Character
	Health         int `json:"saude"`
	Resource       int `json:"recurso"`
	Stamina        int `json:"estamina"`
	Strength       int `json:"força"`
	Intellect      int `json:"intelecto"`
	Agility        int `json:"agilidade"`
	CriticalStrike int `json:"critico"`
	Haste          int `json:"aceleração"`
	Versatility    int `json:"versatilidade"`
	Mastery        int `json:"maestria"`
}

type CharacterWithAchievementsPanel struct {
	Character
	AchievementPanel []AchievementType `json:"painel_de_conquistas"`
}

type AchievementType struct {
	Name   string `json:"nome"`
	Points int    `json:"pontos"`
}

type Achievement struct {
	Name         string        `json:"nome"`
	Points       int           `json:"pontos"`
	Type         string        `json:"tipo,omitempty"`
	SubType      string        `json:"subtipo,omitempty"`
	Date         time.Time     `json:"data,omitempty"`
	Requirements []Achievement `json:"requisitos,omitempty"`
}

type AchievPointsInMonth struct {
	Year   int        `json:"ano"`
	Month  time.Month `json:"mês"`
	Points int        `json:"pontos"`
}

type Mount struct {
	Name        string `json:"nome"`
	Obtainment  string `json:"obtenção"`
	Description string `json:"descrição"`
	IsFavorite  bool   `json:"favorita"`
}

type Item struct {
	Part             string   `json:"parte"`
	Name             string   `json:"nome"`
	Level            int      `json:"nivel"`
	Stamina          int      `json:"estamina"`
	Strength         int      `json:"força"`
	Intellect        int      `json:"intelecto"`
	Agility          int      `json:"agilidade"`
	CriticalStrike   int      `json:"critico"`
	Haste            int      `json:"aceleração"`
	Versatility      int      `json:"versatilidade"`
	Mastery          int      `json:"maestria"`
	LevelRequirement null.Int `json:"nivel_min"`
	Price            null.Int `json:"preço"`
	MaxDurability    null.Int `json:"durabilidade_max"`
}

type Equipment struct {
	Item
	CurDurability null.Int `json:"durabilidade_atual"`
}
