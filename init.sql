CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    chamado_id VARCHAR(20) NOT NULL,
    tipo VARCHAR(20),
    cpf_encrypted BYTEA NOT NULL,
    cpf_blind_index VARCHAR(64) UNIQUE NOT NULL,
    status_anterior TEXT,
    status_novo TEXT,
    titulo TEXT,
    descricao TEXT,
    "timestamp" TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE
);