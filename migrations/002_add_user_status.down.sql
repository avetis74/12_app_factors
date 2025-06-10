-- Rollback: Remove user status field
-- Version: 002
-- Description: Remove status field from users table

DROP INDEX IF EXISTS idx_users_status;
ALTER TABLE users DROP COLUMN IF EXISTS status; 