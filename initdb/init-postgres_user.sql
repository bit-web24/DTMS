CREATE USER bittu WITH PASSWORD 'bittu';

CREATE DATABASE users WITH OWNER = bittu;

\connect users bittu

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\q
