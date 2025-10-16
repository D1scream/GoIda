DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_email_format') THEN
        ALTER TABLE users ADD CONSTRAINT chk_users_email_format 
        CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_name_length') THEN
        ALTER TABLE users ADD CONSTRAINT chk_users_name_length 
        CHECK (length(name) >= 2 AND length(name) <= 100);
    END IF;
END $$;
