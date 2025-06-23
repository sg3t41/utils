-- 記事画像フィールドをarticle_imageに統一
-- featured_imageの値をarticle_imageに移行し、不要なカラムを削除

-- 新しいarticle_imageカラムを追加
ALTER TABLE articles ADD COLUMN IF NOT EXISTS article_image VARCHAR(500);

-- featured_imageの値をarticle_imageに移行
UPDATE articles SET article_image = featured_image WHERE featured_image IS NOT NULL;

-- 古いカラムを削除
ALTER TABLE articles DROP COLUMN IF EXISTS featured_image;
ALTER TABLE articles DROP COLUMN IF EXISTS thumbnail_image;
ALTER TABLE articles DROP COLUMN IF EXISTS image;

-- インデックスを削除・追加
DROP INDEX IF EXISTS idx_articles_featured_image;
CREATE INDEX idx_articles_article_image ON articles(article_image) WHERE article_image IS NOT NULL;