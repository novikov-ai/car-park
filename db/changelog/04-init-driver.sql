CREATE TABLE IF NOT EXISTS driver
(
    id            bigserial PRIMARY KEY,
    enterprise_id bigint,
    active_car_id bigint NULL,
    age int NOT NULL,
    salary int NOT NULL,
    experience int NOT NULL,
    FOREIGN KEY (enterprise_id) REFERENCES enterprise (id)
);