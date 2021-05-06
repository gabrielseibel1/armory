package data

import (
	"database/sql"
	"github.com/gabrielseibel1/armory/model"
)

type Queryer interface {
	CharactersRanking() ([]model.Character, error)
	CharactersThatAchievedMoreThanGuilds() ([]model.Character, error)
	ItemsWithHigherLevelThanHighestPlayerILvl() ([]model.Item, error)

	Character(char string) (model.CharacterWithAttributes, error)
	CharacterAndMounts(char string) (model.Character, error)
	CharacterAndEquipments(char string) (model.CharacterWithAttributes, error)
	CharacterAndAchievements(char string) (model.Character, error)
	CharacterAndAchievementsPerCategory(char, achType, achSubtype string) (model.Character, error)
	CharacterAndAchievementsPanel(char string) (model.CharacterWithAchievementsPanel, error)
	CharacterAchievementPointsPerMonth(char string, minPoints int) ([]model.AchievPointsInMonth, error)

	GuildAchievements(guild string) ([]model.Achievement, error)
	UnearnedAchievements() ([]model.Achievement, error)
	AchievementsWithRequirements() ([]model.Achievement, error)
}

type queryer struct {
	db *sql.DB
}

func NewQueryer(db *sql.DB) Queryer {
	return &queryer{db: db}
}

func (q *queryer) CharactersRanking() ([]model.Character, error) {
	rows, err := q.db.Query(`
		select 
		       p.personagem_nome, 
		       p.titulo, 
		       p.especializacao_nome, 
		       p.classe_nome, 
		       p.raca_nome, 
		       p.guilda_nome,
		       p.nivel, 
		       p.ilvl, 
		       p.faccao, 
		       sum(c.pontos) 
		from personagens_view p
		join personagens_conquistas pc on pc.fk_personagens_id = p.id
		join conquistas c on c.id = pc.fk_conquistas_id 
		group by 
		         p.personagem_nome,
		         p.titulo, 
		         p.especializacao_nome, 
		         p.classe_nome, 
		         p.raca_nome, 
		         p.guilda_nome,
		         p.nivel, 
		         p.ilvl, 
		         p.faccao
		order by sum(c.pontos) desc;
	`)
	if err != nil {
		return nil, err
	}

	var chars []model.Character
	for rows.Next() {
		c, err := q.scanCharacter(rows)
		if err != nil {
			return nil, err
		}
		chars = append(chars, c)
	}
	return chars, nil
}

func (q *queryer) CharactersThatAchievedMoreThanGuilds() ([]model.Character, error) {
	rows, err := q.db.Query(`
		select p.personagem_nome, p.titulo, p.especializacao_nome, p.classe_nome, p.raca_nome, p.guilda_nome,
		p.nivel, p.ilvl, p.faccao, p.escore_conquistas
		from personagens_view p
		where p.escore_conquistas > (
			select sum(c.pontos) from conquistas c
			join personagens_conquistas pc on pc.fk_conquistas_id = c.id
			join personagens p on p.id = pc.fk_personagens_id
			join guildas g on p.fk_guilda_id = g.id
			group by g.id
			order by sum(c.pontos) asc
			limit 1
		);
	`)
	if err != nil {
		return nil, err
	}

	var chars []model.Character
	for rows.Next() {
		c, err := q.scanCharacter(rows)
		if err != nil {
			return nil, err
		}
		chars = append(chars, c)
	}
	return chars, nil
}

