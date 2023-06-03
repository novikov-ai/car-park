CREATE TABLE IF NOT EXISTS enterprise (
    id                      bigserial       PRIMARY KEY,
    title text NOT NULL,
    city text NOT NULL,
    established             date     NOT NULL
);