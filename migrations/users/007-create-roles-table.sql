CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE roles IS 'Справочник ролей пользователей системы';
COMMENT ON COLUMN roles.id IS 'Уникальный идентификатор роли';
COMMENT ON COLUMN roles.name IS 'Название роли (user, admin)';
COMMENT ON COLUMN roles.description IS 'Описание роли и её прав доступа';
COMMENT ON COLUMN roles.created_at IS 'Дата и время создания роли';
COMMENT ON COLUMN roles.updated_at IS 'Дата и время последнего обновления роли';

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
