-- Drop trigger and function
DROP TRIGGER IF EXISTS update_articles_updated_at ON articles;
DROP FUNCTION IF EXISTS update_articles_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_articles_tags;
DROP INDEX IF EXISTS idx_articles_created_at;
DROP INDEX IF EXISTS idx_articles_published_at;
DROP INDEX IF EXISTS idx_articles_author_id;
DROP INDEX IF EXISTS idx_articles_status;

-- Drop articles table
DROP TABLE IF EXISTS articles;