
-- Tabela de países
CREATE TABLE IF NOT EXISTS country (
    id INT NOT NULL PRIMARY KEY,
    name JSONB NOT NULL,
    cod VARCHAR(3) NOT NULL UNIQUE,
    cod_alpha_2 VARCHAR(2) NOT NULL UNIQUE
);