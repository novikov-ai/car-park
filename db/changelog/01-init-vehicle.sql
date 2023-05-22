CREATE TABLE IF NOT EXISTS vehicle (
    id                      bigserial       PRIMARY KEY,
    model_id                bigint          NOT NULL,
    price                   int             NOT NULL,
    manufacture_year        int             NOT NULL,
    mileage                 int             NOT NULL,
    color                   int             NOT NULL,
    vin                     text            NOT NULL,
    created_at              timestamptz     NOT NULL DEFAULT now(),
    updated_at              timestamptz     NOT NULL DEFAULT now(),
    deleted_at              timestamptz     NULL
);

FOREIGN KEY (model_id) REFERENCES model (id);

CREATE INDEX IF NOT EXISTS idx_vehicle_created_at ON vehicle (created_at);
CREATE INDEX IF NOT EXISTS idx_vehicle_updated_at ON vehicle (updated_at);
CREATE INDEX IF NOT EXISTS idx_vehicle_deleted_at ON vehicle (deleted_at);