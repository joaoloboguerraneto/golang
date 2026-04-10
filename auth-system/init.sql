-- Script de inicializacao do banco de dados
-- Este script e executado automaticamente quando o container PostgreSQL e criado

-- Criar tabela de usuarios
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    failed_attempts INT DEFAULT 0,
    is_blocked BOOLEAN DEFAULT FALSE,
    blocked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela de sessoes
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela de tokens de captcha
CREATE TABLE IF NOT EXISTS captcha_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    answer VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar indices para melhor performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_captcha_token ON captcha_tokens(token);
CREATE INDEX IF NOT EXISTS idx_captcha_expires ON captcha_tokens(expires_at);

-- Funcao para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para atualizar updated_at
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Funcao para limpar sessoes expiradas (pode ser chamada periodicamente)
CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
    DELETE FROM captcha_tokens WHERE expires_at < CURRENT_TIMESTAMP OR used = TRUE;
END;
$$ LANGUAGE plpgsql;

-- Comentarios nas tabelas
COMMENT ON TABLE users IS 'Tabela de usuarios do sistema';
COMMENT ON TABLE sessions IS 'Sessoes ativas dos usuarios';
COMMENT ON TABLE captcha_tokens IS 'Tokens de captcha para validacao anti-bot';

COMMENT ON COLUMN users.failed_attempts IS 'Numero de tentativas de login falhas';
COMMENT ON COLUMN users.is_blocked IS 'Indica se o usuario esta bloqueado';
COMMENT ON COLUMN users.blocked_at IS 'Data/hora em que o usuario foi bloqueado';
