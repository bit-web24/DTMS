-- Check if the database exists and drop it if it does
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM pg_database WHERE datname = 'tasks') THEN
        EXECUTE 'DROP DATABASE tasks';
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
CREATE DATABASE tasks WITH OWNER = bittu;

-- Connect to the new database
\c tasks

-- install extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\q
