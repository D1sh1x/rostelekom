-- 002_seed_admin.sql
-- Seed admin user (username: admin, password: admin123, role: manager)
INSERT INTO users (username, password_hash, role, name)
VALUES ('admin', '$2b$12$juazUNnMngAIQ1YXccdj0.L8ZEW4vwOk7yNzWP6E5GFnG3Tq22uPe', 'manager', 'Administrator')
ON CONFLICT (username) DO NOTHING;
