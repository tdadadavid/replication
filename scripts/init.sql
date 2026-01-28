-- this will used across all database tables
-- create table users
CREATE SCHEMA public;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    age INTEGER NOT NULL,
    balance MONEY NOT NULL DEFAULT 0::money,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
