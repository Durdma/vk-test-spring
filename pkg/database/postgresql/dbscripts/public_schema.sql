SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET client_min_messages = warning;
SET row_security = off;
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;
SET search_path = public, pg_catalog;
SET default_tablespace = '';

CREATE TABLE films (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    name varchar(150) NOT NULL,
    description varchar(1000) NOT NULL,
    date date NOT NULL,
    rating float NOT NULL,
    CONSTRAINT films_pk PRIMARY KEY (id)
);

CREATE TYPE FIO AS (
    f_name varchar(150),
    s_name varchar(150),
    patronymic varchar(150)
    );

CREATE TYPE SEX AS ENUM (
    'Мужчина', 'Женщина'
);

CREATE TABLE actors (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    fio FIO,
    birthday date NOT NULL,
    sex SEX,
    CONSTRAINT actors_pk PRIMARY KEY (id)
);

CREATE TABLE actors_films (
    fk_actor_id uuid NOT NULL,
    fk_film_id uuid NOT NULL,
    PRIMARY KEY (fk_actor_id, fk_film_id),
    FOREIGN KEY (fk_actor_id) REFERENCES actors(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    FOREIGN KEY (fk_film_id) REFERENCES films(id) ON DELETE CASCADE ON UPDATE RESTRICT
)

CREATE TYPE ROLES AS ENUM (
    'пользователь', 'администратор'
)

CREATE TABLE users (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    name varchar(50),
    role ROLES NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
)