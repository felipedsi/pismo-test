CREATE TABLE IF NOT EXISTS "operation_types" (
    "operation_type_id" INT PRIMARY KEY,
    "description" TEXT NOT NULL
);

INSERT INTO operation_types (operation_type_id, description) VALUES (1, 'COMPRA A VISTA') ON CONFLICT (operation_type_id) DO NOTHING;
INSERT INTO operation_types (operation_type_id, description) VALUES (2, 'COMPRA PARCELADA') ON CONFLICT (operation_type_id) DO NOTHING;
INSERT INTO operation_types (operation_type_id, description) VALUES (3, 'SAQUE') ON CONFLICT (operation_type_id) DO NOTHING;
INSERT INTO operation_types (operation_type_id, description) VALUES (4, 'PAGAMENTO') ON CONFLICT (operation_type_id) DO NOTHING;
