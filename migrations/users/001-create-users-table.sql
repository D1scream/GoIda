CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

COMMENT ON TABLE users IS 'Профили пользователей системы';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN users.email IS 'Email адрес пользователя (уникальный)';
COMMENT ON COLUMN users.name IS 'Имя пользователя';
COMMENT ON COLUMN users.role_id IS 'Ссылка на роль пользователя';
COMMENT ON COLUMN users.is_deleted IS 'Флаг мягкого удаления (TRUE - удален, FALSE - активен)';
COMMENT ON COLUMN users.created_at IS 'Дата и время создания пользователя';
COMMENT ON COLUMN users.updated_at IS 'Дата и время последнего обновления пользователя';

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted);

