CREATE TABLE IF NOT EXISTS manager(
    id    bigserial PRIMARY KEY,
    full_name text NOT NULL,
    age int NOT NULL,
    salary int NOT NULL
);