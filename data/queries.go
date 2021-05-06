package data

import (
	"database/sql"
	"github.com/gabrielseibel1/armory/model"
)

type Queryer interface {
	Tables() ([]string, error)
	CharacterMounts(char string) ([]model.Mount, error)
	AllCharacterAchievements(char string) ([]model.Achievement, error)

	//*Q0 - todo: query view with character (simple) + spec + class + race + guild

	//*Q1 - todo: query character view + achievements
		//AllCharacterAchievements(char string) ([]model.Achievement, error) // todo: integrate character view
	//*Q2 - todo: query character view + achievement panel
		//CharacterAchievementTypePoints(char string) ([]model.AchievementTypePoints, error)
	//*Q3 - todo: query character view + achievement list per category
		//CharacterAchievements(char, achType, achSubtype string) ([]Achievement, error)
	//*Q4 - todo: guild achievements (distinct members'achievs)
		//GuildAchievements(guild string) ([]Achievement, error)

	//*Q5 - todo: query character view + equipment

	//*Q6 - todo: query character view + mounts
		// CharacterMounts(char string) ([]model.Mount, error) // todo: integrate character view

	//*Q8 - todo: query achievement points history grouped by month since a given year (group by (month, year) having year > x)

	//*Q9 - todo: select all equips where ilvl > ilvl of character with max(ilvl)
	//*Q10 - todo: select guilds with achiev points (composed by members') < achiev points of character with max(achiev points)
	//idea:
	//select p.nome from personagens p
	//where p.escore_conquistas > (
	//		select sum(c.pontos) from conquistas c
	//		join personagens_conquistas pc on pc.fk_conquistas_id = c.id
	//		join personagens p on p.id = pc.fk_personagens_id
	//		join guildas g on p.fk_guilda_id = g.id
	//		group by g.id
	//		order by desc sum(c.pontos)
	//	)

	//*Q11 - todo: query todos os achievements que nenhum personagem conquistou
	//select c.nome from conquistas c
	//where not exists ()

	//Q12 - todo: CharactersAchievementRanking() ([]Character, error) //(players ordered by ach points)
}

type queries struct {
	db *sql.DB
}

func NewQueryer(db *sql.DB) Queryer {
	return &queries{db: db}
}

func (q *queries) Tables() ([]string, error) {
	rows, err := q.db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name;
	`)
	if err != nil {
		return []string{}, err
	}

	var tables []string
	var t string
	for rows.Next() {
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	return tables, nil
}

func (q *queries) CharacterMounts(char string) ([]model.Mount, error) {
	stmt, err := q.db.Prepare(`
		SELECT m.id, m.nome, m.obtencao, m.descricao FROM montarias m
		JOIN personagens_montarias as pm ON (m.id = pm.fk_montaria_id)
		JOIN personagens p ON (pm.fk_personagens_id = p.id)
		WHERE p.nome = $1;
	`)
	if err != nil {
		return []model.Mount{}, err
	}

	rows, err := stmt.Query(char)
	if err != nil {
		return []model.Mount{}, err
	}

	var mounts []model.Mount
	for rows.Next() {
		var mount model.Mount
		err = rows.Scan(&mount.Id, &mount.Name, &mount.Description, &mount.Obtainment)
		if err != nil {
			return nil, err
		}
		mounts = append(mounts, mount)
	}

	return mounts, nil
}

func (q *queries) AllCharacterAchievements(char string) ([]model.Achievement, error) {
	stmt, err := q.db.Prepare(`
		SELECT c.id, c.nome, c.pontos FROM conquistas c
		JOIN personagens_conquistas pc ON (c.id = pc.fk_conquistas_id)
		JOIN personagens p ON (pc.fk_personagens_id = p.id)
		WHERE p.nome = $1;
	`)
	if err != nil {
		return []model.Achievement{}, err
	}

	rows, err := stmt.Query(char)
	if err != nil {
		return []model.Achievement{}, err
	}

	var achs []model.Achievement
	for rows.Next() {
		var a model.Achievement
		err = rows.Scan(&a.Id, &a.Name, &a.Points)
		if err != nil {
			return nil, err
		}
		achs = append(achs, a)
	}

	return achs, nil
}