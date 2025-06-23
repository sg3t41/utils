-- 記事画像フィールド統一のロールバック

-- 古いカラムを復元
ALTER TABLE articles ADD COLUMN IF NOT EXISTS featured_image VARCHAR(500);
ALTER TABLE articles ADD COLUMN IF NOT EXISTS thumbnail_image VARCHAR(500);
ALTER TABLE articles ADD COLUMN IF NOT EXISTS image TEXT;

-- article_imageの値をfeatured_imageに移行
UPDATE articles SET featured_image = article_image WHERE article_image IS NOT NULL;

-- article_imageカラムを削除
ALTER TABLE articles DROP COLUMN IF EXISTS article_image;

-- インデックスを復元
DROP INDEX IF EXISTS idx_articles_article_image;
CREATE INDEX idx_articles_featured_image ON articles(featured_image) WHERE featured_image IS NOT NULL;