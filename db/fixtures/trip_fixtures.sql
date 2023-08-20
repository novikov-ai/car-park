WITH random_dates AS (
    SELECT
        ROUND(RANDOM() * 2 + 1)::int AS vehicle_id,
        TIMESTAMP '2018-01-01' + (RANDOM() * (TIMESTAMP '2023-01-01' - TIMESTAMP '2018-01-01')) AS started_at,
        TIMESTAMP '2018-01-01' + (RANDOM() * (TIMESTAMP '2023-01-01' - TIMESTAMP '2018-01-01')) AS ended_at
    FROM generate_series(1, 100)
)
INSERT INTO trip (vehicle_id, started_at, ended_at, started_point, ended_point)
SELECT
    vehicle_id,
    started_at,
    ended_at,
    CAST((RANDOM() * 41) AS NUMERIC(9, 6)) || ',' || CAST((RANDOM() * 19) AS NUMERIC(9, 6)) AS started_point,
    CAST((RANDOM() * 77) AS NUMERIC(9, 6)) || ',' || CAST((RANDOM() * 169) AS NUMERIC(9, 6)) AS ended_point
FROM random_dates
WHERE ended_at >= started_at;
