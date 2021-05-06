-- Q0 - query view with character (simple) + spec + class + race + guild
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
where p.personagem_nome = 'Arenae';

-- Q1 - query character view + achievements
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
where p.personagem_nome = 'Arenae';

-- Q2 - query character view + achievement panel
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
where p.personagem_nome = 'Arenae'
group by tc.id, p.personagem_nome, p.titulo, p.especializacao_nome, p.classe_nome, p.raca_nome, p.guilda_nome,
         p.nivel, p.ilvl, p.faccao, p.escore_conquistas
order by tc.id;

-- Q3 - query character view + achievement list per category
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
where p.personagem_nome = 'Arenae'
  and tc.nome = 'Exploração'
  and stc.nome = 'Battle for Azeroth';

-- Q4 - guild achievements (distinct members'achievs)
select distinct c.nome, c.pontos
from guildas g
         join personagens p on p.fk_guilda_id = g.id
         join personagens_conquistas pc on pc.fk_personagens_id = p.id
         join conquistas c on c.id = pc.fk_conquistas_id
where g.nome = 'Invicta Sanguine';

-- Q5 - query character view + equipment
-- todo: faltou o nome da parte
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
where p.personagem_nome = 'Arenae';

-- Q6 - query character view + mounts
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
where p.personagem_nome = 'Arenae';

-- Q7 - query achievement points history grouped by month since a given year (group by (month, year) having year > x)
select extract(year from pc."data") as year, extract(month from pc."data") as month, sum(c.pontos)
from personagens_view p
         join personagens_conquistas pc on pc.fk_personagens_id = p.id
         join conquistas c on c.id = pc.fk_conquistas_id
where p.personagem_nome = 'Arenae'
group by extract(month from pc."data"), extract(year from pc."data")
having sum(c.pontos) >= '20'
order by extract(year from pc."data") desc, extract(month from pc."data") desc;


-- Q8 - select all equips where ilvl > ilvl of character with max(ilvl)
select e.nome,
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
where e.nivel > (
    select round(avg(e.nivel))
    from personagens p
             join personagens_equipamentos pe on pe.fk_personagens_id = p.id
             join equipamentos e on e.id = pe.fk_equipamentos_id
    group by p.id
    order by avg(e.nivel) desc
    limit 1
);

-- Q9 - select all characters that have more achievement points than the guild with min(achiev points) (composed by members')
select p.personagem_nome,
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
where p.escore_conquistas > (
    select sum(c.pontos)
    from conquistas c
             join personagens_conquistas pc on pc.fk_conquistas_id = c.id
             join personagens p on p.id = pc.fk_personagens_id
             join guildas g on p.fk_guilda_id = g.id
    group by g.id
    order by sum(c.pontos) asc
    limit 1
);

-- Q10 - select all the achievements that no player has achieved
select c.nome, c.pontos
from conquistas c
where not exists(
        select pc.fk_conquistas_id
        from personagens_conquistas pc
        where pc.fk_conquistas_id = c.id
    );

-- Q11 - select players ordered by achievement points
select p.personagem_nome,
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
group by p.personagem_nome, p.titulo, p.especializacao_nome, p.classe_nome, p.raca_nome, p.guilda_nome,
         p.nivel, p.ilvl, p.faccao
order by sum(c.pontos) desc;

-- Q12 - select achievements with requirements
select c.nome   as conquista_nome,
       c.pontos as conquista_pontos,
       tc.nome  as tipo,
       stc.nome as subtipo,
       r.nome   as requisito_nome,
       r.pontos as requisito_pontos
from conquistas c
         join conquistas_requisitos cr on cr.fk_conquistas_id_original = c.id
         join conquistas r on r.id = cr.fk_conquistas_id_requisito
         join tipos_conquistas tc on c.fk_tipos_conquistas_id = tc.id
         join subtipos_conquistas stc on c.fk_subtipos_conquistas_id = stc.id;