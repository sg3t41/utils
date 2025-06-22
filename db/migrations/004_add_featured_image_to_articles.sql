-- Add featured image support to articles table
ALTER TABLE articles ADD COLUMN featured_image VARCHAR(500);
ALTER TABLE articles ADD COLUMN thumbnail_image VARCHAR(500);

-- Add index for image queries (optional, for performance)
CREATE INDEX idx_articles_featured_image ON articles(featured_image) WHERE featured_image IS NOT NULL;