func (q *queryer) ItemsWithHigherLevelThanHighestPlayerILvl() ([]model.Item, error) {
	rows, err := q.db.Query(`
		select pa.nome,
			   e.nome,
			   e.nivel,
			   e.estamina,
			   e.forca,
			   e.intelecto,
			   e.agilidade,
			   e.critico,
			   e.aceleracao,
			   e.versatilidade,
			   e.maestria,
			   e.nivel_min,
			   e.preco,
			   e.durabilidade_max
		from equipamentos e
			join partes pa on e.fk_parte_id = pa.id
		where e.nivel > (
			select round(avg(e.nivel))
			from personagens p
					 join personagens_equipamentos pe on pe.fk_personagens_id = p.id
					 join equipamentos e on e.id = pe.fk_equipamentos_id
			group by p.id
			order by avg(e.nivel) desc
			limit 1
		);
	`)
	if err != nil {
		return nil, err
	}

	var equips []model.Item
	for rows.Next() {
		var e model.Item
		err := rows.Scan(
			&e.Part,
			&e.Name, &e.Level, &e.Stamina,
			&e.Strength, &e.Intellect, &e.Agility,
			&e.CriticalStrike, &e.Haste, &e.Versatility, &e.Mastery,
			&e.LevelRequirement, &e.Price, &e.MaxDurability,
		)
		if err != nil {
			return nil, err
		}
		equips = append(equips, e)
	}
	return equips, nil
}

func (q *queryer) Character(char string) (model.CharacterWithAttributes, error) {
	stmt, err := q.db.Prepare(`
		select p.personagem_nome,
			   p.titulo,
			   p.especializacao_nome,
			   p.classe_nome,
			   p.raca_nome,
			   p.guilda_nome,
			   p.nivel,
			   p.saude,
			   p.recurso,
			   p.estamina,
			   p.forca,
			   p.agilidade,
			   p.intelecto,
			   p.critico,
			   p.aceleracao,
			   p.maestria,
			   p.versatilidade,
			   p.ilvl,
			   p.faccao,
			   p.adm,
			   p.reputacao,
			   p.escore_conquistas
		from personagens_view p
		where p.personagem_nome = $1;
	`)
	if err != nil {
		return model.CharacterWithAttributes{}, err
	}
	rows, err := stmt.Query(char)
	if err != nil {
		return model.CharacterWithAttributes{}, err
	}

	var c model.CharacterWithAttributes
	if rows.Next() {
		err := rows.Scan(
			&c.Name,
			&c.Title,
			&c.Spec,
			&c.Class,
			&c.Race,
			&c.Guild,
			&c.Level,
			&c.Health,
			&c.Resource,
			&c.Stamina,
			&c.Strength,
			&c.Agility,
			&c.Intellect,
			&c.CriticalStrike,
			&c.Haste,
			&c.Mastery,
			&c.Versatility,
			&c.ILvl,
			&c.Faction,
			&c.IsGuildAdmin,
			&c.GuildReputation,
			&c.AchievementsScore,
		)
		if err != nil {
			return model.CharacterWithAttributes{}, err
		}
	}
	return c, nil
}

func (q *queryer) CharacterAndMounts(char string) (model.Character, error) {
	stmt, err := q.db.Prepare(`
		select p.personagem_nome,
			   p.titulo,
			   p.especializacao_nome,
			   p.classe_nome,
			   p.raca_nome,
			   p.guilda_nome,
			   p.nivel,
			   p.ilvl,
			   p.faccao,
			   p.escore_conquistas,
			   m.nome,
			   m.obtencao,
			   m.descricao,
			   pm.favorita
		from personagens_view p
				 join personagens_montarias pm on pm.fk_personagens_id = p.id
				 join montarias m on m.id = pm.fk_montaria_id
		where p.personagem_nome = $1;
	`)
	if err != nil {
		return model.Character{}, err
	}
	rows, err := stmt.Query(char)
	if err != nil {
		return model.Character{}, err
	}

	var c model.Character
	var mounts []model.Mount
	for rows.Next() {
		var m model.Mount
		err := rows.Scan(
			&c.Name,
			&c.Title,
			&c.Spec,
			&c.Class,
			&c.Race,
			&c.Guild,
			&c.Level,
			&c.ILvl,
			&c.Faction,
			&c.AchievementsScore,
			&m.Name,
			&m.Obtainment,
			&m.Description,
			&m.IsFavorite,
		)
		if err != nil {
			return model.Character{}, err
		}
		mounts = append(mounts, m)
	}
	c.Mounts = mounts
	return c, nil
}

