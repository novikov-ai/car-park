CREATE TABLE IF NOT EXISTS model (
    id                      bigserial       PRIMARY KEY,
    brand_name              text            NOT NULL,
    vehicle_type            int             NOT NULL,
    seats                   int             NOT NULL,
    tank                    int             NOT NULL,
    vehicle_capacity        int             NOT NULL
);

CREATE TABLE IF NOT EXISTS enterprise (
                                          id                      bigserial       PRIMARY KEY,
                                          title text NOT NULL,
                                          city text NOT NULL,
                                          established             date     NOT NULL,
                                          utc int NOT NULL
);

CREATE TABLE IF NOT EXISTS vehicle (
                                       id                      bigserial       PRIMARY KEY,
                                       model_id                bigint          NOT NULL,
                                       enterprise_id           bigint          NULL,
                                       price                   int             NOT NULL,
                                       manufacture_year        int             NOT NULL,
                                       mileage                 int             NOT NULL,
                                       color                   int             NOT NULL,
                                       vin                     text            NOT NULL UNIQUE,
                                       purchased_at            timestamp       NOT NULL,
                                       created_at              timestamptz     NOT NULL DEFAULT now(),
                                       updated_at              timestamptz     NOT NULL DEFAULT now(),
                                       deleted_at              timestamptz     NULL,
                                       FOREIGN KEY (model_id) REFERENCES model (id),
                                       FOREIGN KEY (enterprise_id) REFERENCES enterprise (id)
);

CREATE INDEX IF NOT EXISTS idx_vehicle_created_at ON vehicle (created_at);
CREATE INDEX IF NOT EXISTS idx_vehicle_updated_at ON vehicle (updated_at);
CREATE INDEX IF NOT EXISTS idx_vehicle_deleted_at ON vehicle (deleted_at);

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

CREATE TABLE IF NOT EXISTS driver_vehicle
(
    vehicle_id    bigint NOT NULL,
    driver_id     bigint NOT NULL,
    PRIMARY KEY (vehicle_id, driver_id)
);

CREATE TABLE IF NOT EXISTS manager(
                                      id    bigserial PRIMARY KEY,
                                      full_name text NOT NULL,
                                      age int NOT NULL,
                                      salary int NOT NULL
);

CREATE TABLE IF NOT EXISTS manager_enterprise(
                                                 manager_id    bigint NOT NULL,
                                                 enterprise_id     bigint NOT NULL,
                                                 PRIMARY KEY (manager_id, enterprise_id)
);

CREATE TABLE IF NOT EXISTS trip (
                                    id                      bigserial       PRIMARY KEY,
                                    vehicle_id              bigint           NOT NULL,
                                    started_point           text,
                                    ended_point             text,
                                    started_at              timestamp, --UTC
                                    ended_at                timestamp, --UTC
                                    scheduled_at            timestamp, --UTC
                                    CHECK (ended_at >= started_at)
);

CREATE TABLE IF NOT EXISTS gps_point (
                                         id                      bigserial       PRIMARY KEY,
                                         trip_id                 bigserial       ,
                                         longitude               decimal NOT NULL,
                                         latitude                decimal NOT NULL,
                                         created_at              timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gps_point_vehicle_id ON gps_point (trip_id);

CREATE INDEX idx_gps_point_created_on ON gps_point (created_at);

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