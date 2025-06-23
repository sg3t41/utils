DROP INDEX IF EXISTS idx_users_line_user_id;
ALTER TABLE users 
DROP COLUMN IF EXISTS line_user_id,
DROP COLUMN IF EXISTS profile_image;