func (q *queryer) CharacterAndEquipments(char string) (model.CharacterWithAttributes, error) {
	stmt, err := q.db.Prepare(`
		select p.personagem_nome,
			   p.titulo,
			   p.especializacao_nome,
			   p.classe_nome,
			   p.raca_nome,
			   p.guilda_nome,
			   p.nivel,
			   p.saude,
			   p.recurso,
			   p.estamina,
			   p.forca,
			   p.agilidade,
			   p.intelecto,
			   p.critico,
		       p.aceleracao,
			   p.maestria,
			   p.versatilidade,
			   p.ilvl,
			   p.faccao,
		       p.escore_conquistas,
		       pa.nome,
			   e.nome,
			   e.nivel,
			   e.estamina,
			   e.forca,
			   e.intelecto,
			   e.agilidade,
			   e.critico,
			   e.aceleracao,
			   e.versatilidade,
			   e.maestria,
			   e.nivel_min,
			   e.preco,
			   e.durabilidade_max,
			   pe.durabilidade
		from personagens_view p
				 join personagens_equipamentos pe on pe.fk_personagens_id = p.id
				 join equipamentos e on e.id = pe.fk_equipamentos_id
        		 join partes pa on e.fk_parte_id = pa.id
		where p.personagem_nome = $1;
	`)
	if err != nil {
		return model.CharacterWithAttributes{}, err
	}
	rows, err := stmt.Query(char)
	if err != nil {
		return model.CharacterWithAttributes{}, err
	}

	var c model.CharacterWithAttributes
	var equips []model.Equipment
	for rows.Next() {
		var e model.Equipment
		err := rows.Scan(
			&c.Name,
			&c.Title,
			&c.Spec,
			&c.Class,
			&c.Race,
			&c.Guild,
			&c.Level,
			&c.Health,
			&c.Resource,
			&c.Stamina,
			&c.Strength,
			&c.Agility,
			&c.Intellect,
			&c.CriticalStrike,
			&c.Haste,
			&c.Mastery,
			&c.Versatility,
			&c.ILvl,
			&c.Faction,
			&c.AchievementsScore,
			&e.Part,
			&e.Name,
			&e.Level,
			&e.Stamina,
			&e.Strength,
			&e.Intellect,
			&e.Agility,
			&e.CriticalStrike,
			&e.Haste,
			&e.Versatility,
			&e.Mastery,
			&e.LevelRequirement,
			&e.Price,
			&e.MaxDurability,
			&e.CurDurability,
		)
		if err != nil {
			return model.CharacterWithAttributes{}, err
		}
		equips = append(equips, e)
	}
	c.Equipments = equips
	return c, nil
}

func (q *queryer) CharacterAndAchievements(char string) (model.Character, error) {
	stmt, err := q.db.Prepare(`
		select p.personagem_nome,
			   p.titulo,
			   p.especializacao_nome,
			   p.classe_nome,
			   p.raca_nome,
			   p.guilda_nome,
			   p.nivel,
			   p.ilvl,
			   p.faccao,
			   p.escore_conquistas,
			   c.nome,
			   c.pontos,
			   pc."data"
		from personagens_view p
				 join personagens_conquistas pc on pc.fk_personagens_id = p.id
				 join conquistas c on c.id = pc.fk_conquistas_id
		where p.personagem_nome = $1;
	`)
	if err != nil {
		return model.Character{}, err
	}
	rows, err := stmt.Query(char)
	if err != nil {
		return model.Character{}, err
	}

	var c model.Character
	var achs []model.Achievement
	for rows.Next() {
		var a model.Achievement
		err := rows.Scan(
			&c.Name,
			&c.Title,
			&c.Spec,
			&c.Class,
			&c.Race,
			&c.Guild,
			&c.Level,
			&c.ILvl,
			&c.Faction,
			&c.AchievementsScore,
			&a.Name,
			&a.Points,
			&a.Date,
		)
		if err != nil {
			return model.Character{}, err
		}
		achs = append(achs, a)
	}
	c.Achievements = achs
	return c, nil
}

