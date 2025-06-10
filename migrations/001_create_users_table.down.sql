-- Rollback: Drop users table
-- Version: 001
-- Description: Remove users table and all related objects

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users; 