INSERT INTO model (brand_name, vehicle_type, seats, tank, vehicle_capacity)
VALUES
    ('Mercedes', 0, 5, 60, 1850),
    ('BMW', 0, 2, 55, 1470),
    ('Volkswagen', 1, 9, 120, 3450),
    ('Scania', 2, 2, 240, 55000);

INSERT INTO enterprise (title, city, established)
VALUES
    ('Yandex', 'Moscow', '1997-09-23'),
    ('Get', 'Tel-Aviv', '2010-11-01'),
    ('Uber', 'California', '2009-03-14');
    ('Airbnb', 'California', '2008-08-05');

INSERT INTO vehicle (model_id, enterprise_id, price, manufacture_year, mileage, color, vin)
VALUES
    (1, 1, 1500000, 2015, 75000, 2, '1HGFA16836L135699'),
    (2, 2, 4900000, 2021, 5040, 1, '1FTPX17L71KB86187'),
    (3,1, 645000, 2012, 210000, 3, '5XXGM4A76DG184578'),
    (4,NULL, 140000, 2007, 451032, 6, 'JH4KA7551SC006828');

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