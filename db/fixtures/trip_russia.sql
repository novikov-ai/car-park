WITH random_dates AS (
    SELECT
        ROUND(RANDOM() * 2 + 1)::int AS vehicle_id,
        TIMESTAMP '2015-01-01' + (RANDOM() * (TIMESTAMP '2023-01-01' - TIMESTAMP '2015-01-01')) AS started_at,
        TIMESTAMP '2015-01-01' + (RANDOM() * (TIMESTAMP '2023-01-01' - TIMESTAMP '2015-01-01')) AS ended_at
    FROM generate_series(1, 100)
)
INSERT INTO trip (vehicle_id, started_at, ended_at, started_point, ended_point)
SELECT
    vehicle_id,
    started_at,
    ended_at,
    CONCAT(CAST((RANDOM() * (37.6156 - 30.31) + 30.31) AS NUMERIC(9, 6)), ',', CAST((RANDOM() * (59.94 - 55.7522) + 55.7522) AS NUMERIC(9, 6))) AS started_point,
    CONCAT(CAST((RANDOM() * (37.6156 - 30.31) + 30.31) AS NUMERIC(9, 6)), ',', CAST((RANDOM() * (59.94 - 55.7522) + 55.7522) AS NUMERIC(9, 6))) AS ended_point
FROM random_dates
WHERE
    vehicle_id BETWEEN 1 AND 3
  AND started_at >= (SELECT MIN(started_at) FROM random_dates)
  AND ended_at >= started_at;