func (q *queryer) CharacterAndAchievementsPerCategory(char, achType, achSubtype string) (model.Character, error) {
	stmt, err := q.db.Prepare(`
		select p.personagem_nome,
			   p.titulo,
			   p.especializacao_nome,
			   p.classe_nome,
			   p.raca_nome,
			   p.guilda_nome,
			   p.nivel,
			   p.ilvl,
			   p.faccao,
			   p.escore_conquistas,
			   tc.nome,
			   stc.nome,
			   c.nome,
			   c.pontos,
			   pc."data"
		from personagens_view p
				 join personagens_conquistas pc on pc.fk_personagens_id = p.id
				 join conquistas c on c.id = pc.fk_conquistas_id
				 join tipos_conquistas tc on tc.id = c.fk_tipos_conquistas_id
				 join subtipos_conquistas stc on stc.id = c.fk_subtipos_conquistas_id
		where p.personagem_nome = $1
		  and tc.nome = $2
		  and stc.nome = $3;
	`)
	if err != nil {
		return model.Character{}, err
	}
	rows, err := stmt.Query(char, achType, achSubtype)
	if err != nil {
		return model.Character{}, err
	}

	var c model.Character
	var achs []model.Achievement
	for rows.Next() {
		var a model.Achievement
		err := rows.Scan(
			&c.Name,
			&c.Title,
			&c.Spec,
			&c.Class,
			&c.Race,
			&c.Guild,
			&c.Level,
			&c.ILvl,
			&c.Faction,
			&c.AchievementsScore,
			&a.Type,
			&a.SubType,
			&a.Name,
			&a.Points,
			&a.Date,
		)
		if err != nil {
			return model.Character{}, err
		}
		achs = append(achs, a)
	}
	c.Achievements = achs
	return c, nil
}

func (q *queryer) CharacterAndAchievementsPanel(char string) (model.CharacterWithAchievementsPanel, error) {
	stmt, err := q.db.Prepare(`
		select tc.nome as tipo_conquista,
			   sum(c.pontos),
			   p.personagem_nome,
			   p.titulo,
			   p.especializacao_nome,
			   p.classe_nome,
			   p.raca_nome,
			   p.guilda_nome,
			   p.nivel,
			   p.ilvl,
			   p.faccao,
			   p.escore_conquistas
		from personagens_view p
				 join personagens_conquistas pc on pc.fk_personagens_id = p.id
				 join conquistas c on c.id = pc.fk_conquistas_id
				 join tipos_conquistas tc on tc.id = c.fk_tipos_conquistas_id
		where p.personagem_nome = $1
		group by tc.id, p.personagem_nome, p.titulo, p.especializacao_nome, p.classe_nome, p.raca_nome, p.guilda_nome,
				 p.nivel, p.ilvl, p.faccao, p.escore_conquistas
		order by tc.id;
	`)
	if err != nil {
		return model.CharacterWithAchievementsPanel{}, err
	}
	rows, err := stmt.Query(char)
	if err != nil {
		return model.CharacterWithAchievementsPanel{}, err
	}

	var c model.CharacterWithAchievementsPanel
	var types []model.AchievementType
	for rows.Next() {
		var t model.AchievementType
		err := rows.Scan(
			&t.Name,
			&t.Points,
			&c.Name,
			&c.Title,
			&c.Spec,
			&c.Class,
			&c.Race,
			&c.Guild,
			&c.Level,
			&c.ILvl,
			&c.Faction,
			&c.AchievementsScore,
		)
		if err != nil {
			return model.CharacterWithAchievementsPanel{}, err
		}
		types = append(types, t)
	}
	c.AchievementPanel = types
	return c, nil
}

