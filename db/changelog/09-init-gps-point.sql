CREATE TABLE IF NOT EXISTS gps_point (
    id                      bigserial       PRIMARY KEY,
    trip_id                 bigserial       ,
    longitude               decimal NOT NULL,
    latitude                decimal NOT NULL,
    created_at              timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gps_point_vehicle_id ON gps_point (trip_id);

CREATE INDEX idx_gps_point_created_on ON gps_point (created_at);