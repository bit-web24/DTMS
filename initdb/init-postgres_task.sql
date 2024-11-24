CREATE USER bittu WITH PASSWORD 'bittu';

CREATE DATABASE tasks WITH OWNER = bittu;

\connect tasks bittu

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\q
