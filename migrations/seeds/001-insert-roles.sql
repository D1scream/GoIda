INSERT INTO roles (name, description) VALUES 
('user', 'Обычный пользователь - может создавать и редактировать только свои объекты'),
('admin', 'Администратор - может создавать, редактировать и удалять любые объекты')
ON CONFLICT (name) DO NOTHING;
