-- Pulizia per sicurezza (opzionale)
TRUNCATE TABLE machine CASCADE;

-- Inseriamo una macchina
INSERT INTO machine (uuid, machine_code, status, type, desc_name)
VALUES ('550e8400-e29b-41d4-a716-446655440000', 'PRESSA_01', 'RUNNING', 'HYDRAULIC', 'Pressa Principale Reparto A');

-- Inseriamo due sensori validi
INSERT INTO device (uuid, machine_id, device_code, status, type, description)
VALUES 
('660e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'SENS_TEMP_01', 'RUNNING', 'TEMP', 'Sensore Olio'),
('660e8400-e29b-41d4-a716-446655441111', '550e8400-e29b-41d4-a716-446655440000', 'SENS_PRES_01', 'RUNNING', 'PRESSURE', 'Sensore Pressione');