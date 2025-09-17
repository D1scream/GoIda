CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMENT ON TABLE articles IS 'Статьи пользователей системы';
COMMENT ON COLUMN articles.id IS 'Уникальный идентификатор статьи';
COMMENT ON COLUMN articles.title IS 'Заголовок статьи';
COMMENT ON COLUMN articles.content IS 'Содержимое статьи';
COMMENT ON COLUMN articles.author_id IS 'Ссылка на автора статьи';
COMMENT ON COLUMN articles.created_at IS 'Дата и время создания записи';
COMMENT ON COLUMN articles.updated_at IS 'Дата и время последнего обновления записи';

CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles(author_id);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at);
