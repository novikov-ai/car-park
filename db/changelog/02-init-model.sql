CREATE TABLE IF NOT EXISTS model (
    id                      bigserial       PRIMARY KEY,
    brand_name              text            NOT NULL,
    vehicle_type            int             NOT NULL,
    seats                   int             NOT NULL,
    tank                    int             NOT NULL,
    vehicle_capacity        int             NOT NULL
);