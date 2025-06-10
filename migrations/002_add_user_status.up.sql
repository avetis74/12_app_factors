-- Migration: Add user status field
-- Version: 002
-- Description: Add status field to track user state

ALTER TABLE users 
ADD COLUMN status VARCHAR(20) DEFAULT 'active' NOT NULL;

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- Update existing users to have active status
UPDATE users SET status = 'active' WHERE status IS NULL; 