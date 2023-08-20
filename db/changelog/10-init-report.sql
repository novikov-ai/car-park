CREATE TABLE IF NOT EXISTS report (
    id                      bigserial       PRIMARY KEY,
    vehicle_id              bigserial,
    title                   text           NOT NULL,
    period                  interval,
    started_date            timestamp,
    ended_date              timestamp,
    result                  text,
    report_type             int,
    CHECK (ended_date >= started_date)
);