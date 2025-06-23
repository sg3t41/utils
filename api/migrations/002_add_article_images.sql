-- 記事テーブルに画像関連のカラムを追加
ALTER TABLE articles ADD COLUMN IF NOT EXISTS featured_image VARCHAR(500);
ALTER TABLE articles ADD COLUMN IF NOT EXISTS thumbnail_image VARCHAR(500);

-- インデックスを追加（必要に応じて）
CREATE INDEX IF NOT EXISTS idx_articles_featured_image ON articles(featured_image);