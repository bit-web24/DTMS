-- Check if the database exists and drop it if it does
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM pg_database WHERE datname = 'users') THEN
        PERFORM pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'users';
        EXECUTE 'DROP DATABASE users';
    END IF;
END $$;

-- Create the user if it doesn't exist
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'bittu') THEN
        CREATE USER bittu WITH PASSWORD 'bittu';
    END IF;
END $$;

-- Create the database owned by the user
CREATE DATABASE users WITH OWNER = bittu;

-- Connect to the new database
\c users

-- install extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\q
