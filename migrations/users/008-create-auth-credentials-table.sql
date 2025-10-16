CREATE TABLE IF NOT EXISTS auth_credentials (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMENT ON TABLE auth_credentials IS 'Данные для аутентификации пользователей (логины и пароли)';
COMMENT ON COLUMN auth_credentials.id IS 'Уникальный идентификатор записи аутентификации';
COMMENT ON COLUMN auth_credentials.user_id IS 'Ссылка на пользователя';
COMMENT ON COLUMN auth_credentials.login IS 'Логин для входа в систему';
COMMENT ON COLUMN auth_credentials.password IS 'Пароль пользователя';

CREATE INDEX IF NOT EXISTS idx_auth_credentials_user_id ON auth_credentials(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_credentials_login ON auth_credentials(login);