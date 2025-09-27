INSERT INTO users (email, name, role_id) VALUES 
    ('admin@example.com', 'Admin User', (SELECT id FROM roles WHERE name = 'admin')),
    ('user@example.com', 'Regular User', (SELECT id FROM roles WHERE name = 'user')),
    ('user2@example.com', 'Regular User 2', (SELECT id FROM roles WHERE name = 'user'))
ON CONFLICT (email) DO UPDATE SET 
    name = EXCLUDED.name,
    role_id = EXCLUDED.role_id;