func (q *queryer) CharacterAchievementPointsPerMonth(char string, minPoints int) ([]model.AchievPointsInMonth, error) {
	stmt, err := q.db.Prepare(`
		select extract(year from pc."data") as year, extract(month from pc."data") as month, sum(c.pontos)
		from personagens_view p
				 join personagens_conquistas pc on pc.fk_personagens_id = p.id
				 join conquistas c on c.id = pc.fk_conquistas_id
		where p.personagem_nome = $1
		group by extract(month from pc."data"), extract(year from pc."data")
		having sum(c.pontos) >= $2
		order by extract(year from pc."data") desc, extract(month from pc."data") desc;
	`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(char, minPoints)
	if err != nil {
		return nil, err
	}

	var months []model.AchievPointsInMonth
	for rows.Next() {
		var m model.AchievPointsInMonth
		err := rows.Scan(&m.Year, &m.Month, &m.Points)
		if err != nil {
			return nil, err
		}
		months = append(months, m)
	}
	return months, nil
}

func (q *queryer) GuildAchievements(guild string) ([]model.Achievement, error) {
	stmt, err := q.db.Prepare(`
		select distinct c.nome, c.pontos
		from guildas g
			 join personagens p on p.fk_guilda_id = g.id
			 join personagens_conquistas pc on pc.fk_personagens_id = p.id
			 join conquistas c on c.id = pc.fk_conquistas_id
		where g.nome = $1;
	`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(guild)
	if err != nil {
		return nil, err
	}

	var achs []model.Achievement
	for rows.Next() {
		var a model.Achievement
		err := rows.Scan(&a.Name, &a.Points)
		if err != nil {
			return nil, err
		}
		achs = append(achs, a)
	}
	return achs, nil
}

func (q *queryer) UnearnedAchievements() ([]model.Achievement, error) {
	rows, err := q.db.Query(`
		select c.nome, c.pontos from conquistas c
		where not exists (
			select pc.fk_conquistas_id 
			from personagens_conquistas pc
			where pc.fk_conquistas_id = c.id
		);
	`)
	if err != nil {
		return nil, err
	}

	var achs []model.Achievement
	for rows.Next() {
		var a model.Achievement
		err = rows.Scan(&a.Name, &a.Points)
		if err != nil {
			return nil, err
		}
		achs = append(achs, a)
	}
	return achs, nil
}

func (q *queryer) AchievementsWithRequirements() ([]model.Achievement, error) {
	rows, err := q.db.Query(`
		select c.nome, c.pontos, tc.nome, stc.nome, r.nome, r.pontos
		from conquistas c
		join conquistas_requisitos cr on cr.fk_conquistas_id_original = c.id
		join conquistas r on r.id = cr.fk_conquistas_id_requisito
		join tipos_conquistas tc on c.fk_tipos_conquistas_id = tc.id
		join subtipos_conquistas stc on c.fk_subtipos_conquistas_id = stc.id;
	`)
	if err != nil {
		return nil, err
	}

	var achs []model.Achievement
	for rows.Next() {
		var a model.Achievement
		var r model.Achievement
		err = rows.Scan(&a.Name, &a.Points, &a.Type, &a.SubType, &r.Name, &r.Points)
		if err != nil {
			return nil, err
		}
		if achs == nil || a.Name != achs[len(achs)-1].Name {
			achs = append(achs, a)
		}
		achs[len(achs)-1].Requirements = append(achs[len(achs)-1].Requirements, r)
	}
	return achs, nil
}

func (q *queryer) scanCharacter(rows *sql.Rows) (model.Character, error) {
	var c model.Character
	err := rows.Scan(
		&c.Name,
		&c.Title,
		&c.Spec,
		&c.Class,
		&c.Race,
		&c.Guild,
		&c.Level,
		&c.ILvl,
		&c.Faction,
		&c.AchievementsScore,
	)
	return c, err
}
