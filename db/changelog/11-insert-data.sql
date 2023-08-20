INSERT INTO model (brand_name, vehicle_type, seats, tank, vehicle_capacity)
VALUES
    ('Mercedes', 0, 5, 60, 1850),
    ('BMW', 0, 2, 55, 1470),
    ('Volkswagen', 1, 9, 120, 3450),
    ('Ferrari', 4, 2, 90, 1200),
    ('Scania', 2, 2, 240, 55000);

INSERT INTO enterprise (title, city, established, utc)
VALUES
    ('Yandex', 'Moscow', '1997-09-23', 3),
    ('Get', 'Tel-Aviv', '2010-11-01', 4),
    ('Uber', 'California', '2009-03-14', 8),
    ('Airbnb', 'California', '2008-08-05', -5);

INSERT INTO vehicle (model_id, enterprise_id, price, manufacture_year, mileage, color, vin, purchased_at)
VALUES
    (1, 1, 1500000, 2015, 75000, 2, '1HGFA16836L135699', '2023-03-21 18:11:09'),
    (2, 2, 4900000, 2021, 5040, 1, '1FTPX17L71KB86187', '2022-10-02 06:32:03'),
    (3, 1, 645000, 2012, 210000, 3, '5XXGM4A76DG184578', '2022-05-04 21:12:59'),
    (4, NULL, 140000, 2007, 451032, 6, 'JH4KA7551SC006828', '2021-04-23 04:21:49'),
    (5, 2, 1750000, 2016, 80500, 4, '2C4RDGCG3GR226662', '2022-07-17 14:55:28'),
    (4, 1, 1120000, 2019, 32600, 1, '1G4ZP5SZ8HU122939', '2021-08-09 09:40:12'),
    (3, NULL, 740000, 2014, 125800, 3, '1FTEX1C82AF652509', '2023-01-15 17:30:59'),
    (2, 1, 6500000, 2020, 9000, 2, 'JTHBK1GG5F2153089', '2022-06-28 12:18:07'),
    (1, 2, 2250000, 2018, 62000, 5, '1G1ZA5E08CF172580', '2023-02-09 22:45:34'),
    (3, 1, 320000, 2005, 187250, 1, '1GNKVFKD2HJ145881', '2021-03-07 03:39:16'),
    (5, NULL, 1650000, 2017, 89230, 4, '1GYS4AKJ6FR132387', '2021-06-17 10:27:43'),
    (2, 2, 7800000, 2020, 2800, 3, '1G4HP52K64U105580', '2022-04-02 18:59:28'),
    (4, 1, 530000, 2009, 153200, 5, 'WAUEFBFL4BA792251', '2022-11-30 07:08:56'),
    (1, NULL, 890000, 2014, 72000, 2, '1HGCP26888A113127', '2021-12-23 15:29:41'),
    (3, 2, 420000, 2010, 123500, 1, '2G61U5S30D9192959', '2022-09-14 23:04:37'),
    (5, 1, 1220000, 2019, 45000, 6, '3D7MX38A17G776334', '2023-05-18 09:56:22'),
    (2, NULL, 7200000, 2022, 1200, 4, '2C3CDZC95JH274140', '2021-10-30 16:35:09'),
    (4, 2, 258000, 2016, 168950, 2, 'JH4DA9363LS001532', '2023-07-05 11:28:03'),
    (1, 1, 1480000, 2013, 88900, 5, '2G61T5S30D9182959', '2022-03-12 05:47:19'),
    (3, NULL, 695000, 2008, 146380, 3, '1C4RJEAG2HC889297', '2021-07-29 19:08:51'),
    (2, 2, 4200000, 2021, 3020, 1, 'KS2GK16836L135612', '2023-01-08 14:16:33'),
    (4, 1, 1370000, 2018, 54000, 4, '1GKS2GKC1HR201984', '2021-11-20 08:37:12');

INSERT INTO driver (enterprise_id, active_car_id, age, salary, experience)
VALUES
    (1, 1, 37, 90000, 15),
    (2, NULL, 19, 24000, 1),
    (3, NULL, 28, 120000, 7);

INSERT INTO driver_vehicle (vehicle_id, driver_id)
VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (2, 2),
    (2, 3);

INSERT INTO manager (full_name, age, salary)
VALUES
    ('Ivan Smirnov', 34, 120000),
    ('Marry Green', 27, 90000);

INSERT INTO manager_enterprise (manager_id, enterprise_id)
VALUES
    (1, 1),
    (1, 2),
    (2, 2),
    (2, 3);

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

INSERT INTO report (vehicle_id, title, period, started_date, ended_date, result, report_type)
SELECT
    (random() * 10)::int + 1 as vehicle_id,
            'Report ' || i as title,
    interval '1 day' * (floor(random() * 10)::int + 1) as period,
    now() - interval '90 days' + interval '1 day' * (floor(random() * 90)::int + 1) as started_date,
    now() - interval '1 day' + interval '1 day' * (floor(random() * 90)::int + 2) as ended_date,
    'Result ' || i as result,
     (random() * 5)::int as report_type
FROM generate_series(1, 50) as i;
