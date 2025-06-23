ALTER TABLE users 
ADD COLUMN line_user_id VARCHAR(255) UNIQUE,
ADD COLUMN profile_image VARCHAR(500);

CREATE INDEX idx_users_line_user_id ON users(line_user_id);