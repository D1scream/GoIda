INSERT INTO users (email, name) VALUES 
    ('admin@example.com', 'Admin User'),
    ('user@example.com', 'Regular User')
ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name;
