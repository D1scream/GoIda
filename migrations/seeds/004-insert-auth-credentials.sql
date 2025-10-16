INSERT INTO auth_credentials (user_id, login, password, created_at, updated_at) VALUES
(1, 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'user', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (login) DO UPDATE SET 
    user_id = EXCLUDED.user_id,
    password = EXCLUDED.password,
    updated_at = CURRENT_TIMESTAMP;