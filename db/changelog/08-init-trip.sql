CREATE TABLE IF NOT EXISTS trip (
    id                      bigserial       PRIMARY KEY,
    vehicle_id              bigint           NOT NULL,
    started_point           text,
    ended_point             text,
    started_at              timestamp, --UTC
    ended_at                timestamp, --UTC
    scheduled_at            timestamp, --UTC
    track_length            int,
    max_velocity            int,
    max_acceleration        int,
    CHECK (ended_at >= started_at)
);