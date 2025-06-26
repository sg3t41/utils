-- リンク集テーブルの削除
DROP TRIGGER IF EXISTS trigger_update_links_updated_at ON links;
DROP FUNCTION IF EXISTS update_links_updated_at();
DROP TABLE IF EXISTS links;