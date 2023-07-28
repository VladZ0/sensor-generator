BEGIN;

CREATE TABLE IF NOT EXISTS spieces
(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3'),
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS sensor_groups
(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3'),
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS sensors
(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_id INT NOT NULL,
    index INT NOT NULL,
    x FLOAT NOT NULL,
    y FLOAT NOT NULL,
    z FLOAT NOT NULL,
    data_output_rate INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3'),
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_group
        FOREIGN KEY(group_id)
        REFERENCES sensor_groups(id)
);

CREATE TABLE IF NOT EXISTS sensor_data
(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sensor_id INT NOT NULL,
    temperature FLOAT NOT NULL,
    transparency INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (now() AT TIME ZONE 'utc-3'),
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_sensor
        FOREIGN KEY(sensor_id)
        REFERENCES sensors(id)
);

CREATE TABLE IF NOT EXISTS detected_spieces
(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    spiece_id INT NOT NULL,
    sensor_data_id INT NOT NULL,
    CONSTRAINT fk_spiece
        FOREIGN KEY(spiece_id)
        REFERENCES spieces(id)
);

END;