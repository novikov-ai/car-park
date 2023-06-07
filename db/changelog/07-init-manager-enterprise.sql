CREATE TABLE IF NOT EXISTS manager_enterprise(
    manager_id    bigint NOT NULL,
    enterprise_id     bigint NOT NULL,
    PRIMARY KEY (manager_id, enterprise_id)
);