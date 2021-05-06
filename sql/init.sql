/* logico: */

DROP VIEW IF EXISTS personagens_view;
DROP TABLE IF EXISTS personagens_conquistas;
DROP TABLE IF EXISTS personagens_equipamentos;
DROP TABLE IF EXISTS personagens_montarias;
DROP TABLE IF EXISTS personagens;
DROP TABLE IF EXISTS especializacoes;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS conquistas_requisitos;
DROP TABLE IF EXISTS conquistas;
DROP TABLE IF EXISTS equipamentos;
DROP TABLE IF EXISTS guildas;
DROP TABLE IF EXISTS montarias;
DROP TABLE IF EXISTS partes;
DROP TABLE IF EXISTS racas;
DROP TABLE IF EXISTS subtipos_conquistas;
DROP TABLE IF EXISTS tipos_conquistas;
DROP TYPE IF EXISTS enum_faccao;

CREATE TYPE enum_faccao AS ENUM ('Aliança', 'Horda');

CREATE TABLE personagens
(
    id                   serial PRIMARY KEY,
    fk_guilda_id         serial,
    fk_especializacao_id serial,
    fk_raca_id           serial       NOT NULL,
    nome                 varchar(128) NOT NULL UNIQUE,
    titulo               varchar(128),
    nivel                integer      NOT NULL,
    saude                integer      NOT NULL,
    recurso              integer      NOT NULL,
    estamina             integer      NOT NULL,
    forca                integer      NOT NULL,
    agilidade            integer      NOT NULL,
    intelecto            integer      NOT NULL,
    aceleracao           integer      NOT NULL,
    critico              integer      NOT NULL,
    maestria             integer      NOT NULL,
    versatilidade        integer      NOT NULL,
    ilvl                 integer      NOT NULL,
    faccao               enum_faccao,
    adm                  boolean,
    reputacao            integer,
    escore_conquistas    integer      NOT NULL
);

CREATE TABLE equipamentos
(
    id               serial PRIMARY KEY,
    fk_parte_id      serial       NOT NULL,
    nome             varchar(128) NOT NULL UNIQUE,
    icone            bytea        NOT NULL,
    nivel            integer      NOT NULL,
    estamina         integer      NOT NULL,
    forca            integer      NOT NULL,
    intelecto        integer      NOT NULL,
    agilidade        integer      NOT NULL,
    critico          integer      NOT NULL,
    aceleracao       integer      NOT NULL,
    versatilidade    integer      NOT NULL,
    maestria         integer      NOT NULL,
    nivel_min        integer,
    preco            integer,
    durabilidade_max integer
);

CREATE TABLE partes
(
    id              serial PRIMARY KEY,
    nome            varchar(128) NOT NULL UNIQUE,
    maximo_equipado integer      NOT NULL,
    duravel         boolean      NOT NULL
);

CREATE TABLE conquistas
(
    id                        serial PRIMARY KEY,
    fk_tipos_conquistas_id    serial       NOT NULL,
    fk_subtipos_conquistas_id serial       NOT NULL,
    nome                      varchar(128) NOT NULL UNIQUE,
    pontos                    integer      NOT NULL,
    icone                     bytea        NOT NULL
);

CREATE TABLE tipos_conquistas
(
    id   serial PRIMARY KEY,
    nome varchar(64) NOT NULL UNIQUE
);

CREATE TABLE especializacoes
(
    id           serial PRIMARY KEY,
    fk_classe_id serial      NOT NULL,
    nome         varchar(64) NOT NULL,
    icone        bytea       NOT NULL
);

CREATE TABLE classes
(
    id    serial PRIMARY KEY,
    nome  varchar(64) NOT NULL UNIQUE,
    icone bytea       NOT NULL,
    cor   varchar(64) NOT NULL
);

CREATE TABLE montarias
(
    id        serial PRIMARY KEY,
    nome      varchar(128) NOT NULL UNIQUE,
    icone     bytea        NOT NULL,
    obtencao  varchar(255),
    descricao varchar(255)
);

CREATE TABLE guildas
(
    id    serial PRIMARY KEY,
    nome  varchar(255) NOT NULL UNIQUE,
    icone bytea        NOT NULL
);

CREATE TABLE subtipos_conquistas
(
    id                    serial PRIMARY KEY,
    fk_tipos_conquista_id serial NOT NULL,
    nome                  varchar(64)
);

CREATE TABLE conquistas_requisitos
(
    fk_conquistas_id_original  serial NOT NULL,
    fk_conquistas_id_requisito serial NOT NULL
);

CREATE TABLE personagens_conquistas
(
    fk_conquistas_id  serial NOT NULL,
    fk_personagens_id serial NOT NULL,
    data              date   NOT NULL
);

