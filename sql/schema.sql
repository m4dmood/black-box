CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TYPE machine_status AS ENUM ('IDLE', 'RUNNING', 'OFF', 'FAILURE');
CREATE TYPE device_status AS ENUM ('RUNNING', 'FAILURE', 'DISABLED');
CREATE TYPE log_level AS ENUM ('INFO', 'WARN', 'ALERT');

CREATE TABLE machine (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    network_address VARCHAR(50),
    machine_code VARCHAR(20) UNIQUE NOT NULL,
    status machine_status DEFAULT 'OFF',
    type VARCHAR(50),
    desc_name TEXT
);

CREATE TABLE device (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    machine_id UUID NOT NULL REFERENCES machine(uuid) ON DELETE CASCADE,
    device_code VARCHAR(20) UNIQUE NOT NULL,
    network_address VARCHAR(50),
    status device_status DEFAULT 'DISABLED',
    type VARCHAR(50),
    description TEXT
);

CREATE TABLE registry (
    ts TIMESTAMPTZ NOT NULL,
    device_id UUID NOT NULL REFERENCES device(uuid),
    event_type VARCHAR(20),
    val DOUBLE PRECISION,
    event_message TEXT,
    level log_level DEFAULT 'INFO'
);

SELECT create_hypertable('registry', 'ts');

CREATE TABLE registry_errors (
    ts TIMESTAMPTZ DEFAULT NOW(),
    raw_data TEXT,          -- Salviamo l'intera riga CSV originale
    error_message TEXT,     -- Perché è finita qui? (es. "FK Violation", "Invalid Number")
    processed_at TIMESTAMPTZ DEFAULT NOW()
);