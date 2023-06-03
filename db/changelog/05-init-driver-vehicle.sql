CREATE TABLE IF NOT EXISTS driver_vehicle
(
    vehicle_id    bigint NOT NULL,
    driver_id     bigint NOT NULL,
    PRIMARY KEY (vehicle_id, driver_id)
);