CREATE TABLE personagens_montarias
(
    fk_montaria_id    serial  NOT NULL,
    fk_personagens_id serial  NOT NULL,
    favorita          boolean NOT NULL
);

CREATE TABLE personagens_equipamentos
(
    fk_personagens_id  serial NOT NULL,
    fk_equipamentos_id serial NOT NULL,
    durabilidade       integer
);

CREATE TABLE racas
(
    id   serial PRIMARY KEY,
    nome varchar(64) NOT NULL UNIQUE
);

ALTER TABLE personagens
    ADD CONSTRAINT FK_personagens_2
        FOREIGN KEY (fk_guilda_id)
            REFERENCES guildas (id)
            ON DELETE SET NULL;

ALTER TABLE personagens
    ADD CONSTRAINT FK_personagens_3
        FOREIGN KEY (fk_especializacao_id)
            REFERENCES especializacoes (id)
            ON DELETE SET NULL;

ALTER TABLE personagens
    ADD CONSTRAINT fk_personagens_racas
        FOREIGN KEY (fk_raca_id)
            REFERENCES racas (id)
            ON DELETE SET NULL;

ALTER TABLE equipamentos
    ADD CONSTRAINT FK_equipamentos_2
        FOREIGN KEY (fk_parte_id)
            REFERENCES partes (id)
            ON DELETE CASCADE;

ALTER TABLE conquistas
    ADD CONSTRAINT FK_conquistas_2
        FOREIGN KEY (fk_tipos_conquistas_id)
            REFERENCES tipos_conquistas (id)
            ON DELETE CASCADE;

ALTER TABLE conquistas
    ADD CONSTRAINT FK_conquistas_3
        FOREIGN KEY (fk_subtipos_conquistas_id)
            REFERENCES subtipos_conquistas (id)
            ON DELETE CASCADE;

ALTER TABLE especializacoes
    ADD CONSTRAINT FK_especializacoes_2
        FOREIGN KEY (fk_classe_id)
            REFERENCES classes (id)
            ON DELETE CASCADE;

ALTER TABLE subtipos_conquistas
    ADD CONSTRAINT FK_subtipos_conquistas_2
        FOREIGN KEY (fk_tipos_conquista_id)
            REFERENCES tipos_conquistas (id)
            ON DELETE CASCADE;

ALTER TABLE conquistas_requisitos
    ADD CONSTRAINT FK_conquistas_requisitos_1
        FOREIGN KEY (fk_conquistas_id_original)
            REFERENCES conquistas (id)
            ON DELETE CASCADE;

ALTER TABLE conquistas_requisitos
    ADD CONSTRAINT FK_conquistas_requisitos_2
        FOREIGN KEY (fk_conquistas_id_requisito)
            REFERENCES conquistas (id)
            ON DELETE RESTRICT;

ALTER TABLE personagens_conquistas
    ADD CONSTRAINT FK_personagens_conquistas_1
        FOREIGN KEY (fk_conquistas_id)
            REFERENCES conquistas (id)
            ON DELETE CASCADE;

ALTER TABLE personagens_conquistas
    ADD CONSTRAINT FK_personagens_conquistas_2
        FOREIGN KEY (fk_personagens_id)
            REFERENCES personagens (id)
            ON DELETE CASCADE;

ALTER TABLE personagens_montarias
    ADD CONSTRAINT FK_personagens_montarias_1
        FOREIGN KEY (fk_montaria_id)
            REFERENCES montarias (id)
            ON DELETE CASCADE;

ALTER TABLE personagens_montarias
    ADD CONSTRAINT FK_personagens_montarias_2
        FOREIGN KEY (fk_personagens_id)
            REFERENCES personagens (id)
            ON DELETE CASCADE;

ALTER TABLE personagens_equipamentos
    ADD CONSTRAINT FK_personagens_equipamentos_1
        FOREIGN KEY (fk_personagens_id)
            REFERENCES personagens (id)
            ON DELETE CASCADE;

ALTER TABLE personagens_equipamentos
    ADD CONSTRAINT FK_personagens_equipamentos_2
        FOREIGN KEY (fk_equipamentos_id)
            REFERENCES equipamentos (id)
            ON DELETE CASCADE;

create view personagens_view as
select p.id,
       p.nome as personagem_nome,
       p.titulo,
       e.nome as especializacao_nome,
       c.nome as classe_nome,
       r.nome as raca_nome,
       g.nome as guilda_nome,
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
from personagens p
         join especializacoes e on e.id = p.fk_especializacao_id
         join classes c on c.id = e.fk_classe_id
         join racas r on r.id = p.fk_raca_id
         join guildas g on g.id = p.fk_guilda_id;